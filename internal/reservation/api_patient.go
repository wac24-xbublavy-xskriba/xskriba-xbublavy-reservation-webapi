/*
 * Reservation API
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

 package reservation

import (
   "net/http"

   "github.com/gin-gonic/gin"
)

type PatientAPI interface {

   // internal registration of api routes
   addRoutes(routerGroup *gin.RouterGroup)

    // CreatePatient - Create a new patient
   CreatePatient(ctx *gin.Context)

    // DeletePatient - Deletes a patient
   DeletePatient(ctx *gin.Context)

    // GetPatientById - Get a patient by ID
   GetPatientById(ctx *gin.Context)

    // GetPatientReservations - Get reservations for a specific patient
   GetPatientReservations(ctx *gin.Context)

    // GetPatients - Get a list of all patients
   GetPatients(ctx *gin.Context)

    // RequestExamination - Request an examination for a specific patient
   RequestExamination(ctx *gin.Context)

    // UpdatePatient - Update an existing patient
   UpdatePatient(ctx *gin.Context)

 }

 // partial implementation of PatientAPI - all functions must be implemented in add on files
type implPatientAPI struct {

}

func newPatientAPI() PatientAPI {
  return &implPatientAPI{}
}

func (this *implPatientAPI) addRoutes(routerGroup *gin.RouterGroup) {
  routerGroup.Handle( http.MethodPost, "/patients", this.CreatePatient)
  routerGroup.Handle( http.MethodDelete, "/patients/:patientId", this.DeletePatient)
  routerGroup.Handle( http.MethodGet, "/patients/:patientId", this.GetPatientById)
  routerGroup.Handle( http.MethodGet, "/patients/:patientId/reservations", this.GetPatientReservations)
  routerGroup.Handle( http.MethodGet, "/patients", this.GetPatients)
  routerGroup.Handle( http.MethodPost, "/patients/:patientId/request-examination", this.RequestExamination)
  routerGroup.Handle( http.MethodPut, "/patients/:patientId", this.UpdatePatient)
}

// Copy following section to separate file, uncomment, and implement accordingly
// // CreatePatient - Create a new patient
// func (this *implPatientAPI) CreatePatient(ctx *gin.Context) {
//  	ctx.AbortWithStatus(http.StatusNotImplemented)
// }
//
// // DeletePatient - Deletes a patient
// func (this *implPatientAPI) DeletePatient(ctx *gin.Context) {
//  	ctx.AbortWithStatus(http.StatusNotImplemented)
// }
//
// // GetPatientById - Get a patient by ID
// func (this *implPatientAPI) GetPatientById(ctx *gin.Context) {
//  	ctx.AbortWithStatus(http.StatusNotImplemented)
// }
//
// // GetPatientReservations - Get reservations for a specific patient
// func (this *implPatientAPI) GetPatientReservations(ctx *gin.Context) {
//  	ctx.AbortWithStatus(http.StatusNotImplemented)
// }
//
// // GetPatients - Get a list of all patients
// func (this *implPatientAPI) GetPatients(ctx *gin.Context) {
//  	ctx.AbortWithStatus(http.StatusNotImplemented)
// }
//
// // RequestExamination - Request an examination for a specific patient
// func (this *implPatientAPI) RequestExamination(ctx *gin.Context) {
//  	ctx.AbortWithStatus(http.StatusNotImplemented)
// }
//
// // UpdatePatient - Update an existing patient
// func (this *implPatientAPI) UpdatePatient(ctx *gin.Context) {
//  	ctx.AbortWithStatus(http.StatusNotImplemented)
// }
//

