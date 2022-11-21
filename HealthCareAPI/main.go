package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitRouter() {
	router := mux.NewRouter()

	router.HandleFunc("/members", GetMembers).Methods("GET")
	router.HandleFunc("/members/{id}", GetMember).Methods("GET")
	router.HandleFunc("/members", CreateMember).Methods("POST")
	router.HandleFunc("/member/{id}", UpdateMember).Methods("Put")
	router.HandleFunc("/members/{id}", DeleteMember).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":1000", router))

}

var DB *gorm.DB
var err error

const DNS = "root:password@tcp(127.0.0.1)/memberdb?charset=utf8mb4&parseTime=True&loc=Local"

type Member struct {
	gorm.Model
	FirstName string `json:"firstname"`
	LastName  string `json:"last"`
	Email     string `json: "email"`
}

func InitMigration() {
	DB, err = gorm.Open(mysql.Open(DNS), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Wasn't able to connect to DB")
	}
	DB.AutoMigrate(&Member{})
}

func GetMembers(w http.ResponseWriter, router *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users []Member
	DB.Find(&users)
	json.NewEncoder(w).Encode(users)
}

func GetMember(w http.ResponseWriter, router *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(router)
	var user Member
	DB.First(&user, params["id"])
	json.NewEncoder(w).Encode(user)
}

func CreateMember(w http.ResponseWriter, router *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user Member
	json.NewDecoder(router.Body).Decode(&user)
	DB.Create(&user)
	json.NewEncoder(w).Encode(user)
}

func UpdateMember(w http.ResponseWriter, router *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(router)
	var user Member
	DB.First(&user, params["id"])
	json.NewDecoder(router.Body).Decode(&user)
	DB.Save(&user)
	json.NewEncoder(w).Encode(user)
}

func DeleteMember(w http.ResponseWriter, router *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(router)
	var user Member
	DB.Delete(&user, params["id"])
	json.NewEncoder(w).Encode("The USer is Deleted Successfully!")
}

func main() {
	InitMigration()
	InitRouter()

}
