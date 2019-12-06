package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	_ "github.com/go-sql-driver/mysql"
)

var temp float32

func main() {
	handleRequests()
}

func handleRequests() {
	router := mux.NewRouter()
	router.HandleFunc("/temp", writeToDB).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func writeToDB (w http.ResponseWriter, r *http.Request) {
	reqBody,_ := ioutil.ReadAll(r.Body)
	temp, err := strconv.ParseFloat(string(reqBody),32)
	if err != nil {
		log.Fatal("something other than a float was passed")
	}

	db, err := sql.Open("mysql", DBSource())
	if err != nil {
		panic(err)
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("DB connection working")
	}

	stmt, err := db.Prepare("INSERT INTO temp VALUES (" + strconv.FormatFloat(temp, 'f', 6, 32) + ", CURRENT_TIMESTAMP)")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("updated temp successfully")
	}
}

func DBSource() string {
	file, err := os.Open(".login.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 0
	var user string
	var pass string
	for scanner.Scan() {
		if i == 0 {
			user = scanner.Text()
		} else if i == 1 {
			pass = scanner.Text()
		}
		i++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	output := user + ":" + pass + "@/esp"

	return output
}