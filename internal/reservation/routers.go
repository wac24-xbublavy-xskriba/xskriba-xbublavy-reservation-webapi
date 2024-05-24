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
    "github.com/gin-gonic/gin"
)

func AddRoutes(engine *gin.Engine) {
  group := engine.Group("")
  
  {
    api := newAmbulanceAPI()
    api.addRoutes(group)
  }
  
  {
    api := newPatientAPI()
    api.addRoutes(group)
  }
  
  {
    api := newReservationAPI()
    api.addRoutes(group)
  }
  
}