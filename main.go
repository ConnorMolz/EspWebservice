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
	Moist       float64 `json:"moist"`
	Humidity    float64 `json:"humidity"`
	Temperature float64 `json:"temperature"`
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
	router.HandleFunc("/smart_plant", handleDataPost).Methods("POST")
	// Handle Data download
	router.HandleFunc("/smart_plant", handleDataGet).Methods("GET")

	return router
}

func handleDataPost(w http.ResponseWriter, r *http.Request) {
	newEntry := plantData{}
	err := json.NewDecoder(r.Body).Decode(&newEntry)
	if err != nil {
		fmt.Println(err)
		return
	}

	savePlantData(newEntry)

	_, err = fmt.Fprintf(w, "Ok")
	if err != nil {
		fmt.Printf("%v", err.Error())
		return
	}
}

func handleDataGet(w http.ResponseWriter, r *http.Request) {
	allPlants, err := getAllPlantData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(allPlants)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func savePlantData(plantData plantData) {
	fmt.Println(plantData)
	var connectionString = getDBConnectionString()
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = db.Exec(
		"INSERT INTO smart_plants(date, moist, humidity, temperature) VALUES ($1, $2, $3, $4) ",
		time.Now(),
		plantData.Moist,
		plantData.Humidity,
		plantData.Temperature,
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
	connectionString += os.Getenv("DB_HOST")
	connectionString += "/" + os.Getenv("DB_NAME") + "?sslmode=disable"
	fmt.Println(connectionString)
	return connectionString
}

func getAllPlantData() ([]plantData, error) {
	var connectionString = getDBConnectionString()
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT moist, humidity, temperature FROM smart_plants")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plants []plantData
	for rows.Next() {
		var plant plantData
		err := rows.Scan(&plant.Moist, &plant.Humidity, &plant.Temperature)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Moist: %v, Humidity: %v, Temperature: %v\n", plant.Moist, plant.Humidity, plant.Temperature)
		plants = append(plants, plant)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return plants, nil
}
