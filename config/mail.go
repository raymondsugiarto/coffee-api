package config

type MailConfig struct {
	Brevo Brevo `mapstructure:"brevo"`
}

type Brevo struct {
	Host   string `mapstructure:"host"`
	ApiKey string `mapstructure:"apikey"`
	Sender string `mapstructure:"sender"`
}
