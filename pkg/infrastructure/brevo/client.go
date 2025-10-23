package brevo

import (
	b "github.com/getbrevo/brevo-go/lib"
	"github.com/raymondsugiarto/coffee-api/config"
)

func NewClient() *b.APIClient {
	c := config.GetConfig()
	cfg := b.NewConfiguration()
	//Configure API key authorization: api-key
	cfg.AddDefaultHeader("api-key", c.Mail.Brevo.ApiKey)
	//Configure API key authorization: partner-key
	cfg.AddDefaultHeader("partner-key", c.Mail.Brevo.ApiKey)

	return b.NewAPIClient(cfg)
}
