package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"main/Database"
	"net/http"
	"strconv"
)

var db *sql.DB

func respondStruct(w http.ResponseWriter, status int, data interface{}) {
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "unable to dump response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
}

func getPackage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid response-id format", http.StatusBadRequest)
		return
	}
	pack, err := Database.PackageGet(db, id)
	if err != nil {
		http.Error(w, "package not found", http.StatusNotFound)
		return
	}
	respondStruct(w, http.StatusOK, pack)
}

func updatePackage(w http.ResponseWriter, r *http.Request) {

}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/package/{id:[0-9]+}", getPackage)
	http.Handle("/", r)
}
