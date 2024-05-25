package reservation

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wac24-xbublavy-xskriba/xskriba-xbublavy-reservation-webapi/internal/db_service"
)

// CreatePatient - Create a new patient
func (this *implPatientAPI) CreatePatient(ctx *gin.Context) {
	value, exists := ctx.Get("db_service_patient")
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
  
	db, ok := value.(db_service.DbService[Patient])
	
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

	patient := Patient{}
	err := ctx.BindJSON(&patient)
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

	// Validate the Patient struct
	err = patient.Validate()
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  "Bad Request",
				"message": "Invalid patient data",
				"error":   err.Error(),
			})
		return
	}

	// Create the patient
	if patient.Id == "" {
		patient.Id = uuid.New().String()
	}

	err = db.CreateDocument(ctx, patient.Id, &patient)

	switch err {
	case nil:
		ctx.JSON(
			http.StatusCreated,
			patient,
		)
	case db_service.ErrConflict:
		ctx.JSON(
			http.StatusConflict,
			gin.H{
				"status":  "Conflict",
				"message": "Patient already exists",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to create patient in database",
				"error":   err.Error(),
			},
		)
	}
}

// CreateReservation - Create a new reservation
func (this *implPatientAPI) CreateReservation(ctx *gin.Context) {
	value, exists := ctx.Get("db_service_reservation")
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
  
	db, ok := value.(db_service.DbService[ReservationInput])
	
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

	request := ReservationInput{}
	err := ctx.BindJSON(&request)
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

	reservation := Reservation{}

	// Fetch patient and ambulance from database
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

	// Fetch patient from database
	patientId := ctx.Param("patientId")
	patient, err := patientDB.FindDocument(ctx, patientId)

	if err != nil {
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to fetch patient from database",
				"error":   err.Error(),
			},
		)
		return
	}

	// Fetch ambulance from database
	ambulanceId := request.AmbulanceId
	ambulance, err := ambulanceDB.FindDocument(ctx, ambulanceId)

	if err != nil {
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to fetch ambulance from database",
				"error":   err.Error(),
			},
		)
		return
	}

	// Fill patient and ambulance to reservation
	reservation.Id = uuid.New().String()
	reservation.Patient = *patient
	reservation.Ambulance = *ambulance
	reservation.Start = request.Start
	reservation.End = request.End
	reservation.ExaminationType = request.ExaminationType
	reservation.Message = request.Message

	request.Id = reservation.Id
	request.PatientId = patient.Id

	err = db.CreateDocument(ctx, reservation.Id, &request)

	switch err {
	case nil:
		ctx.JSON(
			http.StatusCreated,
			reservation,
		)
	case db_service.ErrConflict:
		ctx.JSON(
			http.StatusConflict,
			gin.H{
				"status":  "Conflict",
				"message": "Reservation already exists",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to create reservation in database",
				"error":   err.Error(),
			},
		)
	}
}

// DeletePatient - Deletes a patient
func (this *implPatientAPI) DeletePatient(ctx *gin.Context) {
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
	err := db.DeleteDocument(ctx, patientId)
  
	switch err {
	case nil:
		ctx.AbortWithStatus(http.StatusNoContent)
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Patient not found",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to delete patient from database",
				"error":   err.Error(),
			})
	}
}

// GetPatientById - Get a patient by ID
func (this *implPatientAPI) GetPatientById(ctx *gin.Context) {
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
		ctx.JSON(
			http.StatusOK,
			patient,
		)
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Patient not found",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to create patient in database",
				"error":   err.Error(),
			},
		)
	}
}

// GetPatientReservations - Get reservations for a specific patient
func (this *implPatientAPI) GetPatientReservations(ctx *gin.Context) {
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
  
	patientId := ctx.Param("patientId")
	reservationInputs, err := db.GetDocumentsByField(ctx, "patientid", patientId)
  
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "Failed to retrieve reservations from database",
				"error":   err.Error(),
			})
		return
	}
  
	if len(reservationInputs) == 0 {
		reservationInputs = []ReservationInput{}
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

    reservations := make([]Reservation, len(reservationInputs))
    for i, input := range reservationInputs {
        patient, err := patientDB.FindDocument(ctx, input.PatientId)
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

        ambulance, err := ambulanceDB.FindDocument(ctx, input.AmbulanceId)
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

        reservations[i] = Reservation{
            Id:           input.Id,
            Patient:      *patient,
            Ambulance:    *ambulance,
            Start:       input.Start,
            End:         input.End,
            ExaminationType: input.ExaminationType,
            Message:     input.Message,
        }
    }
    
    ctx.JSON(
        http.StatusOK,
        reservations,
    )
}

// GetPatients - Get a list of all patients
func (this *implPatientAPI) GetPatients(ctx *gin.Context) {
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
  
	patients, err := db.GetDocuments(ctx)
  
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "Failed to retrieve patients from database",
				"error":   err.Error(),
			})
		return
	}
  
	if len(patients) == 0 {
		patients = []Patient{}
	}
  
	ctx.JSON(
		http.StatusOK,
		patients,
	)
}

func isDatetimeInsideDate(date time.Time, datetime time.Time) bool {
	year, month, day := date.Date()
	datetimeYear, datetimeMonth, datetimeDay := datetime.Date()

	return year == datetimeYear && month == datetimeMonth && day == datetimeDay
}

// RequestExamination - Request an examination for a specific patient
func (this *implPatientAPI) RequestExamination(ctx *gin.Context) {
	ambulanceValue, ambulanceExists := ctx.Get("db_service_ambulance")

	if !ambulanceExists {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db not found",
				"error":   "db not found",
			})
		return
	}

	ambulanceDB, ambulanceOK := ambulanceValue.(db_service.DbService[Ambulance])

	if !ambulanceOK {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db context is not of required type",
				"error":   "cannot cast db context to db_service.DbService",
			})
		return
	}

	request := RequestExaminationRequest{}
	err := ctx.BindJSON(&request)
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

	ambulances, err := ambulanceDB.GetDocumentsByArrayField(ctx, "medicalexaminations", []string{string(request.ExaminationType)})

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
		ctx.JSON(
			http.StatusOK,
			[]Examination{},
		)
		return
	}

	reservationValue, exists := ctx.Get("db_service_reservation")
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
	
	reservationDB, ok := reservationValue.(db_service.DbService[ReservationInput])
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

	requestDate, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "Failed to parse request date",
				"error":   err.Error(),
			})
		return
	}

	examinations := make([]Examination, 0)

	for _, ambulance := range ambulances {
		reservationInputs, err := reservationDB.GetDocumentsByField(ctx, "ambulanceid", ambulance.Id)

		if err != nil {
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{
					"status":  "Internal Server Error",
					"message": "Failed to retrieve reservations from database",
					"error":   err.Error(),
				})
			return
		}

		closeTime, err := time.Parse("15:04", ambulance.OfficeHours.Close)
		if err != nil {
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{
					"status":  "Internal Server Error",
					"message": "Failed to parse ambulance closing time",
					"error":   err.Error(),
				})
			return
		}
		openTime, err := time.Parse("15:04", ambulance.OfficeHours.Open)
		if err != nil {
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{
					"status":  "Internal Server Error",
					"message": "Failed to parse ambulance opening time",
					"error":   err.Error(),
				})
			return
		}

		duration := closeTime.Sub(openTime)
		interval := time.Minute * 15

		numIntervals := int(duration.Minutes() / interval.Minutes())

		intervals := make([]bool, numIntervals)
		for _, reservationInput := range reservationInputs {
			if !isDatetimeInsideDate(requestDate, reservationInput.Start) {
				continue
			}

			reservationDuration := reservationInput.End.Sub(reservationInput.Start)

			reservationHours, err := time.Parse("15:04", reservationInput.Start.Format("15:04"))
			if err != nil {
				ctx.JSON(
					http.StatusInternalServerError,
					gin.H{
						"status":  "Internal Server Error",
						"message": "Failed to parse reservation start time",
						"error":   err.Error(),
					})
				return
			}

			reservationOffset := reservationHours.Sub(openTime)

			startIndex := int(reservationOffset.Minutes() / interval.Minutes())
			numReservationIntervals := int(reservationDuration.Minutes() / interval.Minutes())

			for i := startIndex; i < startIndex+numReservationIntervals; i++ {
				intervals[i] = true
			}
		}

		for i := 0; i < numIntervals - 3; i++ {
			found := true
			for j := 0; j <= 3; j++ {
				if intervals[i + j] {
					i += j
					found = false
					break
				}
			}

			if found {
				startIndex := i
				endIndex := i + 4

				startTime := requestDate.Add(time.Duration(openTime.Hour()) * time.Hour).Add(time.Duration(openTime.Minute()) * time.Minute).Add(interval * time.Duration(startIndex))
				endTime := requestDate.Add(time.Duration(openTime.Hour()) * time.Hour).Add(time.Duration(openTime.Minute()) * time.Minute).Add(interval * time.Duration(endIndex))

				examinations = append(examinations, Examination{
					Ambulance: ambulance,
					Start: startTime,
					End: endTime,
					ExaminationType: request.ExaminationType,
				})
				
				break
			}
		}
	}

	ctx.JSON(http.StatusOK, examinations)
}

// UpdatePatient - Update an existing patient
func (this *implPatientAPI) UpdatePatient(ctx *gin.Context) {
	updatePatientFunc(ctx, func(c *gin.Context, patient *Patient) (*Patient, interface{}, int) {
		var entry Patient

		if err := c.ShouldBindJSON(&entry); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		// Validate the Patient struct
		var err = entry.Validate()
		if err != nil {
			return nil, gin.H{
				"status":  "Bad Request",
				"message": "Invalid ambulance data",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		if entry.FirstName != "" {
			patient.FirstName = entry.FirstName
		}

		if entry.LastName != "" {
			patient.LastName = entry.LastName
		}

		if entry.Birthday != "" {
			patient.Birthday = entry.Birthday
		}

		if entry.Sex != "" {
			patient.Sex = entry.Sex
		}

		if entry.Bio != "" {
			patient.Bio = entry.Bio
		}

		return patient, patient, http.StatusOK
	})
}

