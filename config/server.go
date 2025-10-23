package config

// ServerList :
type ServerList struct {
	Rest Server
}

// Server is struct for server conf
type Server struct {
	TLS       bool `mapstructure:"tls"`
	Name      string
	Host      string
	Port      int
	SecretKey string
	Timeout   int
}

type MessageBroker struct {
	Adapter string
	Host    string
	Port    string
}
