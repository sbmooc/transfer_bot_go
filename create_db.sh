#!/bin/bash

set -e

sqlite3 transfer_bot_db.db "PRAGMA foreign_keys = ON;"                                                                                                                                                            
sqlite3 transfer_bot_db.db "CREATE TABLE IF NOT EXISTS managers (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, wa_number TEXT, team_name TEXT);"                                                                                         
sqlite3 transfer_bot_db.db "CREATE TABLE IF NOT EXISTS free_agents (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT);"                                                                                                          
sqlite3 transfer_bot_db.db "CREATE TABLE IF NOT EXISTS messages (id INTEGER PRIMARY KEY AUTOINCREMENT, message TEXT, source TEXT, manager_id INTEGER, timestamp INTEGER, FOREIGN KEY (manager_id) REFERENCES manager(ID));"                                                                                                          
sqlite3 transfer_bot_db.db "CREATE TABLE IF NOT EXISTS assigned_free_agents (ID INTEGER PRIMARY KEY AUTOINCREMENT, manager_id INTEGER, free_agent_id INTEGER, FOREIGN KEY (manager_id) REFERENCES manager(ID), FOREIGN KEY      
 (free_agent_id) REFERENCES free_agents(ID));"  

 # Specify the path to the CSV file                                                                                                                                                                                
 players_csv="players.csv"                                                                                                                                                                                               
 insert_players=1
 # Check if the file exists                                                                                                                                                                                        
 if [ ! -f "$players_csv" ]; then                                                                                                                                                                                     
     echo "File $players_csv not found."                                                                                                                                                                              
     insert_players=0
 fi                                                                                                                                                                                                                
 if [ $insert_players -eq 1 ]; then 

 while IFS=',' read -r name 

 do                                                                                                                                                                                                                
    sqlite3 transfer_bot_db.db "INSERT INTO free_agents (name) VALUES ('$name');"

 done < "$players_csv"    

 fi

 managers_csv="managers.csv"                                                                                                                                                                                               
 insert_managers=1
 # Check if the file exists                                                                                                                                                                                        
 if [ ! -f "$managers_csv" ]; then                                                                                                                                                                                     
     echo "File $managers_csv not found."                                                                                                                                                                              
     insert_managers=0
 fi                                                                                                                                                                                                                
 if [ $insert_managers -eq 1 ]; then 

 while IFS=',' read -r waNumber managerName teamName 

 do                                                                                                                                                                                                                
    sqlite3 transfer_bot_db.db "INSERT INTO managers (wa_number, name, team_name) VALUES ('$waNumber', '$managerName', '$teamName');"

 done < "$managers_csv"    

 fi
