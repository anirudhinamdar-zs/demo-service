package department

var AllowedCodes = map[string]struct{}{
	"CSE": {},
	"IT":  {},
	"ECE": {},
	"EEE": {},
	"ME":  {},
}

func IsValidCode(code string) bool {
	_, ok := AllowedCodes[code]
	return ok
}
