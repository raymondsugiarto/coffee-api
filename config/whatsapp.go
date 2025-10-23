package config

type WhatsappConfig struct {
	Saungwa Saungwa
}

type Saungwa struct {
	Url     string
	Appkey  string
	Authkey string
}
