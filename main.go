package main

import (
	"log"
	"os"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
	"github.com/joho/godotenv"
	"github.com/Sirupsen/logrus"
	"encoding/xml"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	r := mux.NewRouter()
	n := negroni.Classic()

	r.HandleFunc("/twiml", twiml)


	n.UseHandler(r)
	logrus.Info("Listening")
	http.ListenAndServe(":" + port, n)
}

func twiml(w http.ResponseWriter, r *http.Request) {
	twiml := TwiML{Say: "Hello World!"}
	x, err := xml.Marshal(twiml)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(x)
}

type TwiML struct {
	XMLName xml.Name `xml:"Response"`
	
	Say string `xml:",omitempty"`
}