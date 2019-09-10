package main

import (
	"./config"
	"./controllers"
	"./database"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := config.ConfigInit(); err != nil {
		log.Fatal(err)
	}
	database.Connect()
	database.CreateTable()
	defer database.Cloce()
	http.HandleFunc("/baseR4/", controllers.ResourceHandler)
	log.Fatal(http.ListenAndServe( os.Getenv("SERVER_HOST")+":"+os.Getenv("SERVER_PORT"), nil))
}