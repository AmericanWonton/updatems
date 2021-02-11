package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var accountSid string
var authToken string
var accountNum string

func sendTextMessage(w http.ResponseWriter, r *http.Request) {
	//Get the byte slice from the JSON
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
	}

	//Declare return information for JSON
	type ReturnMessage struct {
		TheErr     string `json:"TheErr"`
		ResultMsg  string `json:"ResultMsg"`
		SuccOrFail int    `json:"SuccOrFail"`
	}
	theReturnMessage := ReturnMessage{}

	//Declare DataType from JSON
	type TextInfo struct {
		TextMessage string `json:"TextMessage"`
		PhoneACode  int    `json:"PhoneACode"`
		PhoneNumber int    `json:"PhoneNumber"`
	}

	//Marshal the user data into our type
	var dataTextMessage TextInfo
	json.Unmarshal(bs, &dataTextMessage)

	/* Begin assembling Information for the text */
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"
	userAreaNum := strconv.Itoa(dataTextMessage.PhoneACode)
	userPhoneNum := strconv.Itoa(dataTextMessage.PhoneNumber)
	usersNumber := "+" + userAreaNum + userPhoneNum

	msgData := url.Values{}
	msgData.Set("To", usersNumber)  //This is a User's number you send a text to
	msgData.Set("From", accountNum) //This is your Twilio Number where the text is sent from
	msgData.Set("Body", dataTextMessage.TextMessage)
	msgDataReader := *strings.NewReader(msgData.Encode())

	//Send Message
	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//Get response from Twilio API and log messages returned
	resp, theErr := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			fmt.Println(data["sid"])
			//Successful printing, return success
			message := "Successful text sent to: " + usersNumber
			logWriter(message)
			theReturnMessage.ResultMsg = message
			theReturnMessage.TheErr = ""
			theReturnMessage.SuccOrFail = 0
		} else {
			//Un-Successful printing, return success
			message := "Error sending text sent to: " + usersNumber + "\n" + err.Error()
			logWriter(message)
			theReturnMessage.ResultMsg = message
			theReturnMessage.TheErr = message
			theReturnMessage.SuccOrFail = 1
		}
	} else {
		fmt.Println(resp.Status)
		//Un-Successful printing, return success
		message := "Error sending text and getting response sent to: " +
			usersNumber + "\n" + theErr.Error() + resp.Status
		logWriter(message)
		theReturnMessage.ResultMsg = message
		theReturnMessage.TheErr = message
		theReturnMessage.SuccOrFail = 1
	}
	//Format JSON and send response back
	theJSONMessage, err := json.Marshal(theReturnMessage)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
	}
	fmt.Fprint(w, string(theJSONMessage))
}
