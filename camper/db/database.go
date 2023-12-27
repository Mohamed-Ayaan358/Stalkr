// database/database.go

package database

import (
	"database/sql"
	"fmt"
	"log"

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
}

func AddWebsite(website *models.Website) error {
	_, err := db.Exec("INSERT INTO websites (name, url, hash, time, cron_job_id) VALUES (?, ?, ?, ?, ?)", website.Name, website.URL, website.Hash, website.Time, website.Cron_Job_Id)
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
		err := rows.Scan(&website.ID, &website.Name, &website.URL, &website.Hash, &website.Time, &website.Cron_Job_Id)
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

func UpdateHashResult(result string, cronJobID string, id int) error {
	_, err := db.Exec("UPDATE websites SET hash= ?, cron_job_id = ? WHERE id = ?", result, cronJobID, id)
	if err != nil {
		log.Println("Error updating hash result in the database:", err)
		return err
	}

	return nil
}
