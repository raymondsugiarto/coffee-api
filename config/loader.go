package config

import (
	"log"
	"path/filepath"
	"runtime"
	"strings"

	// "github.com/raymondsugiarto/coffee-api/internal/config/sign_in_method"

	"github.com/spf13/viper"
)

// Config :
type Config struct {
	Server ServerList
	// Domain   map[string]sign_in_method.Domain
	Database DatabaseList
	Logger   LoggerConfig
	Aws      AwsConfig
	Mail     MailConfig
	Role     RoleConfig
	Whatsapp WhatsappConfig
	Cron     Cron
}

var configuration Config

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func init() {
	viper.AddConfigPath(basepath + "/resources")
	viper.SetConfigType("yml")
	viper.SetConfigName("server.yml")
	errConf := viper.ReadInConfig()
	if errConf != nil {
		log.Fatalf("Failed to load config: %v", errConf)
	}

	viper.SetConfigName("database.yml")
	errDb := viper.MergeInConfig()
	if errDb != nil {
		log.Fatalf("Cannot read database config: %v", errDb)
	}

	viper.SetConfigName("aws.yml")
	errAws := viper.MergeInConfig()
	if errAws != nil {
		log.Fatalf("cannot read aws config: %v", errAws)
	}

	viper.SetConfigName("logger.yml")
	errLog := viper.MergeInConfig()
	if errLog != nil {
		log.Fatalf("cannot load logger config: %v", errLog)
	}

	viper.SetConfigName("mail.yml")
	errMail := viper.MergeInConfig()
	if errMail != nil {
		log.Fatalf("cannot load mail config: %v", errMail)
	}

	viper.SetConfigName("role.yml")
	errRole := viper.MergeInConfig()
	if errRole != nil {
		log.Fatalf("cannot load role config: %v", errRole)
	}

	viper.SetConfigName("whatsapp.yml")
	errWhatsapp := viper.MergeInConfig()
	if errWhatsapp != nil {
		log.Fatalf("cannot load whatsapp config: %v", errWhatsapp)
	}

	viper.SetConfigName("cron.yml")
	errCron := viper.MergeInConfig()
	if errCron != nil {
		log.Fatalf("cannot load cron config: %v", errCron)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.Unmarshal(&configuration)
}

// GetConfig get config
func GetConfig() *Config {
	return &configuration
}
