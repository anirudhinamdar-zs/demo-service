package employee

type Employee struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	DOB         string `json:"dob"`
	Major       string `json:"major"`
	City        string `json:"city"`
	Department  string `json:"department"`
}

type NewEmployee struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	DOB         string `json:"dob"`
	Major       string `json:"major"`
	City        string `json:"city"`
	Department  string `json:"department"`
}
