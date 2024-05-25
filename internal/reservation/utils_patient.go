package reservation

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wac24-xbublavy-xskriba/xskriba-xbublavy-reservation-webapi/internal/db_service"
)

type patientUpdater = func(
    ctx *gin.Context,
    patient *Patient,
) (updatedPatient *Patient, responseContent interface{}, status int)

func updatePatientFunc(ctx *gin.Context, updater patientUpdater) {
    value, exists := ctx.Get("db_service_patient")
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

    db, ok := value.(db_service.DbService[Patient])
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

    patientId := ctx.Param("patientId")

    patient, err := db.FindDocument(ctx, patientId)

    switch err {
    case nil:
        // continue
    case db_service.ErrNotFound:
        ctx.JSON(
            http.StatusNotFound,
            gin.H{
                "status":  "Not Found",
                "message": "Patient not found",
                "error":   err.Error(),
            },
        )
        return
    default:
        ctx.JSON(
            http.StatusBadGateway,
            gin.H{
                "status":  "Bad Gateway",
                "message": "Failed to load patient from database",
                "error":   err.Error(),
            })
        return
    }

    if !ok {
        ctx.JSON(
            http.StatusInternalServerError,
            gin.H{
                "status":  "Internal Server Error",
                "message": "Failed to cast patient from database",
                "error":   "Failed to cast patient from database",
            })
        return
    }

    updatedPatient, responseObject, status := updater(ctx, patient)

    if updatedPatient != nil {
        err = db.UpdateDocument(ctx, patientId, updatedPatient)
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
                "message": "Ambulance was deleted while processing the request",
                "error":   err.Error(),
            },
        )
    default:
        ctx.JSON(
            http.StatusBadGateway,
            gin.H{
                "status":  "Bad Gateway",
                "message": "Failed to update ambulance in database",
                "error":   err.Error(),
            })
    }

}

func (patient *Patient) Validate() error {
    if !patient.Sex.IsValid() {
        return fmt.Errorf("Invalid sex")
    }
    
    if len(patient.FirstName) == 0 {
        return fmt.Errorf("First name is required")
    }
    
    if len(patient.FirstName) > 20 {
        return fmt.Errorf("First name exceeds maximum length of 20 characters")
    }
    
    if len(patient.LastName) == 0 {
        return fmt.Errorf("Last name is required")
    }
    
    if len(patient.LastName) > 20 {
        return fmt.Errorf("Last name exceeds maximum length of 20 characters")
    }
    
    if len(patient.Bio) > 200 {
        return fmt.Errorf("Bio exceeds maximum length of 200 characters")
    }

    birthday, err := time.Parse("2006-01-02", patient.Birthday)
    if err != nil {
        return fmt.Errorf("Failed to parse birthday: %v", err)
    }
    
    if birthday.After(time.Now()) {
        return fmt.Errorf("Birthday cannot be in the future")
    }
    
    return nil
}

func (patient *PatientInput) Validate() error {
    if !patient.Sex.IsValid() {
        return fmt.Errorf("Invalid sex")
    }
    
    if len(patient.FirstName) == 0 {
        return fmt.Errorf("First name is required")
    }
    
    if len(patient.FirstName) > 20 {
        return fmt.Errorf("First name exceeds maximum length of 20 characters")
    }
    
    if len(patient.LastName) == 0 {
        return fmt.Errorf("Last name is required")
    }
    
    if len(patient.LastName) > 20 {
        return fmt.Errorf("Last name exceeds maximum length of 20 characters")
    }
    
    if len(patient.Bio) > 200 {
        return fmt.Errorf("Bio exceeds maximum length of 200 characters")
    }
    
    birthday, err := time.Parse("2006-01-02", patient.Birthday)
    if err != nil {
        return fmt.Errorf("Failed to parse birthday: %v", err)
    }
    
    if birthday.After(time.Now()) {
        return fmt.Errorf("Birthday cannot be in the future")
    }
    
    return nil
}