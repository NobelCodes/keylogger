package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

const (
	filePath  = "file.txt"                       // Path to the file
	serverURL = "http://localhost:5000/upload/1" // Flask server URL
	countdown = 1 * time.Minute                  // Countdown timer set to 1 minute
)

func main() {
	for {
		// Start the countdown timer
		fmt.Println("Countdown started. Waiting for 1 minute...")
		time.Sleep(countdown)

		// Retry internet check every second if it fails
		for {
			if checkInternet() {
				fmt.Println("Internet connection is active.")
				break
			}
			fmt.Println("No internet connection. Retrying in 1 second...")
			time.Sleep(30 * time.Second)
		}

		// Try uploading the file, retry every second if it fails
		for !uploadFile(filePath) {
			fmt.Println("File upload failed. Retrying in 1 second...")
			time.Sleep(55 * time.Second)
		}

		fmt.Println("File uploaded successfully. Resetting countdown.")
	}
}

// Check internet connectivity by pinging Google
func checkInternet() bool {
	resp, err := http.Get("http://google.com")
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

// Upload the file to the server
func uploadFile(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		return false
	}
	defer file.Close()

	// Create a new buffer and multipart writer
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Add the file part
	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		fmt.Printf("Failed to create form file: %v\n", err)
		return false
	}

	// Read file content and write it to the form part
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		return false
	}
	part.Write(fileBytes)

	// Close the multipart writer
	writer.Close()

	// Create an HTTP POST request
	req, err := http.NewRequest("POST", serverURL, body)
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return false
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request to the server
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to upload file: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Upload successful.")
		return true
	}

	fmt.Printf("Upload failed with status: %s\n", resp.Status)
	return false
}
