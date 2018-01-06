package server

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/hathbanger/goDash/models"
	"github.com/labstack/echo"
	"net/http"
)

func CreateOrganizationController(c echo.Context) error {
	var s models.Organization
	if err := c.Bind(&s); err != nil {
		fmt.Println("s", s)
	}
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userId := claims["id"].(string)

	// organizationName := c.FormValue("organizationName")
	organization := models.NewOrganizationModel(s.OrganizationName, userId)
	err := organization.Save()
	if err != nil {
		return c.JSON(
			http.StatusForbidden,
			"We're sorry! There's already a organization with that orgName..")
	}
	return c.JSON(http.StatusOK, "Organization created!")
}

func GetOrganizationController(c echo.Context) error {
	organizationID := c.Param("organizationID")
	organization, err := models.FindOrganizationModel(organizationID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, organization)
}

func UpdateOrganizationController(c echo.Context) error {
	organizationID := c.Param("organizationID")
	organizationName := c.FormValue("organizationName")
	models.UpdateOrganizationModel(organizationID, organizationName)
	organization, err := models.FindOrganizationModel(organizationID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, organization)
}

func RemoveOrganizationController(c echo.Context) error {
	organizationID := c.Param("organizationID")
	err := models.DeleteOrganizationModel(organizationID)
	if err != nil {
		return c.JSON(http.StatusNotFound, "not able to remove the organization..")
	}

	return c.JSON(http.StatusOK, "Organization deleted!")
}
