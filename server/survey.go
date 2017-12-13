package server

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

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

func SendSurveyController(c echo.Context) error {
	// Set account keys & information
	// accountSid := "ACae26f5f2f727a525f6042bf286ba0306"
	// authToken := "c9791dcc14812f2a5f37c140b12db265"

	// Create possible message bodies
	quotes := [7]string{"I urge you to please notice when you are happy, and exclaim or murmur or think at some point, 'If this isn't nice, I don't know what is.'",
		"Peculiar travel suggestions are dancing lessons from God.",
		"There's only one rule that I know of, babies—God damn it, you've got to be kind.",
		"Many people need desperately to receive this message: 'I feel and think much as you do, care about many of the things you care about, although most people do not care about them. You are not alone.'",
		"That is my principal objection to life, I think: It's too easy, when alive, to make perfectly horrible mistakes.",
		"So it goes.",
		"We must be careful about what we pretend to be."}

	rand.Seed(time.Now().Unix())
	quote := quotes[rand.Intn(len(quotes))]

	data, _ := models.SendQuestionToPhone("+13038593854", quote)
	return c.JSON(http.StatusOK, data)
}

func ReceiveSurveyResponse(c echo.Context) error {
	var s models.ReceiveRes
	if err := c.Bind(&s); err != nil {
		fmt.Println("s", s)
	}

	phoneNumber := s.From[2:len(s.From)]

	survey, err := models.AddResponseToSurveyModel(phoneNumber, s.Body)

	if err != nil {
		fmt.Println(err)
	}

	survey, _ = models.FindSurveyModelByPhoneNumber(phoneNumber)
	if survey.Finished == false {
		campaign, _ := models.FindCampaignModel(survey.Campaign.Hex())
		models.SendQuestionToPhone(phoneNumber, campaign.Questions[len(survey.Answers)])
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
