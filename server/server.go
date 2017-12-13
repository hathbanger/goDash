
package server

import (
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"fmt"
	
)


func Run() {
	e := echo.New()
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	// Restricted Access
	r := e.Group("")
	r.Use(mw.JWT([]byte("secret")))
	// r.GET("", GetUserController)

	// CORS
	e.Use(mw.CORSWithConfig(mw.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	// ROUTES

	e.GET("/api", accessible)
	e.POST("/api/login", LoginUserController)
	e.POST("/api/user", CreateUserController)
	
	r.GET("/api/user", GetUserController)
	e.POST("/api/user/update", UpdateUserController)
	e.POST("/api/user/upload", BulkUserController)
	e.POST("/api/user/delete", RemoveUserController)

	r.POST("/api/organization", CreateOrganizationController)
	e.GET("/api/:organizationID", GetOrganizationController)
	e.POST("/api/:organizationID/update", UpdateOrganizationController)
	e.POST("/api/:organizationID/delete", RemoveOrganizationController)	

	e.POST("/api/survey", CreateSurveyController)
	e.GET("/api/:organizationID/get-surveys", GetSurveysController)
	e.GET("/api/send-survey", SendSurveyController)
	e.POST("/api/bulk-upload", BulkResponseController)
	e.POST("/api/receive-survey", ReceiveSurveyResponse)
	e.POST("/api/survey/update", UpdateSurveyController)
	// e.POST("/:organizationID/delete", RemoveOrganizationController)	


	e.POST("/api/campaign", CreateCampaignController)
	e.GET("/api/:organizationID/get-campaigns", GetCampaignsController)
	e.GET("/api/:organizationID/:campaignId", StartCampaignController)

	e.POST("/api/team", CreateTeamController)
	e.GET("/api/:organizationID/teams", GetAllTeamsController)
	// e.GET("/api/:organizationID/get-campaigns", GetCampaignsController)
	// e.GET("/api/:organizationID/:campaignId", StartCampaignController)
	// e.POST("/:organizationID/update", UpdateOrganizationController)
	// e.POST("/:organizationID/delete", RemoveOrganizationController)	

	e.GET("/api/chatbot", ChatBot)

	fmt.Println("RUNNING from RUN!")

	fmt.Println("Server now running on this port: 1323")
	e.Logger.Fatal(e.Start(":1323"))
}
