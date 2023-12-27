package main

import (
	database "camper/db"
	"camper/routes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db             *sql.DB
	timer          *time.Timer
	mu             sync.Mutex
	checkIntervals = []int{30, 60, 90}
	wg             sync.WaitGroup
)

func main() {
	database.InitDB()
	routes.SetupRoutes()

	wg.Add(2)
	go startTimer()
	go startHTTPServer()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Shutdown procedure
	if err := db.Close(); err != nil {
		log.Fatal(err)
	}

	wg.Wait()

}

func startTimer() {
	// Set up the initial timer
	setTimer()

	for {
		select {
		case <-time.After(1 * time.Second):
			// Check divisibility at each second
			checkDivisibility()
		}
	}
}

func setTimer() {
	mu.Lock()
	defer mu.Unlock()

	// Calculate the duration until the end of the day
	now := time.Now()
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	durationUntilEndOfDay := endOfDay.Sub(now)

	// Set the timer
	timer = time.NewTimer(durationUntilEndOfDay)
}

func checkDivisibility() {
	mu.Lock()
	defer mu.Unlock()

	currentTime := int(time.Now().Unix())
	fmt.Println("Current time:", currentTime)

	// Check divisibility at each second
	for _, interval := range checkIntervals {
		if currentTime%interval == 0 {
			fmt.Println(interval, currentTime%interval)
			go database.QueryInterval(interval)
		}
	}
}

// --------------------------------------

func startHTTPServer() {
	defer wg.Done()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
