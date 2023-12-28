package functions

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func CalculateWebsiteHash(url string) (string, error) {
	// Make an HTTP GET request to the specified URL
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	cleanedHTML := cleanHTML(string(body))

	hash := sha256.Sum256([]byte(cleanedHTML))

	return hex.EncodeToString(hash[:]), nil
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
