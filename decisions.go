// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package main

import (
	"encoding/json"
	"github.com/bmizerany/pat"
	"github.com/darkhelmet/env"
	"github.com/jmcvetta/randutil"
	"io/ioutil"
	"log"
	"net/http"
)

type payload struct {
	Choices []string
}

type decision struct {
	Choice string
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
	var pl payload
	err = json.Unmarshal(body, &pl)
	log.Println(pl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	choice, err := randutil.ChoiceString(pl.Choices)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(decision{choice})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Add("content-type", "application/json")
	w.Write(bytes)
}

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	//
	// Configuration
	//
	port := env.StringDefault("PORT", "9000")
	pwd := env.StringDefault("PWD", "/app")
	//
	// Routing
	//
	mux := pat.New()
	mux.Post("/decide", http.HandlerFunc(Decide))
	http.Handle("/v1/", http.StripPrefix("/v1", mux))
	http.Handle("/", http.FileServer(http.Dir(pwd + "/angular/app")))
	//
	// Start Webserver
	//
	log.Println("Starting webserver on port " + port + "...")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Panicln(err)
	}
}
