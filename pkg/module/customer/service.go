package customer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	b "github.com/getbrevo/brevo-go/lib"
	"github.com/gofiber/fiber/v2/log"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	ea "github.com/raymondsugiarto/coffee-api/pkg/entity/authentication"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/module/approval"
	"github.com/raymondsugiarto/coffee-api/pkg/module/authentication/token"
	"github.com/raymondsugiarto/coffee-api/pkg/module/company"
	"github.com/raymondsugiarto/coffee-api/pkg/module/country"
	customerpoint "github.com/raymondsugiarto/coffee-api/pkg/module/customer/customer_point"
	"github.com/raymondsugiarto/coffee-api/pkg/module/customer/participant"
	"github.com/raymondsugiarto/coffee-api/pkg/module/district"
	investmentitem "github.com/raymondsugiarto/coffee-api/pkg/module/investment/investment_item"
	"github.com/raymondsugiarto/coffee-api/pkg/module/notification"
	pensionbenefitrecipient "github.com/raymondsugiarto/coffee-api/pkg/module/pension_benefit_recipient"
	"github.com/raymondsugiarto/coffee-api/pkg/module/province"
	"github.com/raymondsugiarto/coffee-api/pkg/module/regency"
	"github.com/raymondsugiarto/coffee-api/pkg/module/user"
	usercredential "github.com/raymondsugiarto/coffee-api/pkg/module/user-credential"
	"github.com/raymondsugiarto/coffee-api/pkg/module/village"
	sh "github.com/raymondsugiarto/coffee-api/pkg/shared"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/utils"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type Service interface {
	Create(ctx context.Context, req *entity.CustomerDto) (*entity.CustomerDto, error)
	FindByReferralCode(ctx context.Context, referralCode string) (*entity.CustomerDto, error)
	FindByEmail(ctx context.Context, email string) (*entity.CustomerDto, error)
	ChangePassword(ctx context.Context, dto *ea.UserIdentityVerificationDto) error
	VerifyPhoneNumber(ctx context.Context, dto *ea.UserIdentityVerificationDto) error
	VerifyEmail(ctx context.Context, dto *ea.UserIdentityVerificationDto) error
	FindByIDWithScope(ctx context.Context, id string, scopes []string) (*entity.CustomerDto, error)
	FindByID(ctx context.Context, id string) (*entity.CustomerDto, error)
	FindByUserID(ctx context.Context, userId string) (*entity.CustomerDto, error)
	FindByCompanyID(ctx context.Context, companyId string) (*entity.CustomerDto, error)
	Update(ctx context.Context, dto *entity.CustomerDto) (*entity.CustomerDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.CustomerFindAllRequest) (*pagination.ResultPagination, error)
	FindAllMyReferral(ctx context.Context, req *entity.CustomerFindAllRequest) (*pagination.ResultPagination, error)
	FindAllByCompany(ctx context.Context, req *entity.CustomerFindAllRequest) (*pagination.ResultPagination, error)
	ConfirmationApprovalCallback(ctx context.Context, req *entity.CustomerDto, tx *gorm.DB) (context.Context, error)
	CustomerChangePassword(ctx context.Context, request *entity.PasswordChangeInputDto) error
	UpdateMyProfile(ctx context.Context, req *entity.CustomerDto) (*entity.CustomerDto, error)
	CountByType(ctx context.Context) (*entity.CustomerCount, error)
	CountByTypeThisMonth(ctx context.Context) (*entity.CustomerCount, error)
	CreateBatch(ctx context.Context, req []*entity.CustomerDto) ([]*entity.CustomerDto, error)
	ParseExcelToCustomers(ctx context.Context, reader io.Reader) (*entity.ExcelImportResult, error)
	UpdateEffectiveDateIfFirstPayment(ctx context.Context, investmentID string, paymentAt time.Time) error
	GenerateExcel(ctx context.Context, req *entity.CustomerFindAllRequest, result *pagination.ResultPagination) ([]byte, error)
}

type emailCredentialParams struct {
	email     string
	firstName string
	password  string
	companyID string
	company   *entity.CompanyDto
}

type service struct {
	repository              Repository
	userCredentialService   usercredential.Service
	userService             user.Service
	benefitRecipientService pensionbenefitrecipient.Service
	approvalService         approval.Service
	tokenService            token.Service
	participantService      participant.Service
	customerPointService    customerpoint.Service
	notificationService     notification.Service
	companyService          company.Service
	countryService          country.Service
	provinceService         province.Service
	regencyService          regency.Service
	districtService         district.Service
	villageService          village.Service
	investmentItemService   investmentitem.Service
}

func NewService(
	repository Repository,
	userCredentialService usercredential.Service,
	userService user.Service,
	benefitRecipientService pensionbenefitrecipient.Service,
	tokenService token.Service,
	approvalService approval.Service,
	participantService participant.Service,
	customerPointService customerpoint.Service,
	notificationService notification.Service,
	companyService company.Service,
	countryService country.Service,
	provinceService province.Service,
	regencyService regency.Service,
	districtService district.Service,
	villageService village.Service,
	investmentItemService investmentitem.Service,
) Service {
	return &service{
		repository:              repository,
		userCredentialService:   userCredentialService,
		userService:             userService,
		benefitRecipientService: benefitRecipientService,
		tokenService:            tokenService,
		approvalService:         approvalService,
		participantService:      participantService,
		customerPointService:    customerPointService,
		notificationService:     notificationService,
		companyService:          companyService,
		countryService:          countryService,
		provinceService:         provinceService,
		regencyService:          regencyService,
		districtService:         districtService,
		villageService:          villageService,
		investmentItemService:   investmentItemService,
	}
}

func (s *service) Create(ctx context.Context, dto *entity.CustomerDto) (*entity.CustomerDto, error) {
	organizationID := shared.GetOrganization(ctx).ID

	// validate email
	c, err := s.FindByEmail(ctx, dto.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("errorFindByEmail: %v", err)
			return nil, status.New(status.BadRequest, fmt.Errorf("email %w sudah terdaftar", err))
		}
	}
	if c != nil {
		return nil, status.New(status.BadRequest, fmt.Errorf("email %s sudah terdaftar", dto.Email))
	}

	password, err := gonanoid.Generate(utils.PASSWORD_CHARSET, 8)
	if err != nil {
		log.Errorf("errorGeneratePassword: %v", err)
		return nil, status.New(status.InternalServerError, err)
	}

	hashPassword, _ := utils.HashPassword(password)
	user := &entity.UserDto{
		OrganizationID: organizationID,
		UserType:       entity.CUSTOMER,
		UserCredential: []entity.UserCredentialDto{
			{
				OrganizationID: organizationID,
				Username:       dto.Email,
				Password:       hashPassword,
			},
		},
	}

	if dto.OrganizationID == "" {
		dto.OrganizationID = organizationID
	}
	if dto.User == nil {
		dto.User = user
	}
	dto.User.EmailVerificationStatus = "UNVERIFIED"
	dto.User.PhoneVerificationStatus = "UNVERIFIED"

	var emailParams *emailCredentialParams
	if userType := shared.GetOriginTypeKey(ctx); userType != string(entity.CUSTOMER) {
		emailParams = &emailCredentialParams{
			email:     dto.Email,
			firstName: dto.FirstName,
			password:  password,
			companyID: dto.CompanyID,
			company:   dto.Company, // Pass existing company data if available
		}
	}

	result, err := s.repository.Create(ctx, dto, func(tx *gorm.DB) error {
		if dto.CompanyID != "" {
			if _, err := s.participantService.Create(ctx, dto.ToParticipantCompanyDto(), dto); err != nil {
				return err
			}

			var uid string
			if userCredential := shared.GetUserCredential(ctx); userCredential == nil {
				uid = dto.UserID
			} else {
				uid = userCredential.UserID
			}
			if dto.ApprovalStatus != model.ApprovalStatusApproved {
				if _, err := s.approvalService.CreateWithTx(ctx, dto.ToApprovalSubmitDto(uid), tx); err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	if emailParams != nil {
		go s.sendCredentialEmailIfAllowed(ctx, emailParams)
	}

	return result, nil
}

func (s *service) sendCredentialEmailIfAllowed(ctx context.Context, params *emailCredentialParams) {
	shouldSendEmail := true

	if params.companyID != "" {
		var company *entity.CompanyDto

		if params.company != nil {
			company = params.company
		} else {
			var err error
			company, err = s.companyService.FindByID(ctx, params.companyID)
			if err != nil {
				log.Errorf("errorFindCompanyByID: %v", err)
			}
		}

		// Only send email for PPIP companies, not DKP
		if company != nil && company.CompanyType == string(model.CompanyTypeDKP) {
			shouldSendEmail = false
		}
	}

	if shouldSendEmail {
		s.notificationService.SendEmailTemplate(ctx, &entity.NotificationInputDto{
			TemplateID: sh.TEMPLATE_ID_EMAIL_CHANGE_PASSWORD,
			To: []b.SendSmtpEmailTo{
				{
					Email: params.email,
				},
			},
			Data: map[string]interface{}{
				"name":     params.firstName,
				"email":    params.email,
				"password": params.password,
			},
		})
	}
}

func (s *service) FindByReferralCode(ctx context.Context, referralCode string) (*entity.CustomerDto, error) {
	return s.repository.FindByReferralCode(ctx, referralCode)
}

func (s *service) FindAllMyReferral(ctx context.Context, req *entity.CustomerFindAllRequest) (*pagination.ResultPagination, error) {
	if req.UserID == "" {
		return nil, status.New(status.BadRequest, errors.New("userId is required"))
	}

	customer, err := s.repository.FindByUserID(ctx, req.UserID)
	if err != nil {
		return nil, status.New(status.EntityNotFound, errors.New("customerNotFound"))
	}

	req.CustomerIDParent = customer.ID

	return s.repository.FindAll(ctx, req)
}

func (s *service) FindByEmail(ctx context.Context, email string) (*entity.CustomerDto, error) {
	return s.repository.FindByEmail(ctx, email)
}

func (s *service) VerifyPhoneNumber(ctx context.Context, req *ea.UserIdentityVerificationDto) error {
	return s.userService.UpdatePhoneVerificationStatus(ctx, req.UserID, model.VERIFIED)
}
func (s *service) VerifyEmail(ctx context.Context, req *ea.UserIdentityVerificationDto) error {
	err := s.userService.UpdateEmailVerificationStatus(ctx, req.UserID, model.VERIFIED)
	if err != nil {
		log.Errorf("errorUpdateEmailVerificationStatus: %v", err)
		return err
	}

	userCredentialDto, err := s.userCredentialService.FindByUsername(ctx, req.UserIdentity)
	if err != nil {
		return err
	}

	customer, err := s.FindByUserID(ctx, req.UserID)
	if err != nil {
		log.Errorf("errorFindCustomerByUserID: %v", err)
		return err
	}

	req.Data, err = s.tokenService.GenerateToken(ctx, entity.UserCredentialData{
		ID:         userCredentialDto.ID,
		UserID:     userCredentialDto.User.ID,
		CustomerID: customer.ID,
	})
	if err != nil {
		log.Errorf("errorGenerateToken: %v", err)
		return err
	}
	return nil
}

func (s *service) ChangePassword(ctx context.Context, req *ea.UserIdentityVerificationDto) error {
	userCredentials, err := s.userCredentialService.FindAllByUserID(ctx, req.UserID)
	if err != nil {
		log.Errorf("errorFindAllByUserID: %v", err)
		return err
	}
	data := req.Data.(*ea.UserIdentityVerificationInputPasswordDto)

	for _, userCredential := range userCredentials {
		cp := new(entity.ChangePasswordDto)
		cp.UserCredentialID = userCredential.ID
		cp.Password = data.Password
		if err := s.userCredentialService.ChangePassword(ctx, cp); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) CustomerChangePassword(ctx context.Context, request *entity.PasswordChangeInputDto) error {
	userCredentials, err := s.userCredentialService.FindAllByUserID(ctx, request.UserID)
	if err != nil {
		log.Errorf("errorFindAllByUserID: %v", err)
		return err
	}

	userCredential := userCredentials[0]

	if !utils.CheckPasswordHash(request.CurrentPassword, userCredential.Password) {
		return status.New(status.BadRequest, errors.New("current password is incorrect"))
	}

	if err := s.userCredentialService.ChangePassword(ctx, &entity.ChangePasswordDto{
		UserCredentialID: userCredential.ID,
		Password:         request.NewPassword,
	}); err != nil {
		return err
	}
	return nil
}

func (s *service) FindByIDWithScope(ctx context.Context, id string, scopes []string) (*entity.CustomerDto, error) {
	return s.repository.FindByIDWithScope(ctx, id, scopes)
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.CustomerDto, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *service) FindByUserID(ctx context.Context, userId string) (*entity.CustomerDto, error) {
	dto, err := s.repository.FindByUserID(ctx, userId)
	if err != nil {
		return nil, err
	}
	dto.Point, err = s.customerPointService.GetTotalPoint(ctx, dto.ID)
	if err != nil {
		log.Errorf("errorGetTotalPoint: %v", err)
		return nil, err
	}
	return dto, nil
}

func (s *service) FindByCompanyID(ctx context.Context, companyId string) (*entity.CustomerDto, error) {
	return s.repository.FindByCompanyID(ctx, companyId)
}

func (s *service) Update(ctx context.Context, req *entity.CustomerDto) (*entity.CustomerDto, error) {
	uid := shared.GetUserCredential(ctx).UserID
	return s.repository.Update(ctx, req, func(tx *gorm.DB) error {
		for _, recipient := range req.PensionBenefitRecipients {
			if recipient.ID == "" {
				if _, err := s.benefitRecipientService.CreateWithTx(ctx, recipient, tx); err != nil {
					return err
				}
			} else {
				if _, err := s.benefitRecipientService.UpdateWithTx(ctx, recipient, tx); err != nil {
					return err
				}
			}
		}
		if req.NeedApproval {
			if _, err := s.approvalService.CreateWithTx(ctx, req.ToApprovalSuspendDto(uid), tx); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}

func (s *service) FindAll(ctx context.Context, req *entity.CustomerFindAllRequest) (*pagination.ResultPagination, error) {
	return s.repository.FindAll(ctx, req)
}

func (s *service) FindAllByCompany(ctx context.Context, req *entity.CustomerFindAllRequest) (*pagination.ResultPagination, error) {
	companyID := shared.GetCompanyID(ctx)
	if companyID == nil {
		return nil, status.New(status.EntityNotFound, errors.New("companyIDNotFound"))
	}
	req.CompanyID = *companyID
	return s.repository.FindAll(ctx, req)
}

func (s *service) ConfirmationApprovalCallback(ctx context.Context, req *entity.CustomerDto, tx *gorm.DB) (context.Context, error) {
	// 1. Find customer by ID
	customer, err := s.repository.FindByIDWithScope(ctx, req.ID, []string{"complete"})
	if err != nil {
		return ctx, err
	}
	// 2. Check if customer approval status needs to be updated
	if customer.ApprovalStatus != req.ApprovalStatus {
		simStatusInactive := customer.SIMStatus == model.SIMStatusInactive
		requestApproved := req.ApprovalStatus == model.ApprovalStatusApproved
		if simStatusInactive && requestApproved {
			req.SIMStatus = model.SIMStatusActive
		}
		_, err := s.repository.Update(ctx, req, func(tx *gorm.DB) error { return nil })
		if err != nil {
			return ctx, err
		}
	}
	// 3. If already confirmed or after update, return nil
	return ctx, nil
}

func (s *service) UpdateMyProfile(ctx context.Context, req *entity.CustomerDto) (*entity.CustomerDto, error) {
	uid := shared.GetUserCredential(ctx).UserID
	recipient, err := s.benefitRecipientService.FindAll(ctx, &entity.RecipientFindAllRequest{
		CustomerID: req.ID,
	})
	if err != nil {
		return nil, err
	}

	if recipient.Count > 0 {
		if data, ok := recipient.Data.([]*entity.PensionBenefitRecipientDto); ok && len(data) > 0 {
			if len(req.PensionBenefitRecipients) > 0 {
				req.PensionBenefitRecipients[0].ID = data[0].ID
			}
		}
	}

	req.ApprovalStatus = model.ApprovalStatusKycPending

	return s.repository.Update(ctx, req, func(tx *gorm.DB) error {
		for _, recipient := range req.PensionBenefitRecipients {
			if recipient.ID == "" {
				if _, err := s.benefitRecipientService.CreateWithTx(ctx, recipient, tx); err != nil {
					return err
				}
			} else {
				if _, err := s.benefitRecipientService.UpdateWithTx(ctx, recipient, tx); err != nil {
					return err
				}
			}
		}
		if _, err := s.approvalService.CreateWithTx(ctx, req.ToApprovalKYCDto(uid), tx); err != nil {
			return err
		}
		return nil
	})
}

func (s *service) CountByType(ctx context.Context) (*entity.CustomerCount, error) {
	companyID := shared.GetCompanyID(ctx)

	result, err := s.repository.CountByType(ctx, companyID)
	if err != nil {
		return nil, err
	}

	result.EmployeeTotal = result.DKP + result.PPIPCorporate
	return result, nil
}

func (s *service) CountByTypeThisMonth(ctx context.Context) (*entity.CustomerCount, error) {
	companyID := shared.GetCompanyID(ctx)

	result, err := s.repository.CountByTypeThisMonth(ctx, companyID)
	if err != nil {
		return nil, err
	}

	result.EmployeeTotal = result.DKP + result.PPIPCorporate
	return result, nil
}

func (s *service) CreateBatch(ctx context.Context, req []*entity.CustomerDto) ([]*entity.CustomerDto, error) {
	companyMap := make(map[string]*entity.CompanyDto)

	for _, customer := range req {
		if customer.CompanyID != "" && customer.Company == nil {
			if _, exists := companyMap[customer.CompanyID]; !exists {
				company, err := s.companyService.FindByID(ctx, customer.CompanyID)
				if err == nil {
					companyMap[customer.CompanyID] = company
				} else {
					log.Errorf("errorFindCompanyByID in CreateBatch: %v", err)
				}
			}
			if companyMap[customer.CompanyID] != nil {
				customer.Company = companyMap[customer.CompanyID]
			}
		}
	}

	var result []*entity.CustomerDto
	for _, d := range req {
		dto, err := s.Create(ctx, d)
		if err != nil {
			return nil, err
		}
		result = append(result, dto)
	}
	return result, nil
}

func (s *service) ParseExcelToCustomers(ctx context.Context, reader io.Reader) (*entity.ExcelImportResult, error) {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, reader); err != nil {
		return nil, fmt.Errorf("failed to copy reader: %w", err)
	}

	xlFile, err := excelize.OpenReader(bytes.NewReader(buf.Bytes()))
	if err != nil {
		log.WithContext(ctx).Errorf("errorOpenExcelFile: %v", err)
		return nil, fmt.Errorf("failed to read Excel file: %w", err)
	}
	defer func() {
		if err := xlFile.Close(); err != nil {
			log.WithContext(ctx).Errorf("errorCloseExcelFile: %v", err)
			fmt.Println("failed to close excel file:", err)
		}
	}()

	// Get the first sheet name
	sheets := xlFile.GetSheetList()
	if len(sheets) == 0 {
		log.WithContext(ctx).Error("no sheets found in the Excel file")
		return nil, fmt.Errorf("no sheets found in the Excel file")
	}

	sheetName := sheets[0]

	// Read rows from the sheet
	rows, err := xlFile.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	var results []*entity.CustomerInputDto
	var errs []error

	for i, row := range rows {
		if i == 0 {
			continue // skip header
		}

		rowErr := func(msg string) {
			log.WithContext(ctx).Errorf("errorRow %d: %s: %s", i+1, row, msg)
			errs = append(errs, fmt.Errorf("baris %d: %s", i+1, msg))
		}

		// Ensure the row has enough columns
		if len(row) < 51 {
			rowErr("invalid row format in Excel file")
			continue
		}

		companyCode := row[0]
		if companyCode == "" {
			rowErr("kode perusahaan mandatory")
			continue
		}
		cityOfBirthCode := row[5]
		if cityOfBirthCode == "" {
			rowErr("tempat lahir mandatory")
			continue
		}
		countryCode := row[7]
		if countryCode == "" {
			rowErr("negara tempat lahir mandatory")
			continue
		}
		normalRetirementAge, err := utils.ToInt64(row[9])
		if err != nil {
			rowErr(fmt.Sprintf("invalid Usia Pensiun Normal: %s", err.Error()))
			continue
		}
		citizenCode := row[10]
		if citizenCode == "" {
			rowErr("kewarganegaraan mandatory")
			continue
		}
		provinceCode := row[24]
		if provinceCode == "" {
			rowErr("provinsi alamat mandatory")
			continue
		}
		regencyCode := row[25]
		if regencyCode == "" {
			rowErr("kota alamat mandatory")
			continue
		}
		districtCode := row[26]
		if districtCode == "" {
			rowErr("kecamatan alamat mandatory")
			continue
		}
		villageCode := row[27]
		if villageCode == "" {
			rowErr("kelurahan alamat mandatory")
			continue
		}
		mailingProvinceCode := row[31]
		if mailingProvinceCode == "" {
			rowErr("provinsi alamat domisili mandatory")
			continue
		}
		mailingRegencyCode := row[32]
		if mailingRegencyCode == "" {
			rowErr("kota alamat domisili mandatory")
			continue
		}
		mailingDistrictCode := row[33]
		if mailingDistrictCode == "" {
			rowErr("kecamatan alamat domisili mandatory")
			continue
		}
		mailingVillageCode := row[34]
		if mailingVillageCode == "" {
			rowErr("kelurahan alamat domisili mandatory")
			continue
		}
		countryRecipientCode := row[49]
		if countryRecipientCode == "" {
			rowErr("negara tempat lahir penerima manfaat mandatory")
			continue
		}
		employerAmount, err := utils.ToFloat64(row[41])
		if err != nil {
			rowErr(fmt.Sprintf("invalid Nominal Iuran Pemberi Kerja: %s", err.Error()))
			continue
		}
		employeeAmount, err := utils.ToFloat64(row[42])
		if err != nil {
			rowErr(fmt.Sprintf("invalid Nominal Iuran Peserta: %s", err.Error()))
			continue
		}
		voluntaryAmount, err := utils.ToFloat64(row[43])
		if err != nil {
			rowErr(fmt.Sprintf("invalid Nominal Iuran Sukarela Peserta: %s", err.Error()))
			continue
		}
		educationFundAmount, err := utils.ToFloat64(row[44])
		if err != nil {
			rowErr(fmt.Sprintf("invalid Nominal Dana Pendidikan Anak Peserta: %s", err.Error()))
			continue
		}
		var effectiveDate *string
		trimmedEffectiveDate := strings.Trim(row[45], " ")
		if trimmedEffectiveDate != "" {
			effectiveDate = &trimmedEffectiveDate
		} else {
			effectiveDate = nil
		}

		companyDto, err := s.companyService.FindByCompanyCode(ctx, companyCode)
		if err != nil {
			rowErr(fmt.Sprintf("kode perusahaan '%s': %s", companyCode, err.Error()))
			continue
		}
		cityOfBirthDto, err := s.regencyService.FindByCode(ctx, cityOfBirthCode)
		if err != nil {
			rowErr(fmt.Sprintf("tempat lahir '%s': %s", cityOfBirthCode, err.Error()))
			continue
		}
		countryDto, err := s.countryService.FindByCCA3(ctx, countryCode)
		if err != nil {
			rowErr(fmt.Sprintf("negara lahir '%s': %s", countryCode, err.Error()))
			continue
		}
		citizenDto, err := s.countryService.FindByCCA3(ctx, citizenCode)
		if err != nil {
			rowErr(fmt.Sprintf("kewarganegaraan '%s': %s", citizenCode, err.Error()))
			continue
		}
		provinceDto, err := s.provinceService.FindByCode(ctx, provinceCode)
		if err != nil {
			rowErr(fmt.Sprintf("provinsi alamat '%s': %s", provinceCode, err.Error()))
			continue
		}
		regencyDto, err := s.regencyService.FindByCode(ctx, regencyCode)
		if err != nil {
			rowErr(fmt.Sprintf("kota alamat '%s': %s", regencyCode, err.Error()))
			continue
		}
		districtDto, err := s.districtService.FindByCode(ctx, districtCode)
		if err != nil {
			rowErr(fmt.Sprintf("kecamatan alamat '%s': %s", districtCode, err.Error()))
			continue
		}
		villageDto, err := s.villageService.FindByCode(ctx, villageCode)
		if err != nil {
			rowErr(fmt.Sprintf("kelurahan alamat '%s': %s", villageCode, err.Error()))
			continue
		}
		mailingProvinceDto, err := s.provinceService.FindByCode(ctx, mailingProvinceCode)
		if err != nil {
			rowErr(fmt.Sprintf("provinsi domisili '%s': %s", mailingProvinceCode, err.Error()))
			continue
		}
		mailingRegencyDto, err := s.regencyService.FindByCode(ctx, mailingRegencyCode)
		if err != nil {
			rowErr(fmt.Sprintf("kota domisili '%s': %s", mailingRegencyCode, err.Error()))
			continue
		}
		mailingDistrictDto, err := s.districtService.FindByCode(ctx, mailingDistrictCode)
		if err != nil {
			rowErr(fmt.Sprintf("kecamatan domisili '%s': %s", mailingDistrictCode, err.Error()))
			continue
		}
		mailingVillageDto, err := s.villageService.FindByCode(ctx, mailingVillageCode)
		if err != nil {
			rowErr(fmt.Sprintf("kelurahan domisili '%s': %s", mailingVillageCode, err.Error()))
			continue
		}
		countryRecipientDto, err := s.countryService.FindByCCA3(ctx, countryRecipientCode)
		if err != nil {
			rowErr(fmt.Sprintf("negara lahir penerima manfaat '%s': %s", countryRecipientCode, err.Error()))
			continue
		}

		dto := &entity.CustomerInputDto{
			CompanyID:               companyDto.ID,
			EmployeeID:              row[1],
			FirstName:               row[2],
			Nickname:                row[3],
			DateStart:               &row[4],
			CityOfBirthID:           cityOfBirthDto.ID,
			DateOfBirth:             &row[6],
			CountryOfBirth:          countryDto.ID,
			MotherName:              row[8],
			NormalRetirementAge:     normalRetirementAge,
			Citizenship:             citizenDto.ID,
			Sex:                     row[11],
			MaritalStatus:           row[12],
			Occupation:              row[13],
			Position:                row[14],
			SourceOfFunds:           row[15],
			AnnualIncome:            row[16],
			PurposeOfAccount:        row[17],
			NameOnBankAccount:       row[18],
			BankAccountNumber:       row[19],
			BankName:                row[20],
			IdentificationNumber:    row[21],
			TaxIdentificationNumber: row[22],
			Address:                 row[23],
			ProvinceID:              provinceDto.ID,
			RegencyID:               regencyDto.ID,
			DistrictID:              districtDto.ID,
			VillageID:               villageDto.ID,
			RT:                      row[28],
			RW:                      row[29],
			PostalCode:              villageDto.PostalCode,
			MailingAddress:          row[30],
			MailingProvinceID:       mailingProvinceDto.ID,
			MailingRegencyID:        mailingRegencyDto.ID,
			MailingDistrictID:       mailingDistrictDto.ID,
			MailingVillageID:        mailingVillageDto.ID,
			MailingRT:               row[35],
			MailingRW:               row[36],
			MailingPostalCode:       mailingVillageDto.PostalCode,
			PhoneNumber:             row[37],
			OfficePhone:             row[38],
			MobilePhone:             row[39],
			Email:                   row[40],
			EmployerAmount:          employerAmount,
			EmployeeAmount:          employeeAmount,
			VoluntaryAmount:         voluntaryAmount,
			EducationFundAmount:     educationFundAmount,
			EffectiveDate:           effectiveDate,
			PensionBenefitRecipients: []*entity.PensionBenefitRecipientInputDto{
				{
					Name:                 row[46],
					Relationship:         row[47],
					DateOfBirth:          row[48],
					CountryOfBirth:       countryRecipientDto.ID,
					IdentificationNumber: row[50],
				},
			},
		}
		results = append(results, dto)
	}

	return &entity.ExcelImportResult{
		Data:   results,
		Errors: errs,
	}, nil
}

func (s *service) UpdateEffectiveDateIfFirstPayment(ctx context.Context, investmentID string, paymentAt time.Time) error {
	items, err := s.investmentItemService.FindByInvestmentID(ctx, investmentID)
	if err != nil {
		return err
	}

	for _, item := range items {
		customerID := item.CustomerID
		if customerID == "" {
			continue
		}

		has, err := s.investmentItemService.HasPreviousInvestmentForCustomer(ctx, investmentID, customerID)
		if err != nil {
			return err
		}
		if has {
			continue
		}

		customer, err := s.repository.FindByIDWithScope(ctx, customerID, []string{"complete"})
		if err != nil {
			return err
		}

		if customer.EffectiveDate == nil || *customer.EffectiveDate == "" {
			formatted := paymentAt.Format(time.RFC3339)
			customer.EffectiveDate = &formatted
			if _, err := s.repository.Update(ctx, customer, func(tx *gorm.DB) error {
				return nil
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *service) GenerateExcel(ctx context.Context, req *entity.CustomerFindAllRequest, result *pagination.ResultPagination) ([]byte, error) {
	excel := excelize.NewFile()
	sheetName := "Customers"
	index, _ := excel.NewSheet(sheetName)
	excel.DeleteSheet("Sheet1")

	headers := []string{"Kode Perusahaan", "ID Karyawan", "Nama Lengkap", "Nama Panggilan", "Tanggal Mulai Bekerja", "Tempat Lahir", "Tanggal Lahir", "Negara Lahir", "Nama Ibu", "Usia Pensiun Normal", "Kewarganegaraan", "Jenis Kelamin", "Status Perkawinan", "Pekerjaan", "Jabatan", "Sumber Dana", "Pendapatan Tahunan", "Tujuan Pembukaan Rekening", "Nama Rekening Bank", "Nomor Rekening Bank", "Nama Bank", "Nomor Identitas", "Nomor Pokok Wajib Pajak (NPWP)", "Alamat", "Provinsi Alamat", "Kota Alamat", "Kecamatan Alamat", "Kelurahan Alamat", "RT", "RW", "Kode Pos Alamat", "Alamat Domisili", "Provinsi Alamat Domisili", "Kota Alamat Domisili", "Kecamatan Alamat Domisili", "Kelurahan Alamat Domisili", "RT Domisili", "RW Domisili", "Kode Pos Domisili", "Nomor Telepon", "Telepon Kantor", "Nomor Ponsel", "Email", "Iuran Pemberi Kerja (Employer)", "Iuran Peserta (Employee)", "Iuran Sukarela Peserta (Voluntary)", "Dana Pendidikan Anak Peserta (Education Fund)", "Tanggal Efektif", "Nama Penerima Manfaat", "Hubungan", "Tanggal Lahir Penerima Manfaat", "Negara Lahir Penerima Manfaat", "Nomor Identitas Penerima Manfaat"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		excel.SetCellValue(sheetName, cell, header)
	}

	customers := result.Data.([]*entity.CustomerDto)
	for i, customer := range customers {
		row := i + 2 // Start from row 2 to skip header
		col1, _ := excelize.CoordinatesToCellName(1, row)
		if customer.Company != nil {
			excel.SetCellValue(sheetName, col1, customer.Company.CompanyCode)
		} else {
			excel.SetCellValue(sheetName, col1, customer.CompanyID)
		}
		col2, _ := excelize.CoordinatesToCellName(2, row)
		excel.SetCellValue(sheetName, col2, customer.EmployeeID)
		col3, _ := excelize.CoordinatesToCellName(3, row)
		excel.SetCellValue(sheetName, col3, customer.FirstName)
		col4, _ := excelize.CoordinatesToCellName(4, row)
		excel.SetCellValue(sheetName, col4, customer.Nickname)
		col5, _ := excelize.CoordinatesToCellName(5, row)
		if customer.DateStart != nil {
			excel.SetCellValue(sheetName, col5, *customer.DateStart)
		} else {
			excel.SetCellValue(sheetName, col5, "")
		}
		col6, _ := excelize.CoordinatesToCellName(6, row)
		if customer.CityOfBirth != nil {
			excel.SetCellValue(sheetName, col6, customer.CityOfBirth.Code)
		} else {
			excel.SetCellValue(sheetName, col6, customer.CityOfBirthID)
		}
		col7, _ := excelize.CoordinatesToCellName(7, row)
		if customer.DateOfBirth != nil {
			excel.SetCellValue(sheetName, col7, *customer.DateOfBirth)
		} else {
			excel.SetCellValue(sheetName, col7, "")
		}
		col8, _ := excelize.CoordinatesToCellName(8, row)
		if customer.CountryBirth != nil {
			excel.SetCellValue(sheetName, col8, customer.CountryBirth.CCA3)
		} else {
			excel.SetCellValue(sheetName, col8, customer.CountryOfBirth)
		}
		col9, _ := excelize.CoordinatesToCellName(9, row)
		excel.SetCellValue(sheetName, col9, customer.MotherName)
		col10, _ := excelize.CoordinatesToCellName(10, row)
		if customer.NormalRetirementAge != nil {
			excel.SetCellValue(sheetName, col10, *customer.NormalRetirementAge)
		} else {
			excel.SetCellValue(sheetName, col10, "")
		}
		col11, _ := excelize.CoordinatesToCellName(11, row)
		if customer.Citizen != nil {
			excel.SetCellValue(sheetName, col11, customer.Citizen.CCA3)
		} else {
			excel.SetCellValue(sheetName, col11, customer.Citizenship)
		}
		col12, _ := excelize.CoordinatesToCellName(12, row)
		excel.SetCellValue(sheetName, col12, customer.Sex)
		col13, _ := excelize.CoordinatesToCellName(13, row)
		excel.SetCellValue(sheetName, col13, customer.MaritalStatus)
		col14, _ := excelize.CoordinatesToCellName(14, row)
		excel.SetCellValue(sheetName, col14, customer.Occupation)
		col15, _ := excelize.CoordinatesToCellName(15, row)
		excel.SetCellValue(sheetName, col15, customer.Position)
		col16, _ := excelize.CoordinatesToCellName(16, row)
		excel.SetCellValue(sheetName, col16, customer.SourceOfFunds)
		col17, _ := excelize.CoordinatesToCellName(17, row)
		excel.SetCellValue(sheetName, col17, customer.AnnualIncome)
		col18, _ := excelize.CoordinatesToCellName(18, row)
		excel.SetCellValue(sheetName, col18, customer.PurposeOfAccount)
		col19, _ := excelize.CoordinatesToCellName(19, row)
		excel.SetCellValue(sheetName, col19, customer.NameOnBankAccount)
		col20, _ := excelize.CoordinatesToCellName(20, row)
		excel.SetCellValue(sheetName, col20, customer.BankAccountNumber)
		col21, _ := excelize.CoordinatesToCellName(21, row)
		excel.SetCellValue(sheetName, col21, customer.BankName)
		col22, _ := excelize.CoordinatesToCellName(22, row)
		excel.SetCellValue(sheetName, col22, customer.IdentificationNumber)
		col23, _ := excelize.CoordinatesToCellName(23, row)
		excel.SetCellValue(sheetName, col23, customer.TaxIdentificationNumber)
		col24, _ := excelize.CoordinatesToCellName(24, row)
		excel.SetCellValue(sheetName, col24, customer.Address)
		col25, _ := excelize.CoordinatesToCellName(25, row)
		if customer.Province != nil {
			excel.SetCellValue(sheetName, col25, customer.Province.Code)
		} else {
			excel.SetCellValue(sheetName, col25, customer.ProvinceID)
		}
		col26, _ := excelize.CoordinatesToCellName(26, row)
		if customer.Regency != nil {
			excel.SetCellValue(sheetName, col26, customer.Regency.Code)
		} else {
			excel.SetCellValue(sheetName, col26, customer.RegencyID)
		}
		col27, _ := excelize.CoordinatesToCellName(27, row)
		if customer.District != nil {
			excel.SetCellValue(sheetName, col27, customer.District.Code)
		} else {
			excel.SetCellValue(sheetName, col27, customer.DistrictID)
		}
		col28, _ := excelize.CoordinatesToCellName(28, row)
		if customer.Village != nil {
			excel.SetCellValue(sheetName, col28, customer.Village.Code)
		} else {
			excel.SetCellValue(sheetName, col28, customer.VillageID)
		}
		col29, _ := excelize.CoordinatesToCellName(29, row)
		excel.SetCellValue(sheetName, col29, customer.RT)
		col30, _ := excelize.CoordinatesToCellName(30, row)
		excel.SetCellValue(sheetName, col30, customer.RW)
		col31, _ := excelize.CoordinatesToCellName(31, row)
		excel.SetCellValue(sheetName, col31, customer.PostalCode)
		col32, _ := excelize.CoordinatesToCellName(32, row)
		excel.SetCellValue(sheetName, col32, customer.MailingAddress)
		col33, _ := excelize.CoordinatesToCellName(33, row)
		if customer.MailingProvince != nil {
			excel.SetCellValue(sheetName, col33, customer.MailingProvince.Code)
		} else {
			excel.SetCellValue(sheetName, col33, customer.MailingProvinceID)
		}
		col34, _ := excelize.CoordinatesToCellName(34, row)
		if customer.MailingRegency != nil {
			excel.SetCellValue(sheetName, col34, customer.MailingRegency.Code)
		} else {
			excel.SetCellValue(sheetName, col34, customer.MailingRegencyID)
		}
		col35, _ := excelize.CoordinatesToCellName(35, row)
		if customer.MailingDistrict != nil {
			excel.SetCellValue(sheetName, col35, customer.MailingDistrict.Code)
		} else {
			excel.SetCellValue(sheetName, col35, customer.MailingDistrictID)
		}
		col36, _ := excelize.CoordinatesToCellName(36, row)
		if customer.MailingVillage != nil {
			excel.SetCellValue(sheetName, col36, customer.MailingVillage.Code)
		} else {
			excel.SetCellValue(sheetName, col36, customer.MailingVillageID)
		}
		col37, _ := excelize.CoordinatesToCellName(37, row)
		excel.SetCellValue(sheetName, col37, customer.MailingRT)
		col38, _ := excelize.CoordinatesToCellName(38, row)
		excel.SetCellValue(sheetName, col38, customer.MailingRW)
		col39, _ := excelize.CoordinatesToCellName(39, row)
		excel.SetCellValue(sheetName, col39, customer.MailingPostalCode)
		col40, _ := excelize.CoordinatesToCellName(40, row)
		excel.SetCellValue(sheetName, col40, customer.PhoneNumber)
		col41, _ := excelize.CoordinatesToCellName(41, row)
		excel.SetCellValue(sheetName, col41, customer.OfficePhone)
		col42, _ := excelize.CoordinatesToCellName(42, row)
		excel.SetCellValue(sheetName, col42, customer.MobilePhone)
		col43, _ := excelize.CoordinatesToCellName(43, row)
		excel.SetCellValue(sheetName, col43, customer.Email)
		col44, _ := excelize.CoordinatesToCellName(44, row)
		if customer.EmployerAmount != nil {
			excel.SetCellValue(sheetName, col44, *customer.EmployerAmount)
		} else {
			excel.SetCellValue(sheetName, col44, 0.0)
		}
		col45, _ := excelize.CoordinatesToCellName(45, row)
		if customer.EmployeeAmount != nil {
			excel.SetCellValue(sheetName, col45, *customer.EmployeeAmount)
		} else {
			excel.SetCellValue(sheetName, col45, 0.0)
		}
		col46, _ := excelize.CoordinatesToCellName(46, row)
		if customer.VoluntaryAmount != nil {
			excel.SetCellValue(sheetName, col46, *customer.VoluntaryAmount)
		} else {
			excel.SetCellValue(sheetName, col46, 0.0)
		}
		col47, _ := excelize.CoordinatesToCellName(47, row)
		if customer.EducationFundAmount != nil {
			excel.SetCellValue(sheetName, col47, *customer.EducationFundAmount)
		} else {
			excel.SetCellValue(sheetName, col47, 0.0)
		}
		col48, _ := excelize.CoordinatesToCellName(48, row)
		if customer.EffectiveDate == nil {
			excel.SetCellValue(sheetName, col48, "")
		} else if *customer.EffectiveDate == "" {
			excel.SetCellValue(sheetName, col48, "")
		} else {
			excel.SetCellValue(sheetName, col48, *customer.EffectiveDate)
		}
		col49, _ := excelize.CoordinatesToCellName(49, row)
		col50, _ := excelize.CoordinatesToCellName(50, row)
		col51, _ := excelize.CoordinatesToCellName(51, row)
		col52, _ := excelize.CoordinatesToCellName(52, row)
		col53, _ := excelize.CoordinatesToCellName(53, row)
		if len(customer.PensionBenefitRecipients) > 0 {
			excel.SetCellValue(sheetName, col49, customer.PensionBenefitRecipients[0].Name)
			excel.SetCellValue(sheetName, col50, customer.PensionBenefitRecipients[0].Relationship)
			excel.SetCellValue(sheetName, col51, customer.PensionBenefitRecipients[0].DateOfBirth)
			if customer.PensionBenefitRecipients[0].CountryBirth != nil {
				excel.SetCellValue(sheetName, col52, customer.PensionBenefitRecipients[0].CountryBirth.CCA3)
			} else {
				excel.SetCellValue(sheetName, col52, customer.PensionBenefitRecipients[0].CountryOfBirth)
			}
			excel.SetCellValue(sheetName, col53, customer.PensionBenefitRecipients[0].IdentificationNumber)
		} else {
			excel.SetCellValue(sheetName, col49, "")
			excel.SetCellValue(sheetName, col50, "")
			excel.SetCellValue(sheetName, col51, "")
			excel.SetCellValue(sheetName, col52, "")
			excel.SetCellValue(sheetName, col53, "")
		}
	}

	excel.SetActiveSheet(index)

	var buf bytes.Buffer
	if err := excel.Write(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
