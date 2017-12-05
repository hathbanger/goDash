package server

import (
	"net/http"
	"fmt"


  	
  	"strings"
  	"math/rand"
  	"time"
  	"net/url"
  	"encoding/json"


	"github.com/hathbanger/goDash/models"
	"github.com/labstack/echo"
	// "github.com/dgrijalva/jwt-go"
)

func CreateSurveyController(c echo.Context) error {
	var s models.Survey
	if err := c.Bind(&s); err != nil {
		fmt.Println("s", s)
	}
	fmt.Println("creating survey!")
	survey := models.NewSurveyModel(s.Organization.Hex(), s.Campaign.Hex(), s.User.Hex(), s.Answers)
	err := survey.Save()
	if err != nil {
		return c.JSON(
			http.StatusForbidden,
			"We're sorry! There's already a organization with that orgName..")
	}
	return c.JSON(http.StatusOK, "Survey created!")
}

func GetSurveysController(c echo.Context) error {
	organizationID := c.Param("organizationID")
	organization, err := models.FindOrganizationModel(organizationID)
	if err != nil {
		return err
	}
	surveyArray := []models.Survey{}
	for _, v := range organization.Surveys {
		survey, _ := models.FindSurveyModel(v.Hex())
		surveyArray = append(surveyArray, survey)
	}
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, surveyArray)
}



func UpdateSurveyController(c echo.Context) error {
	var s models.Survey
	if err := c.Bind(&s); err != nil {
		fmt.Println("s", s)
	}
	models.UpdateSurveyModel(s.Id.Hex(), s.Answers)
	survey, err := models.FindSurveyModel(s.Id.Hex())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, survey)
}

func BulkResponseController(c echo.Context) error {
	var br models.BulkResponse
	if err := c.Bind(&br); err != nil {
		fmt.Println("err", err)
	}	

	for _, v := range br.Responses {
		fmt.Println("response", v)
		models.AddResponseToSurveyModel(v.From, v.Body)
	}

	return c.JSON(http.StatusOK, "survey")

}

func SendSurvey(c echo.Context) error {
  // Set account keys & information
  	accountSid := "ACae26f5f2f727a525f6042bf286ba0306"
  	authToken := "c9791dcc14812f2a5f37c140b12db265"
  	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

  	// Create possible message bodies
  	quotes := [7]string{"I urge you to please notice when you are happy, and exclaim or murmur or think at some point, 'If this isn't nice, I don't know what is.'",
                      	"Peculiar travel suggestions are dancing lessons from God.",
                      	"There's only one rule that I know of, babies—God damn it, you've got to be kind.",
                      	"Many people need desperately to receive this message: 'I feel and think much as you do, care about many of the things you care about, although most people do not care about them. You are not alone.'",
                      	"That is my principal objection to life, I think: It's too easy, when alive, to make perfectly horrible mistakes.",
                      	"So it goes.",
                      	"We must be careful about what we pretend to be."}

  	// Set up rand
  	rand.Seed(time.Now().Unix())

  // Pack up the data for our message
  	msgData := url.Values{}
  	
  	msgData.Set("To","+13038593854")
 	msgData.Set("From","+17204087088")

 	msgData.Set("Body",quotes[rand.Intn(len(quotes))])
	msgDataReader := *strings.NewReader(msgData.Encode())

  	// Create HTTP request client
  	client := &http.Client{}
  	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
  	req.SetBasicAuth(accountSid, authToken)
  	req.Header.Add("Accept", "application/json")
  	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

  	// Make HTTP POST request and return message SID
  	resp, _ := client.Do(req)
  	if (resp.StatusCode >= 200 && resp.StatusCode < 300) {
    	var data map[string]interface{}
    	decoder := json.NewDecoder(resp.Body)
    	err := decoder.Decode(&data)
    	if (err == nil) {
      		fmt.Println(data)
    	}
      		return c.JSON(http.StatusOK, data)
  	} else {
    	fmt.Println(resp);
    	var data map[string]interface{}
    	decoder := json.NewDecoder(resp.Body)
    	err := decoder.Decode(&data)    
    	if (err == nil) {
      		fmt.Println(data)
    	}    		
    	return c.JSON(http.StatusOK, data)
  	}
}

func ReceiveSurveyResponse(c echo.Context) error {
	var s models.ReceiveRes
	if err := c.Bind(&s); err != nil {
		fmt.Println("s", s)
	}

	survey, err := models.AddResponseToSurveyModel(s.From, s.Body)

	if err != nil {
		fmt.Println(err)
	}

	return c.JSON(http.StatusOK, survey)
}



// func RemoveOrganizationController(c echo.Context) error {
// 	organizationID := c.Param("organizationID")
// 	err := models.DeleteOrganizationModel(organizationID)
// 	if err != nil {
// 		return c.JSON(http.StatusNotFound, "not able to remove the organization..")
// 	}

// 	return c.JSON(http.StatusOK, "Organization deleted!")
// }

