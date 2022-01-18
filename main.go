package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type User struct {
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	ID          string `json:"identifier"`
	Birthdate   string `json:"birthdate"`
	Gender      string `json:"gender"`
	Phonenumber string `json:"phone_number"`
	Job         string `json:"job"`
}

type Vehicle struct {
	Plate string `json:"plate"`
	Vin   string `json:"vin"`
	Model string `json:"model"`
}

// Variables
var users []User
var vehicles []Vehicle
var db *sql.DB
var err error

func getVehicles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res, err := db.Query("SELECT plate, model, vin FROM vehicles")
	if err != nil {
		panic(err.Error())
	}
	for res.Next() {
		var vehicle Vehicle
		err = res.Scan(&vehicle.Plate, &vehicle.Model, &vehicle.Vin)
		if err != nil {
			panic(err.Error())
		}
		vehicles = append(vehicles, vehicle)
	}
	json.NewEncoder(w).Encode(vehicles)
}

func getVehicle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	println("Requesting vehicle: ", params["plate"])
	row := db.QueryRow("SELECT plate, model, vin FROM vehicles WHERE plate = ?", params["plate"])
	var vehicle Vehicle
	err = row.Scan(&vehicle.Plate, &vehicle.Model, &vehicle.Vin)
	if err != nil {
		panic(err.Error())
	}
	json.NewEncoder(w).Encode(vehicle)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res, err := db.Query("SELECT firstname, lastname, id, birthdate, sex, phone_number, job FROM characters_base")
	if err != nil {
		panic(err.Error())
	}
	for res.Next() {
		var user User
		err = res.Scan(&user.Firstname, &user.Lastname, &user.ID, &user.Birthdate, &user.Gender, &user.Phonenumber, &user.Job)
		if err != nil {
			panic(err.Error())
		}
		users = append(users, user)
	}
	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	println("Requesting user: ", params["id"])
	row := db.QueryRow("SELECT firstname, lastname, id, birthdate, sex, phone_number, job FROM characters_base WHERE id = ?", params["id"])
	var user User
	err = row.Scan(&user.Firstname, &user.Lastname, &user.ID, &user.Birthdate, &user.Gender, &user.Phonenumber, &user.Job)
	if err != nil {
		panic(err.Error())
	}
	json.NewEncoder(w).Encode(user)
}

func main() {
	r := mux.NewRouter()
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "",
		AllowNativePasswords: true,
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "",
	}
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
	println("Successfully connected do MySQL database")
	r.HandleFunc("/api/users", getUsers).Methods("GET")
	r.HandleFunc("/api/vehicles", getVehicles).Methods("GET")
	r.HandleFunc("/api/vehicles/{plate}", getVehicle).Methods("GET")
	r.HandleFunc("/api/users/{id}", getUser).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", r))
}
