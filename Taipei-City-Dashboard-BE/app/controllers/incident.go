package controllers

import (
	"net/http"

	"TaipeiCityDashboardBE/app/models"

	"github.com/gin-gonic/gin"
)

func GetIncident(c *gin.Context) {
	type incidentQuery struct {
		PageSize       int    `form:"pagesize"`
		PageNum        int    `form:"pagenum"`
		Sort           string `form:"sort"`
		Order          string `form:"order"`
	}

	// Get query parameters
	var query incidentQuery
	c.ShouldBindQuery(&query)

	incidents, totalIncidents, resultNum, err := models.GetAllIncident(query.PageSize, query.PageNum, query.Sort, query.Order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "total": totalIncidents, "results": resultNum, "data": incidents})
}

func CreateIncident(c *gin.Context) {
	var incident models.Incident
	// var buf bytes.Buffer
	// _, err := io.Copy(&buf, c.Request.Body)
	// if err != nil {
	// 		// Handle error
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
	// 		return
	// }
	
	// Convert buffer to string
	// requestBody := buf.String()
	// fmt.Println(requestBody)

	// Bind the issue data
	if err := c.ShouldBindJSON(&incident); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if incident.Type == "" || incident.Description == "" || incident.Distance == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "title, description, distance info is required"})
		return
	}

	tmpIncident, err := models.CreateIncident(incident.Type, incident.Description, incident.Distance, incident.Latitude, incident.Longitude, incident.Time)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": tmpIncident})
}

func DeleteIncident(c *gin.Context) {
	var incident models.Incident

	// Bind the issue data
	if err := c.ShouldBindJSON(&incident); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// if incident.Type == "" || incident.Description == "" || incident.Distance == 0 {
	// 	c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "title, description, distance info is required"})
	// 	return
	// }

	tmpIncident, err := models.DeleteIncident(incident.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": tmpIncident})
}


func CreateIncidentType(c *gin.Context) {
	var incidentType models.IncidentType

	if err := c.ShouldBindJSON(&incidentType); (err != nil) || (incidentType.Type == "") {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	tmpIncident, err := models.CreateIncidentType(incidentType.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	if tmpIncident.IsEmpty() {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Incident type " + incidentType.Type + " already exists!!"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": tmpIncident})
}

func UpdateIncidentType(c *gin.Context) {
	var incident models.Incident

	if err := c.ShouldBindJSON(&incident); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	tmpIncident, err := models.UpdateIncidentType(incident.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": tmpIncident})
}
