package main

import (
	database "camper/db"
	"camper/routes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/cron/v3"
	"golang.org/x/net/html"
)

var (
	c *cron.Cron
)

func main() {
	database.InitDB()
	routes.SetupRoutes()

	c = cron.New()
	c.AddFunc("@every 5s", func() {
		checkExistingCronJobs()
	})

	c.Start()
	// select {}

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func checkExistingCronJobs() {
	rows, err := database.GetWebsites()
	if err != nil {
		log.Println("Error querying database:", err)
		return
	}

	for _, website := range rows {
		var id = website.ID
		var url = website.URL
		// var hash = website.Hash
		var entryTime = website.Time
		var cronJobID = website.Cron_Job_Id

		var cron_job_id string

		// Check if a cron job is already set up for the current entry
		if isCronJobScheduled(cronJobID) == "" {
			log.Printf("Cron job already scheduled for entry with ID %d, skipping...", id)
			continue
		}

		cron_job_id = isCronJobScheduled(cronJobID)

		var cronSchedule = fmt.Sprintf("0 %d * * *", entryTime)
		fmt.Println(cronSchedule)

		c.AddFunc("@every 5s", func() {
			result := performAction(url)
			fmt.Println(result)
			database.UpdateHashResult(result, cron_job_id, id)
		})

		log.Printf("Cron job '%s' scheduled for entry with ID %d", cronJobID, id)
	}
}

func isCronJobScheduled(cronJobID string) string {
	entries := c.Entries()
	for _, entry := range entries {
		hexString := fmt.Sprintf("%x", entry.ID)
		if len(hexString) >= 4 && hexString[0:2] == "\\x" {
			hexString = hexString[2:]
		}

		if hexString == cronJobID {
			// fmt.Println("Cron job already scheduled for entry with ID", cronJobID)
			return ""
		} else {
			// fmt.Println("Cron job not scheduled for entry with ID", cronJobID)
			return hexString

		}
	}
	return ""
}

func performAction(url string) string {
	hash, err := calculateWebsiteHash(url)
	if err != nil {
		log.Println("Error calculating hash:", err)
		return "Error calculating hash"
	}

	return fmt.Sprintf("%x", hash)
}

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
