package shared

type LoginFormData struct {
	Email    string
	Password string
	Errors   map[string]string
}
