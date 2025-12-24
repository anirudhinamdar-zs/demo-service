package department

type Department struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Floor       int    `json:"floor"`
	Description string `json:"description"`
}

type NewDepartment struct {
	Name        string `json:"name"`
	Floor       int    `json:"floor"`
	Description string `json:"description"`
}
