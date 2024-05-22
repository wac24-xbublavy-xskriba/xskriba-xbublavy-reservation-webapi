package reservation

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wac24-xbublavy-xskriba/xskriba-xbublavy-reservation-webapi/internal/db_service"
)

// CreateAmbulance - Create a new ambulance
func (this *implAmbulanceAPI) CreateAmbulance(ctx *gin.Context) {
  value, exists := ctx.Get("db_service_ambulance")
  if !exists {
      ctx.JSON(
          http.StatusInternalServerError,
          gin.H{
              "status":  "Internal Server Error",
              "message": "db not found",
              "error":   "db not found",
          })
      return
  }

  db, ok := value.(db_service.DbService[Ambulance])
  if !ok {
      ctx.JSON(
          http.StatusInternalServerError,
          gin.H{
              "status":  "Internal Server Error",
              "message": "db context is not of required type",
              "error":   "cannot cast db context to db_service.DbService",
          })
      return
  }

  ambulance := Ambulance{}
  err := ctx.BindJSON(&ambulance)
  if err != nil {
      ctx.JSON(
          http.StatusBadRequest,
          gin.H{
              "status":  "Bad Request",
              "message": "Invalid request body",
              "error":   err.Error(),
          })
      return
  }

  // Validate the Ambulance struct
  err = ambulance.Validate()
  if err != nil {
      ctx.JSON(
          http.StatusBadRequest,
          gin.H{
              "status":  "Bad Request",
              "message": "Invalid ambulance data",
              "error":   err.Error(),
          })
      return
  }

  if ambulance.Id == "" {
      ambulance.Id = uuid.New().String()
  }

  err = db.CreateDocument(ctx, ambulance.Id, &ambulance)

  switch err {
  case nil:
      ctx.JSON(
          http.StatusCreated,
          ambulance,
      )
  case db_service.ErrConflict:
      ctx.JSON(
          http.StatusConflict,
          gin.H{
              "status":  "Conflict",
              "message": "Ambulance already exists",
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
}

// DeleteAmbulance - Deletes an ambulance
func (this *implAmbulanceAPI) DeleteAmbulance(ctx *gin.Context) {
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
  err := db.DeleteDocument(ctx, ambulanceId)

  switch err {
  case nil:
      ctx.AbortWithStatus(http.StatusNoContent)
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
              "message": "Failed to delete ambulance from database",
              "error":   err.Error(),
          })
  }
}

// GetAmbulanceById - Get an ambulance by ID
func (this *implAmbulanceAPI) GetAmbulanceById(ctx *gin.Context) {
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
      ctx.JSON(
          http.StatusOK,
          ambulance,
      )
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

  // TODO: add reservations to the response
}

// GetAmbulances - Get a list of all ambulances
func (this *implAmbulanceAPI) GetAmbulances(ctx *gin.Context) {
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

  ambulances, err := db.GetDocuments(ctx)

  if err != nil {
      ctx.JSON(
          http.StatusInternalServerError,
          gin.H{
              "status":  "Internal Server Error",
              "message": "Failed to retrieve ambulances from database",
              "error":   err.Error(),
          })
      return
  }

  if len(ambulances) == 0 {
    ambulances = []Ambulance{}
  }

  ctx.JSON(
      http.StatusOK,
      ambulances,
  )
}

// UpdateAmbulance - Update an existing ambulance
func (this *implAmbulanceAPI) UpdateAmbulance(ctx *gin.Context) {
  updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
    var entry Ambulance

    if err := c.ShouldBindJSON(&entry); err != nil {
        return nil, gin.H{
            "status":  http.StatusBadRequest,
            "message": "Invalid request body",
            "error":   err.Error(),
        }, http.StatusBadRequest
    }

    // Validate the Ambulance struct
    var err = entry.Validate()
    if err != nil {
        return nil, gin.H{
                "status":  "Bad Request",
                "message": "Invalid ambulance data",
                "error":   err.Error(),
        }, http.StatusBadRequest
    }

    ambulanceId := ctx.Param("ambulanceId")

    if ambulanceId == "" {
        return nil, gin.H{
            "status":  http.StatusBadRequest,
            "message": "Ambulance ID is required",
        }, http.StatusBadRequest
    }

    fmt.Printf(entry.Id, entry.Address, entry.MedicalExaminations, entry.Name, entry.OfficeHours.Close, entry.OfficeHours.Open)

    fmt.Printf(ambulance.Id, ambulance.Address, ambulance.MedicalExaminations, ambulance.Name, ambulance.OfficeHours.Close, ambulance.OfficeHours.Open)

    if entry.Name != "" {
      ambulance.Name = entry.Name
    }

    if entry.Address != "" {
      ambulance.Address = entry.Address
    }

    if entry.OfficeHours.Open != "" {
      ambulance.OfficeHours.Open = entry.OfficeHours.Open
    }

    if entry.OfficeHours.Close != "" {
      ambulance.OfficeHours.Close = entry.OfficeHours.Close
    }

    if entry.MedicalExaminations != nil {
      ambulance.MedicalExaminations = entry.MedicalExaminations
    }

    return ambulance, ambulance, http.StatusOK
  })
}
