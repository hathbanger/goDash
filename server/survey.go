package server

import (
	"net/http"
	"fmt"
	"github.com/hathbanger/goDash/models"
	"github.com/labstack/echo"
	// "github.com/dgrijalva/jwt-go"
)

func CreateSurveyController(c echo.Context) error {

	// userToken := c.Get("user").(*jwt.Token)
	// claims := userToken.Claims.(jwt.MapClaims)
	// userId := claims["id"].(string)
	var s models.Survey
	if err := c.Bind(&s); err != nil {
		fmt.Println("s", s)
	}
	fmt.Println(s)
	survey := models.NewSurveyModel(s.Organization, "userId", s.Content)
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

// func UpdateOrganizationController(c echo.Context) error {
// 	organizationID := c.Param("organizationID")
// 	organizationName := c.FormValue("organizationName")
// 	models.UpdateOrganizationModel(organizationID, organizationName)
// 	organization, err := models.FindOrganizationModel(organizationID)
// 	if err != nil {
// 		return err
// 	}	

// 	return c.JSON(http.StatusOK, organization)
// }


// func RemoveOrganizationController(c echo.Context) error {
// 	organizationID := c.Param("organizationID")
// 	err := models.DeleteOrganizationModel(organizationID)
// 	if err != nil {
// 		return c.JSON(http.StatusNotFound, "not able to remove the organization..")
// 	}

// 	return c.JSON(http.StatusOK, "Organization deleted!")
// }

