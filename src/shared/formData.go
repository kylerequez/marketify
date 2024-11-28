package shared

type LoginFormData struct {
	Email    string
	Password string
	Errors   map[string]string
}

type SignupFormData struct {
	Firstname  string
	Middlename string
	Lastname   string
	Age        uint
	Gender     string
	Email      string
	Password   string
	RePassword string
	Errors     map[string]string
}
