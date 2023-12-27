package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

var (
	websiteURL     = "http://127.0.0.1:5500/test.html"
	websiteURLChan = make(chan string, 1)
	mu             sync.Mutex
)

// func hello() {
// 	// Start a goroutine to run the ticker and check the website periodically
// 	go checkWebsitePeriodically()

// 	// Define an HTTP endpoint for receiving requests
// 	http.HandleFunc("/hash", hashHandler)
// 	http.HandleFunc("/updateURL", updateURLHandler)

// 	// Start the HTTP server on port 8080
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

func checkWebsitePeriodically() {
	// Create a ticker to run the program every 10 seconds
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mu.Lock()
			// Perform the website content check and hash calculation
			hash, err := calculateWebsiteHash(websiteURL)
			if err != nil {
				log.Println(err)
				mu.Unlock()
				continue
			}

			// Print the content hash
			fmt.Printf("Website URL: %s\nContent Hash (SHA-256): %x\n", websiteURL, hash)
			mu.Unlock()
		}
	}
}

func hashHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	// Perform the website content check and hash calculation
	hash, err := calculateWebsiteHash(websiteURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the content hash
	w.Write([]byte(fmt.Sprintf("Content Hash (SHA-256): %x", hash)))
}

func updateURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	newURL := r.FormValue("url")
	if newURL == "" {
		http.Error(w, "Missing 'url' parameter", http.StatusBadRequest)
		return
	}

	mu.Lock()
	websiteURL = newURL
	mu.Unlock()

	// Notify the goroutine to check the website with the updated URL
	websiteURLChan <- newURL

	w.Write([]byte("Website URL updated"))
}

func calculateWebsiteHash1(url string) ([32]byte, error) {
	// Make an HTTP GET request to the specified URL
	response, err := http.Get(url)
	if err != nil {
		return [32]byte{}, err
	}
	defer response.Body.Close()

	// Read the content of the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return [32]byte{}, err
	}

	// Clean up the HTML
	cleanedHTML := cleanHTML(string(body))

	// Calculate the SHA-256 hash of the cleaned HTML
	hash := sha256.Sum256([]byte(cleanedHTML))

	return hash, nil
}

// cleanHTML removes unnecessary elements from HTML, like whitespace and comments
func cleanHTML1(htmlString string) string {
	tokenizer := html.NewTokenizer(strings.NewReader(htmlString))
	var cleanedHTML string

	for {
		tokenType := tokenizer.Next()

		switch tokenType {
		case html.ErrorToken:
			return cleanedHTML
		case html.TextToken:
			text := strings.TrimSpace(string(tokenizer.Text()))
			if text != "" {
				cleanedHTML += text + " "
			}
		case html.StartTagToken, html.SelfClosingTagToken:
			tagName, _ := tokenizer.TagName()
			cleanedHTML += "<" + string(tagName) + ">"
		}
	}
}
