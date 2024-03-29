package functions

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
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
	fmt.Println(cleanedHTML)
	if strings.Contains(cleanedHTML, "Enable JavaScript and cookies to continue") || strings.Contains(cleanedHTML, "Sorry, we just need to make sure you're not a robot.") {
		return "The site cannot be viewed: JavaScript and cookies are required", errors.New("The site cannot be viewed: JavaScript and cookies are required")
	}

	hash := sha256.Sum256([]byte(cleanedHTML))

	return hex.EncodeToString(hash[:]), nil
}

func cleanHTML(htmlString string) string {
	tokenizer := html.NewTokenizer(strings.NewReader(htmlString))
	var cleanedHTML string
	var insideScriptTag bool

	for {
		tokenType := tokenizer.Next()

		switch tokenType {
		case html.ErrorToken:
			return cleanedHTML
		case html.TextToken:
			text := strings.TrimSpace(string(tokenizer.Text()))
			if text != "" && !insideScriptTag {
				cleanedHTML += text + " "
			}
		case html.StartTagToken, html.SelfClosingTagToken:
			tagName, _ := tokenizer.TagName()
			if strings.EqualFold(string(tagName), "script") {
				insideScriptTag = true
			} else {
				cleanedHTML += "<" + string(tagName) + ">"
			}
		case html.EndTagToken:
			tagName, _ := tokenizer.TagName()
			if strings.EqualFold(string(tagName), "script") {
				insideScriptTag = false
			} else {
				cleanedHTML += "</" + string(tagName) + ">"
			}
		}
	}
}
