package main

import (
	"fmt"
	"log"
	"net/http"
	"html/template"
	"encoding/json"
)

// ################## VARIABLES ################## ///


// Declare some structs to hold responses and generisize the handlers
type genericResponse struct {
	Success bool		`json:"success"`
	Error 	string	`json:"error"`
}

type Response struct {
	Success 	bool		`json:"success"`
	Error 		string	`json:"error"`
}

type GenericHandler struct {
	GET  func(http.ResponseWriter, *http.Request)
	PUT  func(http.ResponseWriter, *http.Request)
	POST func(http.ResponseWriter, *http.Request)
}


// ################## FUNCTIONS ################## ///
// 

// This function server the '/' handler which is the root. It shows the index.html file that is located in src/files/index.html
func display(w http.ResponseWriter, r *http.Request) {
	// Parse the text from the index.html
	templatePath := fmt.Sprintf("index.html")
	t, err := template.ParseFiles(templatePath)
	// If there was an error parsing files then send a response to indicate that
	if err != nil {
		log.Printf("Error parsing template at %s\n", templatePath)
		log.Println(err)

		// Package the response
		response, _ := json.Marshal(genericResponse{Success: false, Error: "Server is running"})

		// Write the response
		w.Write(response)
		return
	}

	log.Println("executing ...");

	// Execute the template file
	t.Execute(w, nil)
}

// This is a default function handler to process requests sent to correct handles but not of the correct method (i.e, GET and PUT)
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("\n\n\n\n\nGET request\n")
	log.Println("w: ", w, "\n")
	log.Println("r: ", r, "\n\n\n\n")

	// Package up an error describing the necessity of using the app
	response, _ := json.Marshal(&Response{Success: false, Error: "Server is running"})

	// Write the response
	w.Write(response)
}

// This function is for serving the HTTP handlers and setting up the GET, PUT, and POST redirects.
func (this GenericHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	if req.Method == "GET" && this.GET != nil {
		this.GET(w, req)
	} else if req.Method == "POST" && this.POST != nil {
		this.POST(w, req)
	} else if req.Method == "PUT" && this.PUT != nil {
		this.PUT(w, req)
	} else {
		http.Error(w, "GenericHandler error", http.StatusInternalServerError)
		log.Fatalf("No handler specified for the request %s", req)
	}
}

// This function is used to set the log printout format
func initializeLogger() {
	// Set the date, time, and line of code
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}


// ################## MAIN ################## ///


func main() {

	initializeLogger()
	
	// Set up handlers for different functions
	http.HandleFunc("/", display)

	log.Printf("\n\n")

	// Start the HTTP server on port 5000
	http.ListenAndServe(":5000", nil)

}