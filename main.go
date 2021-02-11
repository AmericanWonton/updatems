package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

func init() {
	getCreds()
	OAuthGmailService() //Initialize Gmail Services
}

func logWriter(logMessage string) {
	//Logging info

	wd, _ := os.Getwd()
	logDir := filepath.Join(wd, "logging", "logging.txt")
	logFile, err := os.OpenFile(logDir, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)

	defer logFile.Close()

	if err != nil {
		fmt.Println("Failed opening log file")
	}

	log.SetOutput(logFile)

	log.Println(logMessage)
}

func getCreds() {
	file, err := os.Open("security/thecreds.txt")

	if err != nil {
		fmt.Printf("Trouble opening file for Amazon Credentials: %v\n", err.Error())
	}

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)
	var text []string

	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	file.Close()

	senderAddress = text[0]
	senderPWord = text[1]
	theClientID = text[2]
	theClientSecret = text[3]
	theAccessToken = text[4]
	theRefreshToken = text[5]
	accountSid = text[6]
	authToken = text[7]
	accountNum = text[8]
}

//Handles all requests coming in
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	//API for serving items
	myRouter.HandleFunc("/sendEmail", sendEmail).Methods("POST")             //Send Email to User
	myRouter.HandleFunc("/sendTextMessage", sendTextMessage).Methods("POST") //Send text to user
	//API Checking Stuff
	myRouter.HandleFunc("/testPing", testPing).Methods("POST") //Test ping for this service
	log.Fatal(http.ListenAndServe(":80", myRouter))
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) //Randomly Seed

	//Handle Requests
	handleRequests()
}

//Handles the test page
//This is a test API we can ping on our Amazon server
func testPing(w http.ResponseWriter, r *http.Request) {
	//Initialize struct for taking messages
	type TestCrudPing struct {
		TheCrudPing string `json:"TheCrudPing"`
	}
	//Collect JSON from Postman or wherever
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
	}
	//Marshal it into our type
	var postedMessage TestCrudPing
	json.Unmarshal(bs, &postedMessage)

	messageLog := "We've had a ping come in from somewhere: " + postedMessage.TheCrudPing
	logWriter(messageLog)
	fmt.Printf("%v\n", messageLog)

	//Declare data to return
	type ReturnMessage struct {
		TheErr     string `json:"TheErr"`
		ResultMsg  string `json:"ResultMsg"`
		SuccOrFail int    `json:"SuccOrFail"`
	}
	theReturnMessage := ReturnMessage{
		TheErr:     "Here's an error",
		ResultMsg:  "Yo, here's a result",
		SuccOrFail: 0,
	}
	//Send the response back
	theJSONMessage, err := json.Marshal(theReturnMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in addUser: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}
