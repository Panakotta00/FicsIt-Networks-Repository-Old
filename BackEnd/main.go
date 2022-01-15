package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"log"
	"main/Database"
	"net/http"
	"os"
	"strconv"
	"time"
)

var db *pgx.Conn

func respondStruct(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "unable to dump response", http.StatusInternalServerError)
		return
	}
}

func getPackage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid package-id format", http.StatusBadRequest)
		return
	}
	pack, err := Database.PackageGet(db, id)
	if err != nil {
		http.Error(w, "package not found", http.StatusNotFound)
		return
	}
	respondStruct(w, http.StatusOK, pack)
}

func getPackageTags(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalide package-id format", http.StatusBadRequest)
		return
	}
	tags, err := Database.PackageTags(db, id)
	if err != nil {
		log.Println(err)
		http.Error(w, "package not found", http.StatusNotFound)
		return
	}
	respondStruct(w, http.StatusOK, *tags)
}

func getTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid tag-id format", http.StatusBadRequest)
		return
	}
	tag, err := Database.TagGet(db, id)
	if err != nil {
		http.Error(w, "tag not found", http.StatusNotFound)
		return
	}
	respondStruct(w, http.StatusOK, tag)
}

func getRelease(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid release-id format", http.StatusBadRequest)
		return
	}
	release, err := Database.ReleaseGet(db, id)
	if err != nil {
		http.Error(w, "release not found", http.StatusNotFound)
		return
	}
	respondStruct(w, http.StatusOK, release)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid user-id format", http.StatusBadRequest)
		return
	}
	user, err := Database.UserGet(db, id)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	respondStruct(w, http.StatusOK, user)
}

func main() {
	port, err := strconv.Atoi(os.Getenv("FINREPO_DB_PORT"))
	if err != nil || port < 0 {
		log.Fatal("Invalid Port: %v", err)
	}
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", os.Getenv("FINREPO_DB_HOST"), port, os.Getenv("FINREPO_DB_USER"), os.Getenv("FINREPO_DB_PASSWORD"), os.Getenv("FINREPO_DB_DATABASE"))
	db, err = pgx.Connect(context.Background(), connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	r := mux.NewRouter()
	r.HandleFunc("/package/{id:[0-9]+}", getPackage)
	r.HandleFunc("/package/{id:[0-9]+}/tags", getPackageTags)
	r.HandleFunc("/release/{id:[0-9]+}", getRelease)
	r.HandleFunc("/tag/{id:[0-9]+}", getTag)
	r.HandleFunc("/user/{id:[0-9]+}", getUser)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
