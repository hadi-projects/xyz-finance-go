package dto

type ConsumerResponse struct {
	NIK          string  `json:"nik"`
	FullName     string  `json:"full_name"`
	LegalName    string  `json:"legal_name"`
	PlaceOfBirth string  `json:"place_of_birth"`
	DateOfBirth  string  `json:"date_of_birth"`
	Salary       float64 `json:"salary"`
	KTPImage     string  `json:"ktp_image"`
	SelfieImage  string  `json:"selfie_image"`
}

type UserProfileResponse struct {
	UserID   uint              `json:"user_id"`
	Email    string            `json:"email"`
	Consumer *ConsumerResponse `json:"consumer"`
}
