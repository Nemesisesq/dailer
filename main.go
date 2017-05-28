package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/robfig/cron"
	"github.com/urfave/negroni"
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
	r.HandleFunc("/call", call)
	OneOff()

	n.UseHandler(r)
	logrus.Info("Listening on :" + port)
	http.ListenAndServe(":"+port, n)
}

func twiml(w http.ResponseWriter, r *http.Request) {
	//twiml := TwiML{Say: "Lola, It's really important to wear warm clothes in the land. The land is a harsh place, a place full of sandwiches"}
	twiml := TwiML{Play: "https://s3.us-east-2.amazonaws.com/sounds4nem/gary_v_rant_60_mins.mp3"}

	x, err := xml.MarshalIndent(twiml, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/xml")
	w.Write(x)
}

type TwiML struct {
	XMLName xml.Name `xml:"Response"`

	Say  string `xml:",omitempty"`
	Play string `xml:",omitempty"`
}

func call(w http.ResponseWriter, r *http.Request) {
	caller()

	//resp, err := MakeCall("+2165346715")

	//if err != nil {
	//	panic(err)
	//}
	//if resp.StatusCode >= 200 && resp.StatusCode < 300 {
	//	var data map[string]interface{}
	//	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	//	err := json.Unmarshal(bodyBytes, &data)
	//	if err == nil {
	//		fmt.Println(data["sid"])
	//	}
	//} else {
	//	fmt.Println(resp.Status)
	//	w.Write([]byte("Go Royals!"))
	//}

}

func MakeCall(toNum string) (*http.Response, error) {
	accountSid := "AC8babac161b27ec214bed203884635819"
	authToken := "5c575b32cf3208e7a86e849fd0cd697b"
	//callSid := "PNbf2d127871ca9856d3d06e700edbf3a1"
	urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%v/Calls", accountSid)
	v := url.Values{}
	v.Set("To", toNum)
	logrus.Info(toNum)
	v.Set("From", "+12164506822")
	call_in_number := fmt.Sprintf("%vtwiml", os.Getenv("SELF_URL"))
	logrus.Info(call_in_number)
	v.Set("Url", call_in_number)
	rb := *strings.NewReader(v.Encode())
	// Create Client
	client := &http.Client{
		Timeout: time.Second * 20,
	}
	req, err := http.NewRequest("POST", urlStr, &rb)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// make request
	resp, err := client.Do(req)
	return resp, err
}

func caller() {
	numbers := []string{
		"+12163466385",
		"+12165346715",
	}
	for _, v := range numbers {
		resp, _ := MakeCall(v)
		logrus.Info(resp)
	}

}

func OneOff() {

	tz, err := time.LoadLocation("America/New_York")

	if err != nil {
		panic(err)
	}

	c := cron.NewWithLocation(tz)

	MakeCall("+12165346715")

	c.AddFunc("0 0 4 * * 1-5", func() { MakeCall("+12165346715") })
	c.AddFunc("@every 2h", func() { MakeCall("+12165346715") })
	//c.AddFunc("@every 5s", func() { logrus.Info("making call") })
	//c.AddFunc("@hourly",      func() { fmt.Println("Every hour") })
	//c.AddFunc("@every 1h30m", func() { fmt.Println("Every hour thirty") })
	c.Start()
	//..
	// Funcs are invoked in their own goroutine, asynchronously.
	//...
	// Funcs may also be added to a running Cron
	//..
	// Inspect the cron job entries' next and previous run times.
	inspect(c.Entries())
	//..
	//c.Stop()  // Stop the scheduler (does not stop any jobs already running).
}
func inspect(entries []*cron.Entry) {
	for _, value := range entries {
		logrus.Info(*value)

	}
}
