package main

import (
	"fmt"
	"io"
	"net/http"
	"io/ioutil"
	"log"
	"github.com/antonholmquist/jason"
)

func generalTrafficHandler(w http.ResponseWriter, r *http.Request) {
	var cookie, err = r.Cookie("VALIDATED")
	if err == nil {
		fmt.Println("VALIDATED")
	} else {
		fmt.Println("NOT VALIDATED")
	}
	io.WriteString(w, "Success"+cookie.Value)
}

func unScrapyResultHandler(w http.ResponseWriter, r *http.Request) {
	var bytes, err = ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	var payload, jsonErr = jason.NewObjectFromBytes(bytes)
	if jsonErr != nil {
		log.Printf("Error decoding JSON: %v", jsonErr)
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}
	hasLiedLanguages, hasLiedLanguagesErr := payload.GetBoolean("has_lied_languages")
	hasLiedResolution, hasLiedResolutionErr := payload.GetBoolean("has_lied_resolution")
	hasLiedOs, hasLiedOsErr := payload.GetBoolean("has_lied_os")
	hasLiedBrowser, hasLiedBrowserErr := payload.GetBoolean("has_lied_browser")

	if hasLiedLanguagesErr != nil || hasLiedResolutionErr != nil || hasLiedOsErr != nil || hasLiedBrowserErr != nil {
		log.Printf("Error while reading lied Detection Values %v %v %v %v", hasLiedLanguagesErr, hasLiedResolutionErr, hasLiedOsErr, hasLiedBrowserErr)
		http.Error(w, "Error while reading lied Detection Values", http.StatusBadRequest)
		return
	}

	if hasLiedLanguages || hasLiedResolution || hasLiedOs || hasLiedBrowser {
		log.Print("Failure")
		io.WriteString(w, "Failure")
		return
	}
	log.Print("Success")
	io.WriteString(w, "Success")
}

func main() {
	http.HandleFunc("/", generalTrafficHandler)
	http.HandleFunc("/un-scrapy/result", unScrapyResultHandler)
	http.ListenAndServe(":8080", nil)
}
