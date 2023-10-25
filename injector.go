package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
)

// Struct to hold the CSV data
type Record struct {
	Column1 string `csv:"column1"`
	Column2 string `csv:"column2"`
	// Add more columns as needed
}

func main() {
	// Database configuration
	db, err := sql.Open("mysql", "username:password@tcp(127.0.0.1:3306)/dbname")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Directory containing the CSV files
	dir := "./csv_directory"

	// Read all files in the directory
	files, err := filepath.Glob(filepath.Join(dir, "*.csv"))
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		// Open the CSV file
		csvFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			panic(err)
		}
		defer csvFile.Close()

		// Parse the CSV file
		records := []*Record{}
		if err := csv.UnmarshalCSV(csvFile, &records); err != nil {
			panic(err)
		}

		// Insert records into the database
		for _, record := range records {
			_, err := db.Exec("INSERT INTO tablename (column1, column2) VALUES (?, ?)", record.Column1, record.Column2)
			if err != nil {
				panic(err)
			}
		}

		fmt.Printf("Successfully imported %s into the database.\n", file)
	}
}
