package server

import (
	"net/http"
	"fmt"

	"github.com/hathbanger/goDash/models"
	"github.com/labstack/echo"
	// "github.com/dgrijalva/jwt-go"
)

func CreateCampaignController(c echo.Context) error {
	var cam models.Campaign
	if err := c.Bind(&cam); err != nil {
		fmt.Println("c", cam)
	}
	fmt.Println("creating campaign!")
	campaign := models.NewCampaignModel(cam.Organization.Hex(), cam.Questions)
	err := campaign.Save()
	if err != nil {
		return c.JSON(
			http.StatusForbidden,
			"We're sorry! There's already a organization with that orgName..")
	}
	return c.JSON(http.StatusOK, "campaign created!")
}

func GetCampaignsController(c echo.Context) error {
	organizationID := c.Param("organizationID")
	organization, err := models.FindOrganizationModel(organizationID)
	if err != nil {
		return err
	}
	campaignArray := []models.Campaign{}
	for _, v := range organization.Campaigns {
		campaign, _ := models.FindCampaignModel(v.Hex())
		campaignArray = append(campaignArray, campaign)
	}
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, campaignArray)
}

// TODO: create a start campaign endpoint that creates a survey for each user,
//       sends the 'trigger' to dialog flow to start the dialog, and then take the 
// 		 dialogflow response, add the response to the campaign questions and 
//       send it to the user to start the survey.

func StartCampaignController(c echo.Context) error {
	campaignId := c.Param("campaignId")
	campaign, err := models.StartCampaignModel(campaignId)
	if err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return c.JSON(http.StatusOK, campaign)
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

