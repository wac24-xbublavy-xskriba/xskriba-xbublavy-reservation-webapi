/*
 * Reservation Api
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Contact: xbublavy@stuba.sk
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

 package reservation

import (
   "net/http"

   "github.com/gin-gonic/gin"
)

type AmbulanceAPI interface {

   // internal registration of api routes
   addRoutes(routerGroup *gin.RouterGroup)

    // CreateAmbulance - Create a new ambulance
   CreateAmbulance(ctx *gin.Context)

    // DeleteAmbulance - Deletes an ambulance
   DeleteAmbulance(ctx *gin.Context)

    // GetAmbulanceById - Get an ambulance by ID
   GetAmbulanceById(ctx *gin.Context)

    // GetAmbulanceReservationsById - Get reservations for a specific ambulance
   GetAmbulanceReservationsById(ctx *gin.Context)

    // GetAmbulances - Get a list of all ambulances
   GetAmbulances(ctx *gin.Context)

    // UpdateAmbulance - Update an existing ambulance
   UpdateAmbulance(ctx *gin.Context)

 }

 // partial implementation of AmbulanceAPI - all functions must be implemented in add on files
type implAmbulanceAPI struct {

}

func newAmbulanceAPI() AmbulanceAPI {
  return &implAmbulanceAPI{}
}

func (this *implAmbulanceAPI) addRoutes(routerGroup *gin.RouterGroup) {
  routerGroup.Handle( http.MethodPost, "/ambulances", this.CreateAmbulance)
  routerGroup.Handle( http.MethodDelete, "/ambulances/:ambulanceId", this.DeleteAmbulance)
  routerGroup.Handle( http.MethodGet, "/ambulances/:ambulanceId", this.GetAmbulanceById)
  routerGroup.Handle( http.MethodGet, "/ambulances/:ambulanceId/reservations", this.GetAmbulanceReservationsById)
  routerGroup.Handle( http.MethodGet, "/ambulances", this.GetAmbulances)
  routerGroup.Handle( http.MethodPut, "/ambulances/:ambulanceId", this.UpdateAmbulance)
}

// Copy following section to separate file, uncomment, and implement accordingly
// // CreateAmbulance - Create a new ambulance
// func (this *implAmbulanceAPI) CreateAmbulance(ctx *gin.Context) {
//  	ctx.AbortWithStatus(http.StatusNotImplemented)
// }
//
// // DeleteAmbulance - Deletes an ambulance
// func (this *implAmbulanceAPI) DeleteAmbulance(ctx *gin.Context) {
//  	ctx.AbortWithStatus(http.StatusNotImplemented)
// }
//
// // GetAmbulanceById - Get an ambulance by ID
// func (this *implAmbulanceAPI) GetAmbulanceById(ctx *gin.Context) {
//  	ctx.AbortWithStatus(http.StatusNotImplemented)
// }
//
// // GetAmbulanceReservationsById - Get reservations for a specific ambulance
// func (this *implAmbulanceAPI) GetAmbulanceReservationsById(ctx *gin.Context) {
//  	ctx.AbortWithStatus(http.StatusNotImplemented)
// }
//
// // GetAmbulances - Get a list of all ambulances
// func (this *implAmbulanceAPI) GetAmbulances(ctx *gin.Context) {
//  	ctx.AbortWithStatus(http.StatusNotImplemented)
// }
//
// // UpdateAmbulance - Update an existing ambulance
// func (this *implAmbulanceAPI) UpdateAmbulance(ctx *gin.Context) {
//  	ctx.AbortWithStatus(http.StatusNotImplemented)
// }
//

