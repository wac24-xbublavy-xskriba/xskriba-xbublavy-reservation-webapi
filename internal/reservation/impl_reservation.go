package reservation

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wac24-xbublavy-xskriba/xskriba-xbublavy-reservation-webapi/internal/db_service"
)

// DeleteReservation - Deletes a reservation
func (this *implReservationAPI) DeleteReservation(ctx *gin.Context) {
	value, exists := ctx.Get("db_service_reservation")
	if !exists {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service not found",
				"error":   "db_service not found",
			})
		return
	}
  
	db, ok := value.(db_service.DbService[ReservationInput])
	if !ok {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service context is not of type db_service.DbService",
				"error":   "cannot cast db_service context to db_service.DbService",
			})
		return
	}
  
	reservationId := ctx.Param("reservationId")
	err := db.DeleteDocument(ctx, reservationId)
  
	switch err {
	case nil:
		ctx.AbortWithStatus(http.StatusNoContent)
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Reservation not found",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to delete reservation from database",
				"error":   err.Error(),
			})
	}
}

// GetReservationById - Get a reservation by ID
func (this *implReservationAPI) GetReservationById(ctx *gin.Context) {
	value, exists := ctx.Get("db_service_reservation")
	if !exists {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service not found",
				"error":   "db_service not found",
			})
		return
	}
  
	db, ok := value.(db_service.DbService[ReservationInput])
	if !ok {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service context is not of type db_service.DbService",
				"error":   "cannot cast db_service context to db_service.DbService",
			})
		return
	}
  
	reservationId := ctx.Param("reservationId")
	reservationInput, err := db.FindDocument(ctx, reservationId)

	switch err {
	case nil:
		
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Ambulance not found",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to create ambulance in database",
				"error":   err.Error(),
			},
		)
	}

	patientValue, patientExists := ctx.Get("db_service_patient")
	ambulanceValue, ambulanceExists := ctx.Get("db_service_ambulance")

	if !patientExists || !ambulanceExists {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db not found",
				"error":   "db not found",
			})
		return
	}

	patientDB, patientOK := patientValue.(db_service.DbService[Patient])
	ambulanceDB, ambulanceOK := ambulanceValue.(db_service.DbService[Ambulance])

	if !patientOK || !ambulanceOK {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db context is not of required type",
				"error":   "cannot cast db context to db_service.DbService",
			})
		return
	}

	patient, err := patientDB.FindDocument(ctx, reservationInput.PatientId)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "Failed to retrieve patient from database",
				"error":   err.Error(),
			})
		return
	}

	ambulance, err := ambulanceDB.FindDocument(ctx, reservationInput.AmbulanceId)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "Failed to retrieve ambulance from database",
				"error":   err.Error(),
			})
		return
	}

	reservation := Reservation{
		Id: reservationInput.Id,
		Patient: *patient,
		Ambulance: *ambulance,
		Start: reservationInput.Start,
		End: reservationInput.End,
		ExaminationType: reservationInput.ExaminationType,
		Message: reservationInput.Message,
	}

	ctx.JSON(
		http.StatusOK,
		reservation,
	)
}

// UpdateReservation - Update an existing reservation
func (this *implReservationAPI) UpdateReservation(ctx *gin.Context) {
	updateReservationFunc(ctx, func(c *gin.Context, reservationInput *ReservationInput) (updatedReservation *ReservationInput, responseContent interface{}, status int) {
		var entry ReservationInput

		if err := c.ShouldBindJSON(&entry); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		reservationInput.Message = entry.Message

		patientValue, patientExists := ctx.Get("db_service_patient")
		ambulanceValue, ambulanceExists := ctx.Get("db_service_ambulance")

		if !patientExists || !ambulanceExists {
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{
					"status":  "Internal Server Error",
					"message": "db not found",
					"error":   "db not found",
				})
			return
		}

		patientDB, patientOK := patientValue.(db_service.DbService[Patient])
		ambulanceDB, ambulanceOK := ambulanceValue.(db_service.DbService[Ambulance])

		if !patientOK || !ambulanceOK {
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{
					"status":  "Internal Server Error",
					"message": "db context is not of required type",
					"error":   "cannot cast db context to db_service.DbService",
				})
			return
		}

		patient, err := patientDB.FindDocument(ctx, reservationInput.PatientId)
		if err != nil {
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{
					"status":  "Internal Server Error",
					"message": "Failed to retrieve patient from database",
					"error":   err.Error(),
				})
			return
		}

		ambulance, err := ambulanceDB.FindDocument(ctx, reservationInput.AmbulanceId)
		if err != nil {
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{
					"status":  "Internal Server Error",
					"message": "Failed to retrieve ambulance from database",
					"error":   err.Error(),
				})
			return
		}

		reservation := Reservation{
			Id: reservationInput.Id,
			Patient: *patient,
			Ambulance: *ambulance,
			Start: reservationInput.Start,
			End: reservationInput.End,
			ExaminationType: reservationInput.ExaminationType,
			Message: reservationInput.Message,
		}

		return reservationInput, reservation, http.StatusOK
	})
}
