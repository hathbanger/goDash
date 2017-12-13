package server

import (
	"fmt"
	"github.com/hathbanger/goDash/models"
	"github.com/labstack/echo"
	"net/http"
)

func CreateTeamController(c echo.Context) error {
	var t models.Team
	if err := c.Bind(&t); err != nil {
		fmt.Println("err", err)
	}

	team := models.NewTeamModel(t.TeamName, t.TeamType, t.Organization.Hex())
	err := team.Save()
	if err != nil {
		return c.JSON(
			http.StatusForbidden,
			err)
	}
	return c.JSON(http.StatusOK, "Team created!")
}

func GetAllTeamsController(c echo.Context) error {
	organizationID := c.Param("organizationID")
	organization, err := models.FindOrganizationModel(organizationID)
	if err != nil {
		return err
	}

	var teams []*models.Team

	for _, v := range organization.Teams {
		team, _ := models.FindTeamModel(v.Hex())
		teams = append(teams, &team)
	}

	return c.JSON(http.StatusOK, teams)
}

func AddUserToTeamController(c echo.Context) error {
	var t models.Team
	if err := c.Bind(&t); err != nil {
		fmt.Println("err", err)
	}

	team := models.NewTeamModel(t.TeamName, t.TeamType, t.Organization.Hex())
	err := team.Save()
	if err != nil {
		return c.JSON(
			http.StatusForbidden,
			err)
	}
	return c.JSON(http.StatusOK, "Team created!")
}

// func UpdateTeamController(c echo.Context) error {
// 	organizationID := c.Param("organizationID")
// 	organizationName := c.FormValue("organizationName")
// 	models.UpdateOrganizationModel(organizationID, organizationName)
// 	organization, err := models.FindOrganizationModel(organizationID)
// 	if err != nil {
// 		return err
// 	}

// 	return c.JSON(http.StatusOK, organization)
// }

// func RemoveTeamController(c echo.Context) error {
// 	organizationID := c.Param("organizationID")
// 	err := models.DeleteOrganizationModel(organizationID)
// 	if err != nil {
// 		return c.JSON(http.StatusNotFound, "not able to remove the organization..")
// 	}

// 	return c.JSON(http.StatusOK, "Organization deleted!")
// }
