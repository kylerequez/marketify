package shared

type ServerConfig struct {
	AppName  string
	Hostname string
	Port     string
}

type SQLConfig struct {
	URI string
}
