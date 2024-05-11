// Package routes stores all the routes for the Gin router.
package routes

import (
	"TaipeiCityDashboardBE/app/controllers"
	"TaipeiCityDashboardBE/app/middleware"
	"TaipeiCityDashboardBE/global"
	"encoding/json"
	"bytes"
	"io"
	"os"
	"github.com/gin-gonic/gin"
	"fmt"
	"net/http"
	"golang.org/x/net/websocket"
)

// router.go configures all API routes

var (
	Router      *gin.Engine
	RouterGroup *gin.RouterGroup
	clients = make(map[*websocket.Conn]bool) // Connected clients
)

// ConfigureRoutes configures all routes for the API and sets version router groups.
func ConfigureRoutes() {
	Router.Use(middleware.ValidateJWT)
	// API routers
	RouterGroup = Router.Group("/api/" + global.VERSION)
	configureAuthRoutes()
	configureUserRoutes()
	configureComponentRoutes()
	configureDashboardRoutes()
	configureIssueRoutes()
	RouterGroup.GET("/ws", func(c *gin.Context) {
		serveWs(c.Writer, c.Request)
	})
	RouterGroup.POST("/incident", func(c *gin.Context) {
		var buf bytes.Buffer
    _, err := io.Copy(&buf, c.Request.Body)
    if err != nil {
        // Handle error
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
        return
    }
		fmt.Println(string(buf.String()))
		c.JSON(http.StatusOK, gin.H{"message": "Hello, welcome to my Gin app!"})
	})
	RouterGroup.PUT("/write/", func(c *gin.Context) {
		var data any
		if err := c.ShouldBindJSON(&data); err != nil {
			fmt.Println(err)
		}

		// Open a file for writing (create if not exists, truncate if exists)
		file, err := os.Create("incident.geojson")
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()

		// Create a JSON encoder using the file as the output destination
		encoder := json.NewEncoder(file)

		// Encode the struct to JSON and write it to the file
		if err := encoder.Encode(data); err != nil {
			fmt.Println("Error encoding JSON:", err)
			return
		}

		fmt.Println("JSON data written to person.json")	
		c.JSON(http.StatusOK, gin.H{"message": "data write success"})	
	})
}

func configureAuthRoutes() {
	// auth routers
	authRoutes := RouterGroup.Group("/auth")
	authRoutes.Use(middleware.LimitAPIRequests(global.AuthLimitAPIRequestsTimes, global.LimitRequestsDuration))
	authRoutes.Use(middleware.LimitTotalRequests(global.AuthLimitTotalRequestsTimes, global.TokenExpirationDuration))
	authRoutes.POST("/login", controllers.Login)
	// taipeipass login callback
	authRoutes.GET("/callback", controllers.ExecIssoAuth)
	authRoutes.POST("/logout", controllers.IssoLogOut)
}

func configureUserRoutes() {
	userRoutes := RouterGroup.Group("/user")
	userRoutes.Use(middleware.LimitAPIRequests(global.UserLimitAPIRequestsTimes, global.LimitRequestsDuration))
	userRoutes.Use(middleware.LimitTotalRequests(global.UserLimitTotalRequestsTimes, global.TokenExpirationDuration))
	userRoutes.Use(middleware.IsLoggedIn())
	{
		userRoutes.GET("/me", controllers.GetUserInfo)
		userRoutes.PATCH("/me", controllers.EditUserInfo)
	}
	userRoutes.Use(middleware.IsSysAdm())
	{
		userRoutes.GET("/", controllers.GetAllUsers)
		userRoutes.PATCH("/:id", controllers.UpdateUserByID)
	}
}

// configureComponentRoutes configures all component routes.
func configureComponentRoutes() {
	componentRoutes := RouterGroup.Group("/component")

	componentRoutes.Use(middleware.LimitAPIRequests(global.ComponentLimitAPIRequestsTimes, global.LimitRequestsDuration))
	componentRoutes.Use(middleware.LimitTotalRequests(global.ComponentLimitTotalRequestsTimes, global.TokenExpirationDuration))
	{
		componentRoutes.GET("/", controllers.GetAllComponents)
		componentRoutes.
			GET("/:id", controllers.GetComponentByID)
		componentRoutes.
			GET("/:id/chart", controllers.GetComponentChartData)
		componentRoutes.GET("/:id/history", controllers.GetComponentHistoryData)
	}
	componentRoutes.Use(middleware.IsSysAdm())
	{
		componentRoutes.
			PATCH("/:id", controllers.UpdateComponent).
			DELETE("/:id", controllers.DeleteComponent)
		componentRoutes.
			PATCH("/:id/chart", controllers.UpdateComponentChartConfig)
		componentRoutes.PATCH("/:id/map", controllers.UpdateComponentMapConfig)
	}
}

func configureDashboardRoutes() {
	dashboardRoutes := RouterGroup.Group("/dashboard")
	dashboardRoutes.Use(middleware.LimitAPIRequests(global.DashboardLimitAPIRequestsTimes, global.LimitRequestsDuration))
	dashboardRoutes.Use(middleware.LimitTotalRequests(global.DashboardLimitTotalRequestsTimes, global.LimitRequestsDuration))
	{
		dashboardRoutes.
			GET("/", controllers.GetAllDashboards)
		dashboardRoutes.
			GET("/:index", controllers.GetDashboardByIndex)
	}
	dashboardRoutes.Use(middleware.IsLoggedIn())
	{
		dashboardRoutes.POST("/", controllers.CreatePersonalDashboard)
		dashboardRoutes.
			PATCH("/:index", controllers.UpdateDashboard).
			DELETE("/:index", controllers.DeleteDashboard)
	}
	dashboardRoutes.Use(middleware.IsSysAdm())
	{
		dashboardRoutes.POST("/public", controllers.CreatePublicDashboard)
		dashboardRoutes.GET("/check-index/:index", controllers.CheckDashboardIndex)
	}
}

func configureIssueRoutes() {
	issueRoutes := RouterGroup.Group("/issue")
	issueRoutes.Use(middleware.LimitAPIRequests(global.IssueLimitAPIRequestsTimes, global.LimitRequestsDuration))
	issueRoutes.Use(middleware.LimitTotalRequests(global.IssueLimitTotalRequestsTimes, global.LimitRequestsDuration))
	issueRoutes.Use(middleware.IsLoggedIn())
	{
		issueRoutes.
			POST("/", controllers.CreateIssue)
	}
	issueRoutes.Use(middleware.IsSysAdm())
	{
		issueRoutes.
			GET("/", controllers.GetAllIssues)
		issueRoutes.
			PATCH("/:id", controllers.UpdateIssueByID)
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	wsHandler := websocket.Handler(func(ws *websocket.Conn) {
			defer ws.Close()

			// Add client to the clients map
			clients[ws] = true

			// WebSocket connection established
			for {
					var msg string
					// Read message from client
					err := websocket.Message.Receive(ws, &msg)
					if err != nil {
							// Handle error
							break
					}
					// Print message received from client
					fmt.Printf("Received message: %s\n", msg)

					// Broadcast message to all connected clients
					for client := range clients {
							err := websocket.Message.Send(client, msg)
							if err != nil {
									// Handle error
									break
							}
					}
			}
	})

	// Serve WebSocket requests
	wsHandler.ServeHTTP(w, r)
}