package reservation

// IsValid checks if the given medical examination is valid
func (examination MedicalExaminations) IsValid() bool {
	switch examination {
	case X_RAY, MRI, CT, ULTRASOUND, BLOOD_TEST:
		return true
	default:
		return false
	}
}