package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	folderPath      = "C:/Path/To/CSVFiles"      // Specify the folder path where the CSV files are located
	uploadEndpoint  = "http://localhost:8080/upload" // Specify the HTTP upload endpoint
	checkInterval   = 5 * time.Second           // Interval for checking new files
	sentFilesRecord = "sent_files.txt"          // File to keep a record of sent files
)

func main() {
	for {
		// Check for new CSV files in the folder
		files, err := getNewCSVFiles(folderPath)
		if err != nil {
			log.Println("Error getting new CSV files:", err)
			continue
		}

		// Upload each new file to the HTTP endpoint
		for _, file := range files {
			err = uploadFile(file, uploadEndpoint)
			if err != nil {
				log.Println("Error uploading file:", err)
			} else {
				// Record the sent file in the record file
				recordSentFile(file)
			}
		}

		// Wait for the specified interval before checking again
		time.Sleep(checkInterval)
	}
}

// getNewCSVFiles retrieves new CSV files in the specified folder that have not been sent before.
func getNewCSVFiles(folderPath string) ([]string, error) {
	// Read the record file to get the list of already sent files
	sentFiles, err := readSentFilesRecord()
	if err != nil {
		return nil, err
	}

	// Get a list of all CSV files in the folder
	fileList, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	var newFiles []string
	for _, file := range fileList {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".csv") {
			filePath := filepath.Join(folderPath, file.Name())
			// Check if the file has not been sent before
			if !contains(sentFiles, filePath) {
				newFiles = append(newFiles, filePath)
			}
		}
	}

	return newFiles, nil
}

// uploadFile uploads a file to the specified HTTP endpoint.
func uploadFile(filePath, endpoint string) error {
	// Read the file content
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Create a new HTTP request
	client := resty.New()
	resp, err := client.R().
		SetFileReader("file", filepath.Base(filePath), bytes.NewReader(fileData)).
		Post(endpoint)
	if err != nil {
		return err
	}

	// Check the response status
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("upload failed with status code: %d", resp.StatusCode())
	}

	return nil
}

// recordSentFile records a sent file in the record file.
func recordSentFile(filePath string) error {
	recordFile, err := os.OpenFile(sentFilesRecord, os.O_APPEND|os.O_WRONLY|os
