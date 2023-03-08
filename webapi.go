package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// get servers from connections
func (w *WebUI) Servers(rw http.ResponseWriter, r *http.Request) {
	var servers []string = make([]string, len(w.config.Connections))
	for i, connection := range w.config.Connections {
		servers[i] = getSubstring(connection, "Server=", ";")
	}
	renderJSON(rw, servers)
}

// get databases by connection id
func (w *WebUI) Databases(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	// to-do: check range
	connection := w.config.Connections[id]

	dbReader, err := NewDbReader(connection)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dbReader.Close()

	databases, err := dbReader.GetDataBases()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	renderJSON(rw, databases)
}

// edit database description by connection id and database name
func (w *WebUI) Description(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	connection := w.config.Connections[id] // to-do: check range
	database := r.FormValue("database")
	description := r.FormValue("description")

	//log.Printf("%v %v %v %v", id, connection, database, description)
	dbReader, err := NewDbReader(connection)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dbReader.Close()

	if err := dbReader.EditDescription(database, description); err != nil {
		log.Printf("error on edit description: %v", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func renderJSON(rw http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(js)
}

// get first substring between "start" and "end" strings
func getSubstring(str string, start string, end string) string {
	s := strings.Index(str, start)
	if s == -1 {
		return ""
	}
	s += len(start)
	e := strings.Index(str[s:], end)
	if e == -1 {
		return ""
	}
	e += s
	return str[s:e]
}
