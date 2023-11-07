package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Database connection settings
	dbUser := "" // replace with actual database user
	dbPassword := "" // replace with actual database password
	dbName := "" // replace with actual database name

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

			// Insert or update the record into the database
			err = insertOrUpdateRecord(db, record)
			if err != nil {
				log.Println("Error inserting or updating record into the database:", err)
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

func insertOrUpdateRecord(db *sql.DB, record []string) error {
	// Assuming record[1] is numerwyscigu, record[2] is tor, record[3] is surname, record[4] is name, and record[8] is czas
	// Combine surname and name for the imie column
	imie := strings.TrimSpace(record[3]) + " " + strings.TrimSpace(record[4])

	// Prepare the SQL statement for upsert
	stmt, err := db.Prepare(`
		INSERT INTO wyscigi (numerwyscigu, tor, imie, czas)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		tor=VALUES(tor), imie=VALUES(imie), czas=VALUES(czas)
	`) // Modify the table and column names accordingly
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement with the record values
	_, err = stmt.Exec(record[1], record[2], imie, record[8])
	if err != nil {
		return err
	}

	return nil
}
