package reservation

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wac24-xbublavy-xskriba/xskriba-xbublavy-reservation-webapi/internal/db_service"
)

// Validate checks if the Reservation struct is valid
func (reservation *Reservation) Validate() error {
	currentTime := time.Now()

	if reservation.Patient.Validate() != nil {
		return fmt.Errorf("Invalid patient data")
	}

	if reservation.Ambulance.Validate() != nil {
		return fmt.Errorf("Invalid ambulance data")
	}

	if reservation.Start.Before(currentTime) {
		return fmt.Errorf("start time must be in the future")
	}

	if reservation.End.Before(currentTime) {
		return fmt.Errorf("end time must be in the future")
	}

	if reservation.Start.After(reservation.End) {
		return fmt.Errorf("start time must be before end time")
	}

	if !reservation.ExaminationType.IsValid() {
		return fmt.Errorf("Invalid examination type")
	}

	if len(reservation.Message) > 200 {
		return fmt.Errorf("Message exceeds maximum length of 200 characters")
	}

	return nil
}

type reservationUpdater = func(
    ctx *gin.Context,
    reservationInput *ReservationInput,
) (updatedReservation *ReservationInput, responseContent interface{}, status int)

func updateReservationFunc(ctx *gin.Context, updater reservationUpdater) {
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

    reservation, err := db.FindDocument(ctx, reservationId)

    switch err {
    case nil:
        // continue
    case db_service.ErrNotFound:
        ctx.JSON(
            http.StatusNotFound,
            gin.H{
                "status":  "Not Found",
                "message": "Reservation not found",
                "error":   err.Error(),
            },
        )
        return
    default:
        ctx.JSON(
            http.StatusBadGateway,
            gin.H{
                "status":  "Bad Gateway",
                "message": "Failed to load reservation from database",
                "error":   err.Error(),
            })
        return
    }

    if !ok {
        ctx.JSON(
            http.StatusInternalServerError,
            gin.H{
                "status":  "Internal Server Error",
                "message": "Failed to cast reservation from database",
                "error":   "Failed to cast reservation from database",
            })
        return
    }

    updatedReservation, responseObject, status := updater(ctx, reservation)

    if updatedReservation != nil {
        err = db.UpdateDocument(ctx, reservationId, updatedReservation)
    } else {
        err = nil // redundant but for clarity
    }

    switch err {
    case nil:
        if responseObject != nil {
            ctx.JSON(status, responseObject)
        } else {
            ctx.AbortWithStatus(status)
        }
    case db_service.ErrNotFound:
        ctx.JSON(
            http.StatusNotFound,
            gin.H{
                "status":  "Not Found",
                "message": "Reservation was deleted while processing the request",
                "error":   err.Error(),
            },
        )
    default:
        ctx.JSON(
            http.StatusBadGateway,
            gin.H{
                "status":  "Bad Gateway",
                "message": "Failed to update reservation in database",
                "error":   err.Error(),
            })
    }

}