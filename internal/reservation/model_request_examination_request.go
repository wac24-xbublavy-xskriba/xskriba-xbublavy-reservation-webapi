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
	"time"
)

type RequestExaminationRequest struct {

	Date time.Time `json:"date"`

	ExaminationType MedicalExaminations `json:"examinationType"`
}
