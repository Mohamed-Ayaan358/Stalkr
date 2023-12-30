// database/database.go

package database

import (
	"database/sql"
	"fmt"
	"log"

	"camper/functions"
	"camper/models"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("mysql", "root:122302@tcp(localhost:3306)/todos")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to the database")

}

func CloseDB() {
	db.Close()
	fmt.Println("Closed database connection")
	return
}

func AddWebsite(website *models.Website) error {
	_, err := db.Exec("INSERT INTO websites (name, url, hash, time) VALUES (?, ?, ?, ?)", website.Name, website.URL, website.Hash, website.Time)
	if err != nil {
		log.Println("Error adding website:", err)
		return err
	}

	return nil
}

func GetWebsites() ([]models.Website, error) {
	rows, err := db.Query("SELECT * FROM websites")
	if err != nil {
		log.Println("Error getting websites:", err)
		return nil, err
	}
	defer rows.Close()

	var websites []models.Website

	for rows.Next() {
		var website models.Website
		err := rows.Scan(&website.ID, &website.Name, &website.URL, &website.Hash, &website.Time)
		if err != nil {
			log.Println("Error scanning website:", err)
			return nil, err
		}
		websites = append(websites, website)
	}

	return websites, nil
}

func DeleteWebsite(websiteID int) error {
	_, err := db.Exec("DELETE FROM websites WHERE id = ?", websiteID)
	if err != nil {
		log.Println("Error deleting website:", err)
		return err
	}

	return nil
}

func QueryInterval(interval int) {

	// Replace the query logic with your actual database query
	rows, err := db.Query("SELECT * FROM websites WHERE time % ? = 0", interval)
	if err != nil {
		log.Println("Error querying database:", err)
		return
	}
	defer rows.Close()

	var websites []models.Website

	for rows.Next() {
		var website models.Website

		err := rows.Scan(&website.ID, &website.Name, &website.URL, &website.Hash, &website.Time)
		if err != nil {
			log.Println("Error scanning website:", err)
			return

		}

		if website.Time == interval {

			var potentHash, _ = functions.CalculateWebsiteHash(website.URL)
			if website.Hash != potentHash {
				fmt.Printf("Previous hash %x\nContent Hash (SHA-256): %x\n", website.Hash, potentHash)
				db.Exec("UPDATE websites SET hash = ? WHERE id = ?", potentHash, website.ID)
			}
			websites = append(websites, website)
		}

	}
	fmt.Println("Queried", len(websites), websites)
}
