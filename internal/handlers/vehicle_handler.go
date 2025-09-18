package handlers

import (
	"net/http"

	"github.com/asaaitika/fleetmgm-tst/internal/repositories"
	"github.com/gin-gonic/gin"
)

// VehicleHandler handles HTTP requests
type VehicleHandler struct {
	repo *repositories.VehicleRepository
}

// NewVehicleHandler creates new handler instance
func NewVehicleHandler(repo *repositories.VehicleRepository) *VehicleHandler {
	return &VehicleHandler{repo: repo}
}

// GetLastLocation endpoint: GET /vehicles/{vehicle_id}/location
func (h *VehicleHandler) GetLastLocation(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")

	if vehicleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "vehicle_id is required",
		})
		return
	}

	location, err := h.repo.GetLastLocation(vehicleID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Vehicle not found or no location data",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get location",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"vehicle_id": location.VehicleID,
		"latitude":   location.Latitude,
		"longitude":  location.Longitude,
		"timestamp":  location.Timestamp,
	})
}

// GetLocationHistory endpoint: GET /vehicles/{vehicle_id}/history?start=xxx&end=xxx
func (h *VehicleHandler) GetLocationHistory(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")

	var request struct {
		Start int64 `form:"start" binding:"required"`
		End   int64 `form:"end" binding:"required"`
	}

	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "start and end timestamps are required",
		})
		return
	}

	if request.Start > request.End {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "start time must be before end time",
		})
		return
	}

	locations, err := h.repo.GetLocationHistory(vehicleID, request.Start, request.End)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get history",
		})
		return
	}

	var response []gin.H
	for _, loc := range locations {
		response = append(response, gin.H{
			"vehicle_id": loc.VehicleID,
			"latitude":   loc.Latitude,
			"longitude":  loc.Longitude,
			"timestamp":  loc.Timestamp,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"vehicle_id": vehicleID,
		"count":      len(response),
		"history":    response,
	})
}
