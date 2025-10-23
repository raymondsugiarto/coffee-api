package routes

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/raymondsugiarto/coffee-api/config"
	"github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers"
	routes "github.com/raymondsugiarto/coffee-api/pkg/adapter/routes/admin"
	coroute "github.com/raymondsugiarto/coffee-api/pkg/adapter/routes/company"
	cr "github.com/raymondsugiarto/coffee-api/pkg/adapter/routes/customer"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/brevo"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/database"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/middleware"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/middleware/organization"
	"github.com/raymondsugiarto/coffee-api/pkg/module/admin"
	"github.com/raymondsugiarto/coffee-api/pkg/module/approval"
	"github.com/raymondsugiarto/coffee-api/pkg/module/article"
	bp "github.com/raymondsugiarto/coffee-api/pkg/module/benefit_participation"
	"github.com/raymondsugiarto/coffee-api/pkg/module/benefit_type"
	"github.com/raymondsugiarto/coffee-api/pkg/module/claim"
	"github.com/raymondsugiarto/coffee-api/pkg/module/country"
	"github.com/raymondsugiarto/coffee-api/pkg/module/district"
	"github.com/raymondsugiarto/coffee-api/pkg/module/estatement"
	feesetting "github.com/raymondsugiarto/coffee-api/pkg/module/fee_setting"
	"github.com/raymondsugiarto/coffee-api/pkg/module/province"
	"github.com/raymondsugiarto/coffee-api/pkg/module/regency"
	"github.com/raymondsugiarto/coffee-api/pkg/module/report/aum"
	"github.com/raymondsugiarto/coffee-api/pkg/module/report/channel"
	companyparticipant "github.com/raymondsugiarto/coffee-api/pkg/module/report/company_participant"
	contributionsummary "github.com/raymondsugiarto/coffee-api/pkg/module/report/contribution_summary"
	companyreport "github.com/raymondsugiarto/coffee-api/pkg/module/report/ojk_company_report"
	customerreport "github.com/raymondsugiarto/coffee-api/pkg/module/report/ojk_customer_report"
	participantsummary "github.com/raymondsugiarto/coffee-api/pkg/module/report/participant_summary"
	"github.com/raymondsugiarto/coffee-api/pkg/module/report/portfolio"
	"github.com/raymondsugiarto/coffee-api/pkg/module/ticket"
	userlog "github.com/raymondsugiarto/coffee-api/pkg/module/user_log"
	"github.com/raymondsugiarto/coffee-api/pkg/module/village"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	approvaltype "github.com/raymondsugiarto/coffee-api/pkg/module/approval/type"
	"github.com/raymondsugiarto/coffee-api/pkg/module/authentication"
	"github.com/raymondsugiarto/coffee-api/pkg/module/authentication/token"
	"github.com/raymondsugiarto/coffee-api/pkg/module/bank"
	"github.com/raymondsugiarto/coffee-api/pkg/module/company"
	"github.com/raymondsugiarto/coffee-api/pkg/module/customer"
	bankcustomer "github.com/raymondsugiarto/coffee-api/pkg/module/customer/bank_customer"
	customerpoint "github.com/raymondsugiarto/coffee-api/pkg/module/customer/customer_point"
	"github.com/raymondsugiarto/coffee-api/pkg/module/customer/participant"
	unitlink "github.com/raymondsugiarto/coffee-api/pkg/module/customer/unit_link"
	"github.com/raymondsugiarto/coffee-api/pkg/module/investment"
	investmentdistribution "github.com/raymondsugiarto/coffee-api/pkg/module/investment/investment_distribution"
	investmentitem "github.com/raymondsugiarto/coffee-api/pkg/module/investment/investment_item"
	investmentpayment "github.com/raymondsugiarto/coffee-api/pkg/module/investment/investment_payment"
	transactionfee "github.com/raymondsugiarto/coffee-api/pkg/module/investment/transaction_fee"
	investmentproduct "github.com/raymondsugiarto/coffee-api/pkg/module/investment_product"
	netassetvalue "github.com/raymondsugiarto/coffee-api/pkg/module/net_asset_value"
	"github.com/raymondsugiarto/coffee-api/pkg/module/notification"
	pensionbenefitrecipient "github.com/raymondsugiarto/coffee-api/pkg/module/pension_benefit_recipient"
	"github.com/raymondsugiarto/coffee-api/pkg/module/permission/permission"
	rolepermission "github.com/raymondsugiarto/coffee-api/pkg/module/permission/role_permission"
	"github.com/raymondsugiarto/coffee-api/pkg/module/redeem"
	summaryaum "github.com/raymondsugiarto/coffee-api/pkg/module/report/summary_aum"
	transactionhistory "github.com/raymondsugiarto/coffee-api/pkg/module/report/transaction_history"
	"github.com/raymondsugiarto/coffee-api/pkg/module/reward"
	"github.com/raymondsugiarto/coffee-api/pkg/module/role"
	"github.com/raymondsugiarto/coffee-api/pkg/module/user"
	usercredential "github.com/raymondsugiarto/coffee-api/pkg/module/user-credential"
	useridentityverification "github.com/raymondsugiarto/coffee-api/pkg/module/user_identity_verification"
	"github.com/raymondsugiarto/coffee-api/pkg/module/whatsapp"
)

func InitRouter(app fiber.Router) {

	dbConn := database.DBConn
	app.Use(logger.New(), organization.New(), middleware.DefaultResponseHandler())

	brevoClient := brevo.NewClient()

	// Country
	countryRepo := country.NewRepository(dbConn)
	countryService := country.NewService(countryRepo)

	// Province
	provinceRepo := province.NewRepository(dbConn)
	provinceService := province.NewService(provinceRepo)

	// Regency
	regencyRepo := regency.NewRepository(dbConn)
	regencyService := regency.NewService(regencyRepo)

	// District
	districtRepo := district.NewRepository(dbConn)
	districtService := district.NewService(districtRepo)

	// Village
	villageRepo := village.NewRepository(dbConn)
	villageService := village.NewService(villageRepo)

	// Notification
	notificationRepo := notification.NewRepository(dbConn)
	notificationService := notification.NewService(brevoClient, notificationRepo)

	// User Credential
	userCredentialRepo := usercredential.NewRepository(dbConn)
	userCredentialService := usercredential.NewService(userCredentialRepo)

	// Token
	tokenService := token.NewService()

	// User
	userRepo := user.NewRepository(dbConn)
	userService := user.NewService(userRepo, userCredentialService)

	// Pension Benefit Recipient
	benefitRecipientRepo := pensionbenefitrecipient.NewRepository(dbConn)
	benefitRecipientService := pensionbenefitrecipient.NewService(benefitRecipientRepo)

	// Approval
	approvalRepo := approval.NewRepository(dbConn)
	approvalService := approval.NewService(approvalRepo, notificationService)

	// Investment Item
	investmentItemRepo := investmentitem.NewRepository(dbConn)
	investmentItemService := investmentitem.NewService(investmentItemRepo)

	// Customer Point
	customerPointRepo := customerpoint.NewRepository(dbConn)
	customerPointService := customerpoint.NewService(customerPointRepo)

	// Customer
	customerRepo := customer.NewRepository(dbConn)

	// Whatsapp
	whatsappService := whatsapp.NewService()

	// Admin
	adminRepo := admin.NewRepository(dbConn)
	adminService := admin.NewService(adminRepo, config.GetConfig().Role.Company)

	// Company
	companyRepo := company.NewRepository(dbConn)
	companyService := company.NewService(companyRepo, adminService, config.GetConfig().Role.Company, approvalService, userCredentialService)

	// Reward
	rewardRepo := reward.NewRepository(dbConn)
	rewardService := reward.NewService(rewardRepo)
	StorageFileRouter(app.Group("/storage"))

	// Role Permission
	rolePermissionRepo := rolepermission.NewRepository(dbConn)
	rolePermissionService := rolepermission.NewService(rolePermissionRepo)

	// Role
	roleRepo := role.NewRepository(dbConn)
	roleService := role.NewService(roleRepo)

	// Permission
	permissionRepo := permission.NewRepository(dbConn)
	permissionService := permission.NewService(permissionRepo)

	// Investment Product
	investmentProductRepo := investmentproduct.NewRepository(dbConn)
	investmentProductService := investmentproduct.NewService(investmentProductRepo)

	// Transaction Fee
	transactionFeeRepo := transactionfee.NewRepository(dbConn)
	transactionFeeService := transactionfee.NewService(transactionFeeRepo)

	// Unit Link
	unitLinkRepo := unitlink.NewRepository(dbConn)

	// Participant
	participantRepo := participant.NewRepository(dbConn)

	unitLinkService := unitlink.NewService(unitLinkRepo, investmentItemService, investmentProductService)
	participantService := participant.NewService(participantRepo, unitLinkService)

	// Fee setting
	feeSettingRepo := feesetting.NewRepository(dbConn)
	feeSettingService := feesetting.NewService(feeSettingRepo)

	// Net Asset Value
	netAssetValueRepo := netassetvalue.NewRepository(dbConn)
	netAssetValueService := netassetvalue.NewService(netAssetValueRepo, unitLinkService, feeSettingService, transactionFeeService)
	investmentProductService.SetCallbackNetAssetValue(netAssetValueService.FindByInvestmentProductAndDate)

	customerService := customer.NewService(
		customerRepo, userCredentialService, userService,
		benefitRecipientService, tokenService, approvalService,
		participantService, customerPointService, notificationService,
		companyService, countryService, provinceService, regencyService, districtService, villageService,
		investmentItemService)

	// Redeem
	redeemRepo := redeem.NewRepository(dbConn)
	redeemService := redeem.NewService(redeemRepo, rewardRepo, customerPointRepo, customerService, customerPointService, rewardService, dbConn)

	userIdentityVerificationRepo := useridentityverification.NewRepository(dbConn)
	userIdentityVerificationService := useridentityverification.NewService(userIdentityVerificationRepo, customerService, whatsappService, notificationService, companyService)

	authenticationService := authentication.NewService(
		userCredentialService,
		customerService,
		userIdentityVerificationService,
		notificationService,
		companyService,
		whatsappService,
		tokenService,
		customerPointService,
		adminService,
	)

	auth := app.Group("/api/auth")
	authCompany := app.Group("/api/company/auth")
	authAdmin := app.Group("/api/admin/auth")
	authCustomer := app.Group("/api/customer/auth")

	AuthRouter(auth, userService, authenticationService, userIdentityVerificationService)
	CompanyAuthRouter(authCompany, userService, authenticationService, userIdentityVerificationService)
	RegencyRouter(authCompany.Group("/regencies"), regencyService)
	AdminAuthRouter(authAdmin, userService, authenticationService, userIdentityVerificationService)
	CustomerAuthRouter(authCustomer, userService, authenticationService, userIdentityVerificationService)

	// Middleware
	api := app.Group("/api", middleware.Protected())

	CountryRouter(api.Group("/countries"), countryService)
	ProvinceRouter(api.Group("/provinces"), provinceService)
	RegencyRouter(api.Group("/regencies"), regencyService)
	DistrictRouter(api.Group("/districts"), districtService)
	VillageRouter(api.Group("/villages"), villageService)

	// User Log
	userLogRepo := userlog.NewRepository(dbConn)
	userLogService := userlog.NewService(userLogRepo)

	// Investment Payment
	investmentPaymentRepo := investmentpayment.NewRepository(dbConn)
	investmentPaymentService := investmentpayment.NewService(investmentPaymentRepo, customerService, notificationService, participantService, approvalService, feeSettingService)

	// Investment Distribution
	investmentDistributionRepo := investmentdistribution.NewRepository(dbConn)
	investmentDistributionService := investmentdistribution.NewService(investmentDistributionRepo, participantService, investmentProductService, userLogService)

	// Investment
	investmentRepo := investment.NewRepository(dbConn)
	investmentService := investment.NewService(
		investmentRepo, participantService, investmentItemService,
		investmentPaymentService, investmentDistributionService, unitLinkService, netAssetValueService,
		notificationService, approvalService, feeSettingService, customerService,
	)

	// Claim
	claimRepo := claim.NewRepository(dbConn)
	claimService := claim.NewService(claimRepo, approvalService, unitLinkService, participantService)

	// Bank
	bankRepo := bank.NewRepository(dbConn)
	bankService := bank.NewService(bankRepo)

	// Article
	articleRepo := article.NewRepository(dbConn)
	articleService := article.NewService(articleRepo)

	// Benefit Type
	benefitTypeRepo := benefit_type.NewRepository(dbConn)
	benefitTypeService := benefit_type.NewService(benefitTypeRepo)

	// Benefit Registration
	benefitRegistrationRepo := bp.NewRepository(dbConn)
	benefitRegistrationService := bp.NewService(benefitRegistrationRepo, investmentService, investmentDistributionService, participantService, benefitTypeService, feeSettingService)

	// E-Statement
	eStatementService := estatement.NewEStatementService(investmentItemService, notificationService, customerService)

	// Company Report
	companyReportRepo := companyreport.NewRepository(dbConn)
	companyReportService := companyreport.NewOJKCompanyReportService(companyReportRepo)

	// Customer Report
	customerReportRepo := customerreport.NewRepository(dbConn)
	customerReportService := customerreport.NewOJKCustomerReportService(customerReportRepo, customerService)

	// Company Participant Report
	companyParticipantRepo := companyparticipant.NewRepository(dbConn)
	companyParticipantService := companyparticipant.NewService(companyParticipantRepo)

	// Ticket
	ticketRepo := ticket.NewRepository(dbConn)
	ticketService := ticket.NewService(ticketRepo, approvalService, customerService, userService)
	// Approval Callback
	approvalCompanyService := approvaltype.NewCompanyService(companyService)
	approvalService.SetCallbackCompany(approvalCompanyService)
	approvalCustomerService := approvaltype.NewCustomerService(customerService)
	approvalService.SetCallbackCustomer(approvalCustomerService)
	approvalClaimService := approvaltype.NewClaimService(claimService)
	approvalService.SetCallbackClaim(approvalClaimService)
	approvalInvestmentService := approvaltype.NewInvestmentService(investmentService, investmentPaymentService)
	approvalService.SetCallbackInvestment(approvalInvestmentService)
	approvalTicketService := approvaltype.NewTicketService(ticketService)
	approvalService.SetCallbackTicket(approvalTicketService)
	approvalBenefitParticipationSvc := approvaltype.NewBenefitParticipationService(benefitRegistrationService, investmentService, investmentPaymentService)
	approvalService.SetCallbackBenefitParticipation(approvalBenefitParticipationSvc)

	// Bank Customer
	bankCustomerRepo := bankcustomer.NewRepository(dbConn)
	bankCustomerService := bankcustomer.NewService(bankCustomerRepo)

	// Channel Report
	channelRepo := channel.NewRepository(dbConn)
	channelService := channel.NewService(channelRepo, netAssetValueService)
	participantSummaryRepo := participantsummary.NewRepository(dbConn)
	participantSummaryService := participantsummary.NewService(participantSummaryRepo, benefitTypeService)
	contributionSummaryRepo := contributionsummary.NewRepository(dbConn)
	contributionSummaryService := contributionsummary.NewService(contributionSummaryRepo)
	summaryAumRepo := summaryaum.NewRepository(dbConn)
	summaryAumService := summaryaum.NewService(summaryAumRepo)
	aumRepo := aum.NewRepository(dbConn)
	aumService := aum.NewService(aumRepo)

	// Transaction History Report
	transactionHistoryRepo := transactionhistory.NewRepository(dbConn)
	transactionHistoryService := transactionhistory.NewService(transactionHistoryRepo)

	// Portfolio Report
	portfolioRepo := portfolio.NewRepository(dbConn)
	portfolioService := portfolio.NewService(portfolioRepo)

	// Router Customer
	cr.ProfileRouter(api, customerService)
	CustomerInvestmentProductRouter(api.Group("/investment-products"), investmentProductService)
	StorageFileRouter(api.Group("/storage"))
	CustomerInvestmentRouter(api.Group("/investments"), investmentService, investmentPaymentService)
	CustomerBenefitRecipientRouter(api.Group("/benefit-recipients"), benefitRecipientService)
	CustomerBankRouter(api.Group("/banks"), bankService)
	CustomerParticipantRouter(api.Group("/participants"), participantService)
	cr.CustomerRouter(api.Group("/referrals"), customerService)
	CustomerClaimRouter(api.Group("/claims"), claimService)
	CustomerRewardRouter(api.Group("/rewards"), rewardService)
	CustomerRedeemRouter(api.Group("/redeems"), redeemService)
	cr.UnitLinkRouter(api, unitLinkService)
	CustomerTicketRouter(api.Group("/tickets"), ticketService)
	CustomerNotificationRouter(api.Group("/notifications"), notificationService)
	CustomerArticleRouter(api.Group("/articles"), articleService)
	EstatementRouter(api.Group("/estatements"), eStatementService)
	cr.BenefitParticipationRouter(api.Group("/benefit-participations"), benefitRegistrationService)
	CustomerBenefitTypeRouter(api.Group("/benefit-types"), benefitTypeService)
	cr.CustomerBankCustomerRouter(api.Group("/bank-customers"), bankCustomerService)

	// Router Admin
	apiAdmin := api.Group("/admin")
	routes.AdminRouter(apiAdmin, adminService)
	routes.AdminMyRouter(apiAdmin.Group("/me"), adminService, permissionService)
	routes.RoleRouter(apiAdmin.Group("/roles"), roleService, rolePermissionService)
	routes.PermissionRouter(apiAdmin.Group("/permissions"), permissionService)
	InvestmentProductRouter(apiAdmin.Group("/investment-products"), investmentProductService)
	StorageFileRouter(apiAdmin.Group("/storage"))
	routes.NetAssetValueRouter(apiAdmin.Group("/net-asset-values"), netAssetValueService)
	routes.CompanyRouter(apiAdmin.Group("/company"), companyService)
	cr.AdminCustomerRouter(apiAdmin.Group("/customers"), customerService)
	AdminBenefitRecipientRouter(apiAdmin.Group("/benefit-recipients"), benefitRecipientService)
	AdminBankRouter(apiAdmin.Group("/banks"), bankService)
	routes.ApprovalRouter(apiAdmin.Group("/approvals"), approvalService)
	routes.FeeSettingRouter(apiAdmin.Group("/fee-settings"), feeSettingService)
	routes.AdminInvestmentPaymentRouter(apiAdmin, investmentPaymentService)
	AdminArticleRouter(apiAdmin.Group("/articles"), articleService)
	AdminRewardRouter(apiAdmin.Group("/rewards"), rewardService)
	AdminRedeemRouter(apiAdmin.Group("/redeems"), redeemService)
	routes.AdminUnitLinkRouter(apiAdmin, unitLinkService)
	AdminBenefitTypeRouter(apiAdmin.Group("/benefit-types"), benefitTypeService)
	AdminDashboardRouter(apiAdmin.Group("/dashboard"), customerService, investmentProductService, unitLinkService, companyService)
	AdminCountryRouter(apiAdmin.Group("/countries"), countryService)
	AdminProvinceRouter(apiAdmin.Group("/provinces"), provinceService)
	AdminRegencyRouter(apiAdmin.Group("/regencies"), regencyService)
	AdminDistrictRouter(apiAdmin.Group("/districts"), districtService)
	AdminVillageRouter(apiAdmin.Group("/villages"), villageService)
	CustomerTicketRouter(apiAdmin.Group("/tickets"), ticketService)
	cr.AdminUploadRouter(apiAdmin.Group("/import"), customerService)
	routes.ReportRouter(apiAdmin.Group("/report"), channelService, participantSummaryService, contributionSummaryService, summaryAumService, aumService, companyReportService, customerReportService, transactionHistoryService, portfolioService)

	// Router Company
	apiCompany := api.Group("/company")
	CustomerInvestmentProductRouter(apiCompany.Group("/investment-products"), investmentProductService)
	coroute.CompanyCustomerRouter(apiCompany.Group("/customers"), customerService)

	coroute.CompanyMyRouter(apiCompany.Group("/me"), companyService, permissionService)
	CompanyInvestmentRouter(apiCompany.Group("/investments"), investmentService, investmentItemService, investmentPaymentService, investmentDistributionService)
	CompanyBenefitRecipientRouter(apiCompany.Group("/benefit-recipients"), benefitRecipientService)
	CompanyParticipantRouter(apiCompany.Group("participants"), participantService)
	CompanyClaimRouter(apiCompany.Group("/claims"), claimService)
	CompanyDashboardRouter(apiCompany.Group("/dashboard"), customerService, unitLinkService, investmentDistributionService, investmentPaymentService)
	CompanyCountryRouter(apiCompany.Group("/countries"), countryService)
	CompanyProvinceRouter(apiCompany.Group("/provinces"), provinceService)
	CompanyRegencyRouter(apiCompany.Group("/regencies"), regencyService)
	CompanyDistrictRouter(apiCompany.Group("/districts"), districtService)
	CompanyVillageRouter(apiCompany.Group("/villages"), villageService)
	coroute.FindAllCompanyUserLog(apiCompany.Group("/user-logs"), userLogService)
	cr.CompanyUploadRouter(apiCompany.Group("/import"), customerService)
	routes.AdminCompanyRouter(apiCompany.Group("/users"), adminService)
	coroute.ReportRouter(apiCompany.Group("/report"), companyReportService, customerReportService, companyParticipantService)

	go startScheduler(netAssetValueService)
}

func AdminBenefitTypeRouter(app fiber.Router, svc benefit_type.Service) {
	app.Post("", handlers.CreateBenefitType(svc))
	app.Get("", handlers.FindAllBenefitType(svc))
	app.Get("/:id", handlers.FindBenefitTypeByID(svc))
	app.Put("/:id", handlers.UpdateBenefitType(svc))
	app.Delete("/:id", handlers.DeleteBenefitType(svc))
}

func CustomerBenefitTypeRouter(app fiber.Router, svc benefit_type.Service) {
	app.Get("", handlers.FindAllActiveBenefitType(svc))
	app.Get("/:id", handlers.FindBenefitTypeByID(svc))
}

func startScheduler(netAssetValueService netassetvalue.Service) {
	// create a scheduler
	s, err := gocron.NewScheduler()
	if err != nil {
		// handle error
	}

	cfg := config.GetConfig()
	if cfg.Cron.MonthlyFee == "true" {
		// add a job to the scheduler
		_, err = s.NewJob(
			gocron.MonthlyJob(
				1,
				gocron.NewDaysOfTheMonth(1),
				gocron.NewAtTimes(
					gocron.NewAtTime(0, 0, 0),
				),
			),
			gocron.NewTask(
				netAssetValueService.ExecuteMonthlyFee,
			),
		)
	}
	if err != nil {
		// handle error
	}

	// start the scheduler
	s.Start()
}
