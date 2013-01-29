// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/bmizerany/pat"
	"github.com/darkhelmet/env"
	"github.com/jmcvetta/randutil"
	"github.com/stathat/go"
	"io/ioutil"
	"labix.org/v2/mgo"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	db         *mgo.Database
	statPrefix string // Prefix for stat names to track on StatHat
	ezkey      string // EZ Key for StatHat
)

type Choice struct {
	Text string
}

type DecisionRequest struct {
	Quandary string
	Choices  []Choice
}

type DecisionResponse struct {
	Decision string
}

type Decision struct {
	Quandary  string
	Choices   []string
	Winner    string
	Ip        string
	Timestamp time.Time
	UserAgent string
}

// Decide receives a JSON payload containing several strings, and returns a JSON
// message containing one of said strings, selected at random.
func Decide(w http.ResponseWriter, req *http.Request) {
	//
	// Parse Payload
	//
	if req.ContentLength <= 0 {
		msg := "Content-Length must be greater than 0."
		http.Error(w, msg, http.StatusLengthRequired)
		return
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var dreq DecisionRequest
	err = json.Unmarshal(body, &dreq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if dreq.Quandary == "" {
		msg := "Must supply a quandary"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	//
	// Discard empty choices
	//
	validChoices := []string{}
	for _, choice := range dreq.Choices {
		if choice.Text != "" {
			validChoices = append(validChoices, choice.Text)
		}
	}
	if len(validChoices) < 2 {
		msg := fmt.Sprintln("Must provide at least 2 choices, but you provided", len(validChoices))
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	//
	// Make the decision
	//
	winner, err := randutil.ChoiceString(validChoices)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//
	// Save to Database
	//
	c := db.C("quandaries")
	ip := req.Header.Get("X-Forwarded-For")
	if ip != "" {
		ip = strings.Split(ip, ",")[0]
	} else {
		ip = strings.Split(req.RemoteAddr, ":")[0]
	}
	d := Decision{
		Quandary:  dreq.Quandary,
		Choices:   validChoices,
		Winner:    winner,
		Ip:        ip,
		Timestamp: time.Now(),
		UserAgent: req.UserAgent(),
	}
	err = c.Insert(&d)
	if err != nil {
		msg := "MongoDB Error: " + err.Error()
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	//
	// Generate response
	//
	dres := DecisionResponse{winner}
	blob, err := json.Marshal(dres)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Add("content-type", "application/json")
	w.Write(blob)
}

// timeTrack logs the time since start to StatHat.  
// Inspired by: http://blog.stathat.com/2012/10/15/go_http_request_time.html
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	reqtimeName := statPrefix + ":reqtime:" + name
	hitsName := statPrefix + ":hits:" + name
	stathat.PostEZValue(reqtimeName, ezkey, elapsed.Seconds())
	stathat.PostEZCountOne(hitsName, ezkey)
}

// track logs the time take by an HTTP request to StatHat
func track(fn http.HandlerFunc, name string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer timeTrack(time.Now(), name)
		fn(w, req)
	}
}

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	//
	// Configuration
	//
	port := env.StringDefault("PORT", "9000")
	pwd := env.StringDefault("PWD", "/app")
	mongoUrl := env.StringDefault("MONGOLAB_URI", "localhost")
	statPrefix = env.StringDefault("STATHAT_PREFIX", "")
	ezkey = env.StringDefault("STATHAT_EZKEY", "")
	//
	// Connect to MongoDB
	//
	log.Println("Connecting to MongoDB on " + mongoUrl + "...")
	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		log.Panicln(err)
	}
	defer session.Close()
	db = session.DB("")
	_, err = db.CollectionNames()
	if err != nil {
		log.Println("Setting db name to 'decisions'.")
		db = session.DB("decisions")
	}
	//
	// Routing
	//
	mux := pat.New()
	mux.Post("/decide", track(http.HandlerFunc(Decide), "/v1/decide"))
	http.Handle("/v1/", http.StripPrefix("/v1", mux))
	http.Handle("/", http.FileServer(http.Dir(pwd+"/app")))
	//
	// Start Webserver
	//
	log.Println("Starting webserver on port " + port + "...")
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Panicln(err)
	}
}
