package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Database connection settings
	dbUser := "your_username"
	dbPassword := "your_password"
	dbName := "your_database_name"

	// Open database connection
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", dbUser, dbPassword, dbName))
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()

	// HTTP endpoint to receive the CSV file
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		file, _, err := r.FormFile("file") // Assuming the input field name is "file"
		if err != nil {
			http.Error(w, "Failed to retrieve the file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Create a new CSV reader
		reader := csv.NewReader(file)

		// Skip the header row if needed
		// reader.Read()

		// Read and insert each row into the database
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println("Error reading CSV record:", err)
				continue
			}

			// Insert the record into the database
			err = insertRecord(db, record)
			if err != nil {
				log.Println("Error inserting record into the database:", err)
				continue
			}
		}

		// Send a success response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("CSV file uploaded successfully"))
	})

	// Start the HTTP server
	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func insertRecord(db *sql.DB, record []string) error {
	// Prepare the SQL statement
	stmt, err := db.Prepare("INSERT INTO your_table_name (col1, col2, col3) VALUES (?, ?, ?)") // Modify the table and column names accordingly
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement with the record values
	_, err = stmt.Exec(record[0], record[1], record[2]) // Modify the column indices according to your CSV structure
	if err != nil {
		return err
	}

	return nil
}
