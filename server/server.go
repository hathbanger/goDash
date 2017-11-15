
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

	e.GET("/", accessible)
	e.POST("/login", LoginUserController)
	e.POST("/user", CreateUserController)

	r.GET("/user", GetUserController)
	e.POST("/user/update", UpdateUserController)
	e.POST("/user/delete", RemoveUserController)

	r.POST("/organization", CreateOrganizationController)
	e.GET("/:organizationID", GetOrganizationController)
	e.POST("/:organizationID/update", UpdateOrganizationController)
	e.POST("/:organizationID/delete", RemoveOrganizationController)	

	e.POST("/survey", CreateSurveyController)
	e.GET("/:organizationID/get-surveys", GetSurveysController)
	// e.POST("/:organizationID/update", UpdateOrganizationController)
	// e.POST("/:organizationID/delete", RemoveOrganizationController)	


	fmt.Println("RUNNING from RUN!")

	fmt.Println("Server now running on this port: 1323")
	e.Logger.Fatal(e.Start(":1323"))
}
