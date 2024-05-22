package reservation

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateReservation - Create a new reservation
func (this *implReservationAPI) CreateReservation(ctx *gin.Context) {
 	ctx.AbortWithStatus(http.StatusNotImplemented)
}

// DeleteReservation - Deletes a reservation
func (this *implReservationAPI) DeleteReservation(ctx *gin.Context) {
 	ctx.AbortWithStatus(http.StatusNotImplemented)
}

// GetAmbulanceReservationsById - Get an ambulance by ID along with the list of medical examination reservations
func (this *implReservationAPI) GetAmbulanceReservationsById(ctx *gin.Context) {
 	ctx.AbortWithStatus(http.StatusNotImplemented)
}

// GetReservationById - Get a reservation by ID
func (this *implReservationAPI) GetReservationById(ctx *gin.Context) {
 	ctx.AbortWithStatus(http.StatusNotImplemented)
}

// UpdateReservation - Update an existing reservation
func (this *implReservationAPI) UpdateReservation(ctx *gin.Context) {
 	ctx.AbortWithStatus(http.StatusNotImplemented)
}
