package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

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
