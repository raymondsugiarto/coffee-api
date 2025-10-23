package routes

import (
	"github.com/gofiber/fiber/v2"
	handlers "github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers/authentication"
	"github.com/raymondsugiarto/coffee-api/pkg/module/authentication"
	"github.com/raymondsugiarto/coffee-api/pkg/module/user"
	useridentityverification "github.com/raymondsugiarto/coffee-api/pkg/module/user_identity_verification"
)

func AuthRouter(app fiber.Router,
	service user.Service,
	authService authentication.Service,
	userIdentityVerificationService useridentityverification.Service,
) {
	app.Post("/sign-in", handlers.SignIn(authService))
	// app.Post("/admin/auth/sign-in", handlers.SignIn(authService))
	// app.Post("/company/auth/sign-in", handlers.SignIn(authService))
	// app.Post("/sign-up", handlers.SignUp(authService))

	app.Put("/user-identity-verification/:id/password", handlers.VerifyUserIdentityVerificationForPassword(userIdentityVerificationService))
	app.Put("/user-identity-verification/:id/email", handlers.VerifyUserIdentityVerificationForEmail(userIdentityVerificationService))
	app.Post("/user-identity-verification/:id", handlers.ResendUserIdentityVerification(userIdentityVerificationService))
}

func CompanyAuthRouter(app fiber.Router,
	service user.Service,
	authService authentication.Service,
	userIdentityVerificationService useridentityverification.Service,
) {

	app.Post("/sign-in", handlers.SignIn(authService))
	app.Post("/sign-up", handlers.SignUpCompany(authService))
	app.Post("/password", handlers.ForgotPasswordCompany(authService))
	app.Put("/user-identity-verification/:id/password", handlers.VerifyUserIdentityVerificationForPassword(userIdentityVerificationService))
	app.Put("/user-identity-verification/:id/email", handlers.VerifyUserIdentityVerificationForEmail(userIdentityVerificationService))
	app.Post("/user-identity-verification/:id", handlers.ResendUserIdentityVerification(userIdentityVerificationService))
}

func AdminAuthRouter(app fiber.Router,
	service user.Service,
	authService authentication.Service,
	userIdentityVerificationService useridentityverification.Service,
) {
	app.Post("/sign-in", handlers.SignIn(authService))
}

func CustomerAuthRouter(app fiber.Router,
	service user.Service,
	authService authentication.Service,
	userIdentityVerificationService useridentityverification.Service,
) {
	app.Post("/sign-up", handlers.SignUpCustomer(authService))
	app.Post("/password", handlers.ForgotPasswordCustomer(authService))
}
