package reservation

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wac24-xbublavy-xskriba/xskriba-xbublavy-reservation-webapi/internal/db_service"
)

type ambulanceUpdater = func(
    ctx *gin.Context,
    ambulance *Ambulance,
) (updatedAmbulance *Ambulance, responseContent interface{}, status int)

func updateAmbulanceFunc(ctx *gin.Context, updater ambulanceUpdater) {
    value, exists := ctx.Get("db_service_ambulance")
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

    db, ok := value.(db_service.DbService[Ambulance])
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

    ambulanceId := ctx.Param("ambulanceId")

    ambulance, err := db.FindDocument(ctx, ambulanceId)

    switch err {
    case nil:
        // continue
    case db_service.ErrNotFound:
        ctx.JSON(
            http.StatusNotFound,
            gin.H{
                "status":  "Not Found",
                "message": "Ambulance not found",
                "error":   err.Error(),
            },
        )
        return
    default:
        ctx.JSON(
            http.StatusBadGateway,
            gin.H{
                "status":  "Bad Gateway",
                "message": "Failed to load ambulance from database",
                "error":   err.Error(),
            })
        return
    }

    if !ok {
        ctx.JSON(
            http.StatusInternalServerError,
            gin.H{
                "status":  "Internal Server Error",
                "message": "Failed to cast ambulance from database",
                "error":   "Failed to cast ambulance from database",
            })
        return
    }

    updatedAmbulance, responseObject, status := updater(ctx, ambulance)

    if updatedAmbulance != nil {
        err = db.UpdateDocument(ctx, ambulanceId, updatedAmbulance)
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

// CheckDuplicates checks if there are duplicate values in the MedicalExaminations slice
func CheckMedicalExaminationDuplicates(examinations []MedicalExaminations) (bool, []MedicalExaminations) {
	examinationsMap := make(map[MedicalExaminations]bool)
	var duplicates []MedicalExaminations
	for _, exam := range examinations {
		if _, exists := examinationsMap[exam]; exists {
			duplicates = append(duplicates, exam)
		}
		examinationsMap[exam] = true
	}
	return len(duplicates) > 0, duplicates
}

// ValidateMedicalExaminations checks if all medical examinations are valid
func ValidateMedicalExaminations(examinations []MedicalExaminations) (bool, []MedicalExaminations) {
	var incorrectExams []MedicalExaminations
    for _, exam := range examinations {
        if !exam.IsValid() {
            incorrectExams = append(incorrectExams, exam)
        }
    }

    return len(incorrectExams) == 0, incorrectExams
}

// Validate checks if the Ambulance struct is valid
func (a *Ambulance) Validate() error {
    // Check if OfficeHours are valid
    if !a.OfficeHours.IsValid() {
        return fmt.Errorf("Invalid office hours. Open time must be before close time and both times must be in the future.")
    }

    // Check if the name is empty
    if len(a.Name) == 0 {
        return fmt.Errorf("Name is required.")
    }

    // Check if the name is at least one character long and max 50 characters long
    if len(a.Name) > 50 {
        return fmt.Errorf("Invalid name. Name must be at least one character long and max 50 characters long.")
    }

    // Check if the address is empty
    if len(a.Address) == 0 {
        return fmt.Errorf("Address is required.")
    }

    // Check if the address is at least one character long and max 50 characters long
    if len(a.Address) > 50 {
        return fmt.Errorf("Invalid address. Address must be at least one character long and max 50 characters long.")
    }

    // Check if the examinations are not empty
    if len(a.MedicalExaminations) == 0 {
        return fmt.Errorf("The medical examinations are empty")
    }

    // Check if MedicalExaminations are valid
    validExams, incorrectExams := ValidateMedicalExaminations(a.MedicalExaminations)
    if !validExams {
        return fmt.Errorf("Invalid medical examinations: %v", incorrectExams)
    }

    // Check if MedicalExaminations contain duplicates
    duplicateExams, duplicates := CheckMedicalExaminationDuplicates(a.MedicalExaminations)

    if duplicateExams {
        return fmt.Errorf("The medical examinations contain duplicate values: %v", duplicates)
    }

    return nil
}

func (a *AmbulanceInput) Validate() error {
    // Check if OfficeHours are valid
    if !a.OfficeHours.IsValid() {
        return fmt.Errorf("Invalid office hours. Open time must be before close time and both times must be in the future.")
    }

    // Check if the name is empty
    if len(a.Name) == 0 {
        return fmt.Errorf("Name is required.")
    }

    // Check if the name is at least one character long and max 50 characters long
    if len(a.Name) > 50 {
        return fmt.Errorf("Invalid name. Name must be at least one character long and max 50 characters long.")
    }

    // Check if the address is empty
    if len(a.Address) == 0 {
        return fmt.Errorf("Address is required.")
    }

    // Check if the address is at least one character long and max 50 characters long
    if len(a.Address) > 50 {
        return fmt.Errorf("Invalid address. Address must be at least one character long and max 50 characters long.")
    }

    // Check if MedicalExaminations are valid
    validExams, incorrectExams := ValidateMedicalExaminations(a.MedicalExaminations)
    if !validExams {
        return fmt.Errorf("Invalid medical examinations: %v", incorrectExams)
    }

    // Check if MedicalExaminations contain duplicates
    duplicateExams, duplicates := CheckMedicalExaminationDuplicates(a.MedicalExaminations)
    if duplicateExams {
        return fmt.Errorf("The medical examinations contain duplicate values: %v", duplicates)
    }

    return nil
}