package main

import (
        "bytes"
        "database/sql"
        "encoding/csv"
        "fmt"
        "io"
        "log"
        "net/http"
        "os"
        "strconv"
        "time"

        _ "github.com/go-sql-driver/mysql"
)

const (
        dbUsername = "your_username"
        dbPassword = "your_password"
        dbName     = "your_database"
)

func main() {
        // Open a connection to the MySQL database
        db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", dbUsername, dbPassword, dbName))
        if err != nil {
                log.Fatal("Failed to connect to the database:", err)
        }
        defer db.Close()

        // HTTP endpoint to receive the CSV data
        http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
                // Read the CSV file from the request body
                file, _, err := r.FormFile("csv")
                if err != nil {
                        http.Error(w, "Failed to read the CSV file from the request", http.StatusBadRequest)
                        return
                }
                defer file.Close()

                // Create a buffer to read the file contents
                buf := bytes.NewBuffer(nil)
                if _, err := io.Copy(buf, file); err != nil {
                        http.Error(w, "Failed to read the CSV file", http.StatusInternalServerError)
                        return
                }

                // Parse the CSV data from the buffer
                reader := csv.NewReader(bytes.NewReader(buf.Bytes()))
                records, err := reader.ReadAll()
                if err != nil {
                        http.Error(w, "Failed to parse the CSV data", http.StatusInternalServerError)
                        return
                }

                // Iterate over the records and update the MySQL database
                for _, record := range records {
                        // Assuming the first column in the CSV file is the "wyscig" number
                        wyscigNumberStr := record[0]
                        wyscigNumber, err := strconv.Atoi(wyscigNumberStr)
                        if err != nil {
                                log.Printf("Failed to convert 'wyscig' number '%s' to integer: %s\n", wyscigNumberStr, err)
                                continue
                        }

                        // Assuming the second column in the CSV file is the time value
                        timeValue := record[1]

                        // Parse the time value from the CSV file
                        parsedTime, err := time.Parse("2006-01-02 15:04:05", timeValue)
                        if err != nil {
                                log.Printf("Failed to parse time value '%s': %s\n", timeValue, err)
                                continue
                        }

                        // Update the field in the MySQL database with the time value based on the 'wyscig' number
                        query := "UPDATE your_table SET your_field = CONCAT(your_field, ?) WHERE wyscig = ?"
                        _, err = db.Exec(query, parsedTime.Format("2006-01-02 15:04:05"), wyscigNumber)
                        if err != nil {
                                log.Printf("Failed to update field for 'wyscig' number '%d': %s\n", wyscigNumber, err)
                                continue
                        }

                        log.Printf("Field updated for 'wyscig' number '%d' with time value '%s'\n", wyscigNumber, parsedTime)
                }

                // Send a success response
                w.WriteHeader(http.StatusOK)
        })

        // Start the HTTP server
        log.Fatal(http.ListenAndServe(":8080", nil))
}
