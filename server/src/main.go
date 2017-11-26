package main

import (
	"fmt"
	"io"
	"net/http"
	"io/ioutil"
	"log"
	"github.com/antonholmquist/jason"
	"time"
)

const VALIDATOR_HTML = "<html><head><script src='http://127.0.0.1:8080/un-scrapy-client-bundle.js'></script></head></html>"
func generalTrafficHandler(w http.ResponseWriter, r *http.Request) {
	var cookie, err = r.Cookie("VALIDATED")
	if err != nil || cookie.Value == "" {
		fmt.Println("NOT VALIDATED" )
		io.WriteString(w, VALIDATOR_HTML)
		return
	}
	// TODO: Proxy pass to the backend and pipe the response
	log.Println("VALIDATED")
	io.WriteString(w, "Success"+cookie.Value)
}

func unScrapyResultHandler(w http.ResponseWriter, r *http.Request) {
	var bytes, err = ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	log.Print(string(bytes))
	var payload, jsonErr = jason.NewObjectFromBytes(bytes)
	if jsonErr != nil {
		log.Printf("Error decoding JSON: %v", jsonErr)
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}
	hasLiedLanguages, hasLiedLanguagesErr := payload.GetBoolean("fp", "has_lied_languages")
	hasLiedResolution, hasLiedResolutionErr := payload.GetBoolean("fp", "has_lied_resolution")
	hasLiedOs, hasLiedOsErr := payload.GetBoolean("fp", "has_lied_os")
	hasLiedBrowser, hasLiedBrowserErr := payload.GetBoolean("fp", "has_lied_browser")

	if hasLiedLanguagesErr != nil || hasLiedResolutionErr != nil || hasLiedOsErr != nil || hasLiedBrowserErr != nil {
		log.Printf("Error while reading lied Detection Values %v %v %v %v", hasLiedLanguagesErr, hasLiedResolutionErr, hasLiedOsErr, hasLiedBrowserErr)
		http.Error(w, "Error while reading lied Detection Values", http.StatusBadRequest)
		return
	}

	if hasLiedLanguages || hasLiedResolution || hasLiedOs || hasLiedBrowser {
		log.Print("Failure")
		io.WriteString(w, "{\"status\" : \"Failure\"}")
		return
	}
	http.SetCookie(w, &http.Cookie{Name:"VALIDATED", Value:"true", Expires: time.Now().Add(10 * time.Minute), Path: "/", MaxAge:(60*10)})
	log.Print("Success")
	io.WriteString(w, "{\"status\" : \"Success\"}")
}

func main() {
	http.HandleFunc("/", generalTrafficHandler)
	http.HandleFunc("/un-scrapy/result", unScrapyResultHandler)
	http.ListenAndServe(":9090", nil)
}
