package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dAdAbird/searchy/db"
	"github.com/dAdAbird/searchy/db/pg"
)

type session struct {
	db db.DB
}

func main() {
	dbConn := flag.String("db", "postgres://user:pass@localhost/db?sslmode=disable", "Database connecion URL")
	listenOn := flag.String("h", ":8080", "Addr to listen on")

	flag.Parse()

	db, err := pg.New(*dbConn)
	if err != nil {
		log.Fatalln("Unable to connect to db:", err)
	}

	sess := session{db: db}

	http.HandleFunc("/search", sess.handleQuery)
	s := &http.Server{
		Addr:           *listenOn,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("Listening on", *listenOn)
	log.Fatal(s.ListenAndServe())
}

func (s *session) handleQuery(w http.ResponseWriter, r *http.Request) {
	var limit, offset int

	qv := r.URL.Query()

	if val, ok := qv["l"]; ok {
		limit, _ = strconv.Atoi(val[0])
	}
	if limit == 0 {
		limit = 10
	}

	if val, ok := qv["from"]; ok {
		offset, _ = strconv.Atoi(val[0])
	}

	sites, err := s.db.Search("disintegration", limit, offset)

	if err != nil {
		http.Error(w, "Internal error: "+err.Error(), 500)
		log.Println("Request error:", err)
		return
	}

	data, err := json.Marshal(sites)

	if err != nil {
		http.Error(w, "Internal error", 500)
		log.Println("marshal error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
