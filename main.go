package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"time"
)

// Data structs from the ESP
type plantData struct {
	moist       float64
	humidity    float64
	temperature float64
}

func main() {
	// Startup
	fmt.Println("ESP Webservice is starting...")

	router := initRouter()

	var port = getPort()
	server := http.Server{
		Addr:    port,
		Handler: router,
	}

	fmt.Printf("ESP Webservice is started and running on port%v\n", port)
	fmt.Printf("Press Ctrl+C to quit.\n")
	fmt.Printf("Start time: %v\n", time.Now())

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func initRouter() *mux.Router {
	var router = mux.NewRouter()

	// Handle Data upload
	router.HandleFunc("/", handleDataPost).Methods("POST")

	return router
}

func handleDataPost(w http.ResponseWriter, r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&plantData{})
	if err != nil {
		fmt.Println(err)
		return
	}

	savePlantData(plantData{})

	_, err = fmt.Fprintf(w, "Ok")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func savePlantData(plantData plantData) {
	var connectionString = getDBConnectionString()
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = db.Exec(
		"INSERT INTO smart_plants(date, moist, humidity, temperature) VALUES ($1, $2, $3, $4) ",
		time.Now(),
		plantData.moist,
		plantData.humidity,
		plantData.temperature,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = db.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func getPort() string {
	port := ":" + os.Getenv("PORT")
	if port == ":" {
		port = ":8080"
	}
	return port
}

func getDBConnectionString() string {
	var connectionString = "postgresql://"
	connectionString += os.Getenv("DB_USER") + ":"
	connectionString += os.Getenv("DB_PASSWORD") + "@"
	connectionString += "host=" + os.Getenv("DB_HOST") + ";"
	connectionString += "/" + os.Getenv("DB_NAME") + "?sslmode=disable"
	return connectionString
}
