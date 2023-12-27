package main

import (
	"crypto/sha256"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func calculateWebsiteHash(url string) ([32]byte, error) {
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
func cleanHTML(htmlString string) string {
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
