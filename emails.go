package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

/* INFORMATION FOR OUR EMAIL VARIABLES */
var senderAddress string
var senderPWord string

var GmailService *gmail.Service //This gets initialized in init

var theClientID string
var theClientSecret string
var theAccessToken string
var theRefreshToken string

type UserJSON struct {
	TheName    string `json:"TheName"`
	TheEmail   string `json:"TheEmail"`
	TheMessage string `json:"TheMessage"`
}

//Initialized at begininning of program
func OAuthGmailService() {
	config := oauth2.Config{
		ClientID:     theClientID,
		ClientSecret: theClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://" + "localhost:80",
	}

	token := oauth2.Token{
		AccessToken:  theAccessToken,
		RefreshToken: theRefreshToken,
		TokenType:    "Bearer",
		Expiry:       time.Now(),
	}

	var tokenSource = config.TokenSource(context.Background(), &token)

	srv, err := gmail.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		errMsg := "Unable to retrieve Gmail client: " + err.Error()
		fmt.Println(errMsg)
		logWriter(errMsg)
	}

	GmailService = srv
	if GmailService != nil {
		succMsg := "Email service is initialized"
		logWriter(succMsg)
	}
}

func sendEmail(w http.ResponseWriter, r *http.Request) {
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
	type EmailInfo struct {
		EmailMessage string `json:"EmailString"`
		EmailAddress string `json:"EmailAddressString"`
		EmailSubject string `json:"EmailSubject"`
	}

	//Marshal the user data into our type
	var dataEmail EmailInfo
	json.Unmarshal(bs, &dataEmail)

	var message gmail.Message

	emailTo := "To: " + dataEmail.EmailMessage + "\r\n"
	subject := "Subject: " + dataEmail.EmailSubject + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	msg := []byte(emailTo + subject + mime + "\n" + dataEmail.EmailMessage)

	message.Raw = base64.URLEncoding.EncodeToString(msg)

	// Send the message
	_, err = GmailService.Users.Messages.Send("me", &message).Do()
	if err != nil {
		errMsg := "Error sending this message to the User: " + err.Error()
		theReturnMessage.TheErr = errMsg
		theReturnMessage.ResultMsg = errMsg
		theReturnMessage.SuccOrFail = 1
		fmt.Println(errMsg)
		logWriter(errMsg)
	} else {
		theReturnMessage.TheErr = ""
		theReturnMessage.ResultMsg = "Succussfully sent email to: " + dataEmail.EmailAddress
		theReturnMessage.SuccOrFail = 0
	}

	theJSONMessage, err := json.Marshal(theReturnMessage)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
	}
	fmt.Fprint(w, string(theJSONMessage))
}
