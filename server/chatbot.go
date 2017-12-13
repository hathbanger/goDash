package server

import (
	"net/http"
	"fmt"
	"os"

	df "github.com/meinside/dialogflow-go"
	"github.com/labstack/echo"
)

func ChatBot(c echo.Context) error {
	// variables for test
	


	// setup a client
 // for verbose messages


	response, err := SendToChatBot("hey")

	fmt.Println("err", err)

	return c.JSON(http.StatusOK, response)
	// return c.JSON(http.StatusOK, response)
}

func SendToChatBot(messageToBot string) (string, error) {
	token := os.Getenv("API_AI_ACCESS_TOKEN") // XXX - your token here
	client := df.NewClient(token)
	sessionId := "test_0123456789"	
	//client.Verbose = false
	client.Verbose = true
	// var response string

	if response, err := client.QueryText(df.QueryRequest{
		Query:     []string{messageToBot},
		SessionId: sessionId,
		Language:  df.English,
	}); err == nil {
		fmt.Printf(">>> response = %+v\n", response.Result.Fulfillment)
		return	"string", err
	} else {
		fmt.Printf("*** error: %s\n", err)
		return	"response.Result.Fulfillment", err
	}

}