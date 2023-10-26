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

var (
	folderPath      string = "C:/Path/To/CSVFiles" // Default folder path
	uploadEndpoint  string = "http://localhost:8080/upload"
	checkInterval   time.Duration = 5 * time.Second
	sentFilesRecord string = "sent_files.txt"
)

func main() {
	// Check if a folder path is provided as a command-line argument
	if len(os.Args) > 1 {
		folderPath = os.Args[1] // Set the folderPath from the command-line argument
	}

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

func getNewCSVFiles(folderPath string) ([]string, error) {
	sentFiles, err := readSentFilesRecord()
	if err != nil {
		return nil, err
	}

	fileList, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	var newFiles []string
	for _, file := range fileList {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".csv") {
			filePath := filepath.Join(folderPath, file.Name())
			if !contains(sentFiles, filePath) {
				newFiles = append(newFiles, filePath)
			}
		}
	}

	return newFiles, nil
}

func uploadFile(filePath, endpoint string) error {
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	client := resty.New()
	resp, err := client.R().
		SetFileReader("file", filepath.Base(filePath), bytes.NewReader(fileData)).
		Post(endpoint)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("upload failed with status code: %d", resp.StatusCode())
	}

	return nil
}

func recordSentFile(filePath string) error {
	recordFile, err := os.OpenFile(sentFilesRecord, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer recordFile.Close()

	_, err = recordFile.WriteString(filePath + "\n")
	if err != nil {
		return err
	}

	return nil
}

func readSentFilesRecord() ([]string, error) {
	data, err := ioutil.ReadFile(sentFilesRecord)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // If the file doesn't exist, return nil
		}
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	return lines, nil
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
