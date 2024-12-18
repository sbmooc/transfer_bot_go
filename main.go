package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
	"time"
	"transfer_bot/openai"
	"transfer_bot/whatsapp"
)

type WebhookRequest struct {
	WaNumber  string
	Message   string
	Timestamp string
}

func (wr *WebhookRequest) UnmarshalJSON(data []byte) error {
	type Tempstruct struct {
		Entry []struct {
			Changes []struct {
				Value struct {
					Messages []struct {
						From      string `json:"from"`
						Timestamp string `json:"timestamp"`
						Text      struct {
							Body string `json:"body"`
						} `json:"text"`
					} `json:"messages"`
				} `json:"value"`
			} `json:"changes"`
		} `json:"entry"`
	}

	var temp Tempstruct
	if err := json.Unmarshal(data, &temp); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	if len(temp.Entry) == 0 ||
		len(temp.Entry[0].Changes) == 0 ||
		len(temp.Entry[0].Changes[0].Value.Messages) == 0 {
		return fmt.Errorf("invalid JSON structure: missing required fields")
	}

	message := temp.Entry[0].Changes[0].Value.Messages[0]
	wr.WaNumber = message.From
	wr.Timestamp = message.Timestamp
	wr.Message = message.Text.Body

	return nil
}

func main() {
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", "transfer_bot_db.db")
	if err != nil {
		log.Fatal(err)
	}
	wc := whatsapp.NewWhatsappClient("https://graph.facebook.com/v21.0/430180156836144/messages", "EAAHRbEpiYfUBO55r3LSYWo7PpXdj3UVSLpfrRLLW9B21DAS2PNyDQNzc70VsWzQ19CboZCf7qcVKnOh5zglsbqRLxOMk4SPxdXepR7Tzhk56FF4URE80pMoT79nZBYdx8b6QoVZBD5v6DvdoAHcP9ZBhZBxPUxScnSDl9iC4FAGucsTW7y2ZBC8EZAeHhwMBhwZBmGYW8TTCIfdAr4hwDqpZBEu7ZAJSQj7kjUFZCIZD")
	err = wc.ValidateConfiguration()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	http.HandleFunc("/messages", writeMessageToDB(wc, db))
	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
	// Keep the main function running
	select {}
}

func writeMessageToDB(wc *whatsapp.WhatsappClient, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// This is just to configure the webhook
			fmt.Fprintf(w, r.URL.Query().Get("hub.challenge"))
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Decode the JSON data from the request body
		var requestData WebhookRequest
		body, err := io.ReadAll(r.Body)

		err = requestData.UnmarshalJSON(body)

		if err != nil {
			log.Println(err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		var id int

		err = db.QueryRow("SELECT id from managers where wa_number = ?", requestData.WaNumber).Scan(&id)

		if err != nil {
			err = wc.SendMessage("Number not found", requestData.WaNumber)
			if err != nil {
				http.Error(w, "Error sending message back", http.StatusBadRequest)
				return
			}
			http.Error(w, "Bad request, number not found.", http.StatusBadRequest)
			return
		}

		timestamp := time.Now()

		_, err = db.Exec("INSERT INTO messages (wa_number, message, timestamp) VALUES (?, ?, ?)", requestData.WaNumber, requestData.Message, timestamp)
		if err != nil {
			log.Printf("Failed to insert message into database: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = wc.SendMessage("OK", requestData.WaNumber)

		fmt.Printf("id: %d\n", id)

	}
}

func postData(url string, message string) error {
	data := map[string]string{"message": message}
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
