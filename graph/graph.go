// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

package main

import (
	"github.com/darkhelmet/env"
	"github.com/jmcvetta/neoism"
	"labix.org/v2/mgo"
	"log"
	"time"
	"strings"
)

var (
	mongo      *mgo.Database
	neo        *neoism.Database
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

func WriteGraph(d *Decision) error {
	tsProps := neoism.Props{"timestamp": d.Timestamp,}
	quanStr := strings.TrimSpace(strings.ToLower(d.Quandary))
	quandary, _, err := neo.GetOrCreateNode("quandaries", "quandary", neoism.Props{"quandary":quanStr,})
	if err != nil {
		return err
	}
	addrProps := neoism.Props{"ip": d.Ip, "user_agent": d.UserAgent,}
	address, _, err := neo.GetOrCreateNode("addresses", "ip", addrProps)
	if err != nil {
		return err
	}
	_, err = address.Relate("asked", quandary.Id(), tsProps)
	if err != nil {
		return err
	}
	winner := strings.ToLower(d.Winner)
	winner = strings.TrimSpace(winner)
	for _, choice := range d.Choices {
		choice = strings.ToLower(choice)
		choice = strings.TrimSpace(choice)
		c, _, err := neo.GetOrCreateNode("choices", "choice", neoism.Props{"choice":choice,})
		if err != nil {
			return err
		}
		_, err = quandary.Relate("has_choice", c.Id(), tsProps)
		if err != nil {
			return err
		}
		if choice == winner {
			_, err = quandary.Relate("winning_choice", c.Id(), tsProps)
			if err != nil {
				return err
			}

		}
	}
	log.Println(d.Quandary, d.Choices)
	return nil
}

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	//
	// Configuration
	//
	mongoUrl := env.StringDefault("MONGOLAB_URI", "localhost:27017")
	neoUrl := env.StringDefault("NEO4J_URL", "http://localhost:7474/db/data")
	//
	// Connect to MongoDB
	//
	log.Println("Connecting to MongoDB on " + mongoUrl + "...")
	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		log.Panicln(err)
	}
	defer session.Close()
	mongo = session.DB("")
	_, err = mongo.CollectionNames()
	if err != nil {
		log.Println("Setting db name to 'decisions'.")
		mongo = session.DB("decisions")
	}
	//
	// Connect to Neo4j
	//
	log.Println("Connecting to Neo4j on " + neoUrl + "...")
	neo, err = neoism.Connect(neoUrl)
	if err != nil {
		log.Println("Cannot connect to Neo4j:")
		log.Println(err)
	}
	c := mongo.C("quandaries")
	iter := c.Find(nil).Iter()
	defer iter.Close()
	d := Decision{}
	for iter.Next(&d) {
		err := WriteGraph(&d)
		if err != nil {
			log.Fatal(err)
		}
	}
}
