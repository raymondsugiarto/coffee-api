package whatsapp

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/config"
	"golang.org/x/net/context"
)

type Service interface {
	SendMessage(ctx context.Context, to string, message string)
}

type service struct {
}

func NewService() Service {
	return &service{}
}

func (s *service) SendMessage(ctx context.Context, to string, message string) {
	agent := fiber.Post(config.GetConfig().Whatsapp.Saungwa.Url)

	args := fiber.AcquireArgs()
	args.Set("appkey", config.GetConfig().Whatsapp.Saungwa.Appkey)
	args.Set("authkey", config.GetConfig().Whatsapp.Saungwa.Authkey)
	args.Set("to", to)
	args.Set("message", message)
	defer fiber.ReleaseArgs(args)

	req := agent.Request()
	if req == nil {
		fmt.Println("error: failed to acquire request")
		return
	}

	req.Header.SetMethod(fiber.MethodPost)
	agent.Form(args)
	agent.Debug()

	statusCode, body, errs := agent.Bytes()
	if len(errs) > 0 {
		fmt.Printf("error: %v\n", errs)
		return
	}

	fmt.Printf("statusCode: %d\n", statusCode)
	fmt.Printf("body: %s\n", body[:100])
}
