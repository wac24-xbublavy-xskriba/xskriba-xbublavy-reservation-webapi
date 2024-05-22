package reservation

func (sex Sex) IsValid() bool {
	switch sex {
	case FEMALE, MALE:
		return true
	default:
		return false
	}
}