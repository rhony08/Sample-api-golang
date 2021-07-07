package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/rhony08/golang_api/module/connection"
)

const (
	baseURL = "http://www.omdbapi.com/"
)

func search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {
		var page = r.FormValue("page")
		pageParam, err := strconv.Atoi(page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		var searchWord = r.FormValue("search")

		var result []connection.Movie

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			result, err = connectionOutcall.SearchMovieData(searchWord, pageParam)
			wg.Done()
		}()

		go func() {
			_ = connectionOutcall.SaveLogData(connection.DEBUG_TYPE, "user hit search", time.Now())
			wg.Done()
		}()

		wg.Wait()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resultData, _ := json.Marshal(result)

		w.Write(resultData)
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func movieByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {
		var id = r.FormValue("id")
		var result connection.Movie
		var err error

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			result, err = connectionOutcall.GetDetailMovie(id)
			wg.Done()
		}()

		go func() {
			_ = connectionOutcall.SaveLogData(connection.DEBUG_TYPE, "user hit search", time.Now())
			wg.Done()
		}()

		wg.Wait()

		if err != nil {
			http.Error(w, "Movie not found", http.StatusBadRequest)
			return
		}

		resultData, _ := json.Marshal(result)

		w.Write(resultData)
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

var connectionOutcall connection.ConnectionInterface

func initConnection() {
	// load .env file from given path
	// we keep it empty it will load .env from current directory
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	connectionOutcall = connection.GetNewConnectionInterface(connection.OptionConnection{
		BaseURL: baseURL,
		// need to change to config
		ApiKey: os.Getenv("APIKEY"),
	})
}

func main() {
	initConnection()

	http.HandleFunc("/search", search)
	http.HandleFunc("/movie_by_id", movieByID)

	fmt.Println("starting web server at http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
