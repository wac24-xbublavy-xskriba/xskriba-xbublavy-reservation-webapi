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

type MedicalExaminations string

// List of MedicalExaminations
const (
	X_RAY MedicalExaminations = "X-ray"
	MRI MedicalExaminations = "MRI"
	CT MedicalExaminations = "CT"
	ULTRASOUND MedicalExaminations = "Ultrasound"
	BLOOD_TEST MedicalExaminations = "Blood Test"
)
