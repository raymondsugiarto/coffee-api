package entity

import (
	"mime/multipart"

	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type CustomerImportDto struct {
	Document *multipart.FileHeader `json:"Document"`
}

type CustomerInputDto struct {
	OrganizationID           string                             `json:"organizationID,omitempty"`
	Email                    string                             `json:"email,omitempty"`
	PhoneNumber              string                             `json:"phoneNumber,omitempty"`
	FirstName                string                             `json:"firstName,omitempty"`
	SIMNumber                string                             `json:"simNumber,omitempty"`
	SIMStatus                string                             `json:"simStatus,omitempty"`
	ApprovalStatus           string                             `json:"approvalStatus,omitempty"`
	CompanyID                string                             `json:"companyId,omitempty"`
	EmployeeID               string                             `json:"employeeId,omitempty"`
	Nickname                 string                             `json:"nickname,omitempty"`
	DateStart                *string                            `json:"dateStart,omitempty"`
	PlaceOfBirth             string                             `json:"placeOfBirth,omitempty"`
	DateOfBirth              *string                            `json:"dateOfBirth,omitempty"`
	CountryOfBirth           string                             `json:"countryOfBirth,omitempty"`
	MotherName               string                             `json:"motherName,omitempty"`
	NormalRetirementAge      *int64                             `json:"normalRetirementAge,omitempty"`
	Citizenship              string                             `json:"citizenship,omitempty"`
	Sex                      string                             `json:"sex,omitempty"`
	MaritalStatus            string                             `json:"maritalStatus,omitempty"`
	Occupation               string                             `json:"occupation,omitempty"`
	Position                 string                             `json:"position,omitempty"`
	SourceOfFunds            string                             `json:"sourceOfFunds,omitempty"`
	AnnualIncome             string                             `json:"annualIncome,omitempty"`
	PurposeOfAccount         string                             `json:"purposeOfAccount,omitempty"`
	NameOnBankAccount        string                             `json:"nameOnBankAccount,omitempty"`
	BankAccountNumber        string                             `json:"bankAccountNumber,omitempty"`
	BankName                 string                             `json:"bankName,omitempty"`
	IdentificationNumber     string                             `json:"identificationNumber,omitempty"`
	TaxIdentificationNumber  string                             `json:"taxIdentificationNumber,omitempty"`
	Address                  string                             `json:"address,omitempty"`
	MailingAddress           string                             `json:"mailingAddress,omitempty"`
	OfficeAddress            string                             `json:"officeAddress,omitempty"`
	OfficePhone              string                             `json:"officePhone,omitempty"`
	MobilePhone              string                             `json:"mobilePhone,omitempty"`
	EmployerPercentage       *float64                           `json:"employerPercentage,omitempty"`
	EmployerAmount           *float64                           `json:"employerAmount,omitempty"`
	EmployeePercentage       *float64                           `json:"employeePercentage,omitempty"`
	EmployeeAmount           *float64                           `json:"employeeAmount,omitempty"`
	EffectiveDate            *string                            `json:"effectiveDate,omitempty"`
	PaymentMethod            string                             `json:"paymentMethod,omitempty"`
	IdentityCardFile         *multipart.FileHeader              `json:"identityCardFile,omitempty"`
	CustomerPhoto            *multipart.FileHeader              `json:"customerPhoto,omitempty"`
	PensionBenefitRecipients []*PensionBenefitRecipientInputDto `json:"pensionBenefitRecipients,omitempty"`
	NeedApproval             bool                               `json:"needApproval"`
	TaxIdentityCardFile      *multipart.FileHeader              `json:"taxIdentityCardFile,omitempty"`
	CityOfBirthID            string                             `json:"cityOfBirthID,omitempty"`
	ProvinceID               string                             `json:"provinceID,omitempty"`
	RegencyID                string                             `json:"regencyID,omitempty"`
	DistrictID               string                             `json:"districtID,omitempty"`
	VillageID                string                             `json:"villageID,omitempty"`
	RT                       string                             `json:"rt,omitempty"`
	RW                       string                             `json:"rw,omitempty"`
	PostalCode               string                             `json:"postalCode,omitempty"`
	MailingProvinceID        string                             `json:"mailingProvinceID,omitempty"`
	MailingRegencyID         string                             `json:"mailingRegencyID,omitempty"`
	MailingDistrictID        string                             `json:"mailingDistrictID,omitempty"`
	MailingVillageID         string                             `json:"mailingVillageID,omitempty"`
	MailingRT                string                             `json:"mailingRT,omitempty"`
	MailingRW                string                             `json:"mailingRW,omitempty"`
	MailingPostalCode        string                             `json:"mailingPostalCode,omitempty"`
	VoluntaryAmount          *float64                           `json:"voluntaryAmount,omitempty"`
	EducationFundAmount      *float64                           `json:"educationFundAmount,omitempty"`
}

func (u *CustomerInputDto) ToDto() *CustomerDto {
	dto := &CustomerDto{
		OrganizationID:          u.OrganizationID,
		Email:                   u.Email,
		PhoneNumber:             u.PhoneNumber,
		FirstName:               u.FirstName,
		SIMNumber:               u.SIMNumber,
		SIMStatus:               model.SIMStatus(u.SIMStatus),
		ApprovalStatus:          model.ApprovalStatus(u.ApprovalStatus),
		CompanyID:               u.CompanyID,
		EmployeeID:              u.EmployeeID,
		Nickname:                u.Nickname,
		DateStart:               u.DateStart,
		PlaceOfBirth:            u.PlaceOfBirth,
		DateOfBirth:             u.DateOfBirth,
		CountryOfBirth:          u.CountryOfBirth,
		MotherName:              u.MotherName,
		NormalRetirementAge:     u.NormalRetirementAge,
		Citizenship:             u.Citizenship,
		Sex:                     u.Sex,
		MaritalStatus:           u.MaritalStatus,
		Occupation:              u.Occupation,
		Position:                u.Position,
		SourceOfFunds:           u.SourceOfFunds,
		AnnualIncome:            u.AnnualIncome,
		PurposeOfAccount:        u.PurposeOfAccount,
		NameOnBankAccount:       u.NameOnBankAccount,
		BankAccountNumber:       u.BankAccountNumber,
		BankName:                u.BankName,
		IdentificationNumber:    u.IdentificationNumber,
		TaxIdentificationNumber: u.TaxIdentificationNumber,
		Address:                 u.Address,
		MailingAddress:          u.MailingAddress,
		OfficeAddress:           u.OfficeAddress,
		OfficePhone:             u.OfficePhone,
		MobilePhone:             u.MobilePhone,
		EmployerPercentage:      u.EmployerPercentage,
		EmployerAmount:          u.EmployerAmount,
		EmployeePercentage:      u.EmployeePercentage,
		EmployeeAmount:          u.EmployeeAmount,
		EffectiveDate:           u.EffectiveDate,
		PaymentMethod:           u.PaymentMethod,
		NeedApproval:            u.NeedApproval,
		CityOfBirthID:           u.CityOfBirthID,
		ProvinceID:              u.ProvinceID,
		RegencyID:               u.RegencyID,
		DistrictID:              u.DistrictID,
		VillageID:               u.VillageID,
		RT:                      u.RT,
		RW:                      u.RW,
		PostalCode:              u.PostalCode,
		MailingProvinceID:       u.MailingProvinceID,
		MailingRegencyID:        u.MailingRegencyID,
		MailingDistrictID:       u.MailingDistrictID,
		MailingVillageID:        u.MailingVillageID,
		MailingRT:               u.MailingRT,
		MailingRW:               u.MailingRW,
		MailingPostalCode:       u.MailingPostalCode,
		VoluntaryAmount:         u.VoluntaryAmount,
		EducationFundAmount:     u.EducationFundAmount,
	}

	dto.PensionBenefitRecipients = make([]*PensionBenefitRecipientDto, len(u.PensionBenefitRecipients))
	for i, recipient := range u.PensionBenefitRecipients {
		dto.PensionBenefitRecipients[i] = recipient.ToDto()
	}
	return dto
}

type CustomerDto struct {
	ID                       string                        `json:"id"`
	OrganizationID           string                        `json:"organizationID"`
	Username                 string                        `json:"username"`
	Email                    string                        `json:"email"`
	PhoneNumber              string                        `json:"phoneNumber"`
	ReferralCode             string                        `json:"referralCode"`
	FirstName                string                        `json:"firstName"`
	LastName                 string                        `json:"lastName"`
	SIMNumber                string                        `json:"simNumber"`
	SIMStatus                model.SIMStatus               `json:"simStatus"`
	ApprovalStatus           model.ApprovalStatus          `json:"approvalStatus"`
	EmployeeID               string                        `json:"employeeId"`
	Nickname                 string                        `json:"nickname"`
	DateStart                *string                       `json:"dateStart"`
	PlaceOfBirth             string                        `json:"placeOfBirth"`
	DateOfBirth              *string                       `json:"dateOfBirth"`
	CountryOfBirth           string                        `json:"countryOfBirth"`
	MotherName               string                        `json:"motherName"`
	NormalRetirementAge      *int64                        `json:"normalRetirementAge"`
	Citizenship              string                        `json:"citizenship"`
	Sex                      string                        `json:"sex"`
	MaritalStatus            string                        `json:"maritalStatus"`
	Occupation               string                        `json:"occupation"`
	Position                 string                        `json:"position"`
	SourceOfFunds            string                        `json:"sourceOfFunds"`
	AnnualIncome             string                        `json:"annualIncome"`
	PurposeOfAccount         string                        `json:"purposeOfAccount"`
	NameOnBankAccount        string                        `json:"nameOnBankAccount"`
	BankAccountNumber        string                        `json:"bankAccountNumber"`
	BankName                 string                        `json:"bankName"`
	IdentificationNumber     string                        `json:"identificationNumber"`
	TaxIdentificationNumber  string                        `json:"taxIdentificationNumber"`
	Address                  string                        `json:"address"`
	MailingAddress           string                        `json:"mailingAddress"`
	OfficeAddress            string                        `json:"officeAddress"`
	OfficePhone              string                        `json:"officePhone"`
	MobilePhone              string                        `json:"mobilePhone"`
	EmployerPercentage       *float64                      `json:"employerPercentage"`
	EmployerAmount           *float64                      `json:"employerAmount"`
	EmployeePercentage       *float64                      `json:"employeePercentage"`
	EmployeeAmount           *float64                      `json:"employeeAmount"`
	EffectiveDate            *string                       `json:"effectiveDate"`
	PaymentMethod            string                        `json:"paymentMethod"`
	IdentityCardFile         string                        `json:"identityCardFile"`
	CustomerPhoto            string                        `json:"customerPhoto"`
	CustomerIDParent         string                        `json:"customerIDParent"`
	CustomerParent           *CustomerDto                  `json:"customerParent"`
	PensionBenefitRecipients []*PensionBenefitRecipientDto `json:"pensionBenefitRecipients"`
	NeedApproval             bool                          `json:"needApproval"`
	TaxIdentityCardFile      string                        `json:"taxIdentityCardFile"`
	CityOfBirthID            string                        `json:"cityOfBirthId"`
	ProvinceID               string                        `json:"provinceId"`
	RegencyID                string                        `json:"regencyId"`
	DistrictID               string                        `json:"districtId"`
	VillageID                string                        `json:"villageId"`
	RT                       string                        `json:"rt"`
	RW                       string                        `json:"rw"`
	PostalCode               string                        `json:"postalCode"`
	MailingProvinceID        string                        `json:"mailingProvinceId"`
	MailingRegencyID         string                        `json:"mailingRegencyId"`
	MailingDistrictID        string                        `json:"mailingDistrictId"`
	MailingVillageID         string                        `json:"mailingVillageId"`
	MailingRT                string                        `json:"mailingRt"`
	MailingRW                string                        `json:"mailingRw"`
	MailingPostalCode        string                        `json:"mailingPostalCode"`
	VoluntaryAmount          *float64                      `json:"voluntaryAmount"`
	EducationFundAmount      *float64                      `json:"educationFundAmount"`

	CompanyID       string       `json:"companyId"`
	Company         *CompanyDto  `json:"company"`
	UserID          string       `json:"userID"`
	User            *UserDto     `json:"user"`
	Point           float64      `json:"point"`
	CityOfBirth     *RegencyDto  `json:"cityOfBirth"`
	Province        *ProvinceDto `json:"province"`
	Regency         *RegencyDto  `json:"regency"`
	District        *DistrictDto `json:"district"`
	Village         *VillageDto  `json:"village"`
	MailingProvince *ProvinceDto `json:"mailingProvince"`
	MailingRegency  *RegencyDto  `json:"mailingRegency"`
	MailingDistrict *DistrictDto `json:"mailingDistrict"`
	MailingVillage  *VillageDto  `json:"mailingVillage"`
	CountryBirth    *CountryDto  `json:"countryBirth"`
	Citizen         *CountryDto  `json:"citizen"`
}

func (dto *CustomerDto) ToModel() *model.Customer {
	m := &model.Customer{
		OrganizationID:          dto.OrganizationID,
		Email:                   dto.Email,
		PhoneNumber:             dto.PhoneNumber,
		ReferralCode:            dto.ReferralCode,
		FirstName:               dto.FirstName,
		LastName:                dto.LastName,
		SIMNumber:               dto.SIMNumber,
		SIMStatus:               dto.SIMStatus,
		ApprovalStatus:          dto.ApprovalStatus,
		CustomerID:              dto.EmployeeID,
		Nickname:                dto.Nickname,
		DateStart:               dto.DateStart,
		PlaceOfBirth:            dto.PlaceOfBirth,
		DateOfBirth:             dto.DateOfBirth,
		CountryOfBirth:          dto.CountryOfBirth,
		MotherName:              dto.MotherName,
		NormalRetirementAge:     dto.NormalRetirementAge,
		Citizenship:             dto.Citizenship,
		Sex:                     dto.Sex,
		MaritalStatus:           dto.MaritalStatus,
		Occupation:              dto.Occupation,
		Position:                dto.Position,
		SourceOfFunds:           dto.SourceOfFunds,
		AnnualIncome:            dto.AnnualIncome,
		PurposeOfOpeningAccount: dto.PurposeOfAccount,
		NameOnBankAccount:       dto.NameOnBankAccount,
		BankAccountNumber:       dto.BankAccountNumber,
		BankName:                dto.BankName,
		IdentificationNumber:    dto.IdentificationNumber,
		TaxIdentificationNumber: dto.TaxIdentificationNumber,
		Address:                 dto.Address,
		MailingAddress:          dto.MailingAddress,
		OfficeAddress:           dto.OfficeAddress,
		PhoneOffice:             dto.OfficePhone,
		MobilePhone:             dto.MobilePhone,
		EmployerPercentage:      dto.EmployerPercentage,
		EmployerAmount:          dto.EmployerAmount,
		CustomerPercentage:      dto.EmployeePercentage,
		CustomerAmount:          dto.EmployeeAmount,
		EffectiveDate:           dto.EffectiveDate,
		PaymentMethod:           dto.PaymentMethod,
		CustomerIDParent:        dto.CustomerIDParent,
		UserID:                  dto.UserID,
		CompanyID:               dto.CompanyID,
		CityOfBirthID:           dto.CityOfBirthID,
		ProvinceID:              dto.ProvinceID,
		RegencyID:               dto.RegencyID,
		DistrictID:              dto.DistrictID,
		VillageID:               dto.VillageID,
		RT:                      dto.RT,
		RW:                      dto.RW,
		PostalCode:              dto.PostalCode,
		MailingProvinceID:       dto.MailingProvinceID,
		MailingRegencyID:        dto.MailingRegencyID,
		MailingDistrictID:       dto.MailingDistrictID,
		MailingVillageID:        dto.MailingVillageID,
		MailingRT:               dto.MailingRT,
		MailingRW:               dto.MailingRW,
		MailingPostalCode:       dto.MailingPostalCode,
		VoluntaryAmount:         dto.VoluntaryAmount,
		EducationFundAmount:     dto.EducationFundAmount,
	}

	m.PensionBenefitRecipients = make([]*model.PensionBenefitRecipient, len(dto.PensionBenefitRecipients))
	for i, recipient := range dto.PensionBenefitRecipients {
		m.PensionBenefitRecipients[i] = recipient.ToModel()
	}
	if dto.IdentityCardFile != "" {
		m.IdentityCardFile = dto.IdentityCardFile
	}
	if dto.CustomerPhoto != "" {
		m.CustomerPhoto = dto.CustomerPhoto
	}
	if dto.TaxIdentityCardFile != "" {
		m.TaxIdentityCardFile = dto.TaxIdentityCardFile
	}
	if dto.ID != "" {
		m.ID = dto.ID
	}
	if dto.User != nil {
		m.User = dto.User.ToModel()
	}
	return m
}

func (dto *CustomerDto) FromModel(m *model.Customer) *CustomerDto {
	dto.ID = m.ID
	dto.OrganizationID = m.OrganizationID
	dto.Email = m.Email
	dto.PhoneNumber = m.PhoneNumber
	dto.ReferralCode = m.ReferralCode
	dto.FirstName = m.FirstName
	dto.LastName = m.LastName
	dto.SIMNumber = m.SIMNumber
	dto.SIMStatus = m.SIMStatus
	dto.ApprovalStatus = m.ApprovalStatus
	dto.EmployeeID = m.CustomerID
	dto.Nickname = m.Nickname
	dto.DateStart = m.DateStart
	dto.PlaceOfBirth = m.PlaceOfBirth
	dto.DateOfBirth = m.DateOfBirth
	dto.CountryOfBirth = m.CountryOfBirth
	dto.MotherName = m.MotherName
	dto.NormalRetirementAge = m.NormalRetirementAge
	dto.Citizenship = m.Citizenship
	dto.Sex = m.Sex
	dto.MaritalStatus = m.MaritalStatus
	dto.Occupation = m.Occupation
	dto.Position = m.Position
	dto.SourceOfFunds = m.SourceOfFunds
	dto.AnnualIncome = m.AnnualIncome
	dto.PurposeOfAccount = m.PurposeOfOpeningAccount
	dto.NameOnBankAccount = m.NameOnBankAccount
	dto.BankAccountNumber = m.BankAccountNumber
	dto.BankName = m.BankName
	dto.IdentificationNumber = m.IdentificationNumber
	dto.TaxIdentificationNumber = m.TaxIdentificationNumber
	dto.Address = m.Address
	dto.MailingAddress = m.MailingAddress
	dto.OfficeAddress = m.OfficeAddress
	dto.OfficePhone = m.PhoneOffice
	dto.MobilePhone = m.MobilePhone
	dto.EmployerPercentage = m.EmployerPercentage
	dto.EmployerAmount = m.EmployerAmount
	if m.EmployerAmount == nil {
		dto.EmployerAmount = new(float64)
		*dto.EmployerAmount = 0
	}
	dto.EmployeePercentage = m.CustomerPercentage
	dto.EmployeeAmount = m.CustomerAmount
	if m.CustomerAmount == nil {
		dto.EmployeeAmount = new(float64)
		*dto.EmployeeAmount = 0
	}
	dto.EffectiveDate = m.EffectiveDate
	dto.PaymentMethod = m.PaymentMethod
	dto.CustomerIDParent = m.CustomerIDParent
	dto.IdentityCardFile = m.IdentityCardFile
	dto.CustomerPhoto = m.CustomerPhoto
	dto.TaxIdentityCardFile = m.TaxIdentityCardFile
	dto.PensionBenefitRecipients = make([]*PensionBenefitRecipientDto, len(m.PensionBenefitRecipients))
	for i, recipient := range m.PensionBenefitRecipients {
		dto.PensionBenefitRecipients[i] = (&PensionBenefitRecipientDto{}).FromModel(recipient)
	}
	dto.CompanyID = m.CompanyID
	if m.Company != nil {
		dto.Company = (&CompanyDto{}).FromModel(m.Company)
	}
	dto.UserID = m.UserID
	if m.User != nil {
		dto.User = (&UserDto{}).FromModel(m.User)
	}
	dto.CityOfBirthID = m.CityOfBirthID
	if m.CityOfBirth != nil {
		dto.CityOfBirth = (&RegencyDto{}).FromModel(m.CityOfBirth)
	}
	dto.ProvinceID = m.ProvinceID
	if m.Province != nil {
		dto.Province = (&ProvinceDto{}).FromModel(m.Province)
	}
	dto.RegencyID = m.RegencyID
	if m.Regency != nil {
		dto.Regency = (&RegencyDto{}).FromModel(m.Regency)
	}
	dto.DistrictID = m.DistrictID
	if m.District != nil {
		dto.District = (&DistrictDto{}).FromModel(m.District)
	}
	dto.VillageID = m.VillageID
	if m.Village != nil {
		dto.Village = (&VillageDto{}).FromModel(m.Village)
	}
	dto.RT = m.RT
	dto.RW = m.RW
	dto.PostalCode = m.PostalCode
	dto.MailingProvinceID = m.MailingProvinceID
	if m.MailingProvince != nil {
		dto.MailingProvince = (&ProvinceDto{}).FromModel(m.MailingProvince)
	}
	dto.MailingRegencyID = m.MailingRegencyID
	if m.MailingRegency != nil {
		dto.MailingRegency = (&RegencyDto{}).FromModel(m.MailingRegency)
	}
	dto.MailingDistrictID = m.MailingDistrictID
	if m.MailingDistrict != nil {
		dto.MailingDistrict = (&DistrictDto{}).FromModel(m.MailingDistrict)
	}
	dto.MailingVillageID = m.MailingVillageID
	if m.MailingVillage != nil {
		dto.MailingVillage = (&VillageDto{}).FromModel(m.MailingVillage)
	}
	dto.MailingRT = m.MailingRT
	dto.MailingRW = m.MailingRW
	dto.MailingPostalCode = m.MailingPostalCode
	if m.CountryBirth != nil {
		dto.CountryBirth = (&CountryDto{}).FromModel(m.CountryBirth)
	}
	dto.VoluntaryAmount = m.VoluntaryAmount
	if m.VoluntaryAmount == nil {
		dto.VoluntaryAmount = new(float64)
		*dto.VoluntaryAmount = 0
	}
	dto.EducationFundAmount = m.EducationFundAmount
	if m.EducationFundAmount == nil {
		dto.EducationFundAmount = new(float64)
		*dto.EducationFundAmount = 0
	}

	if m.Citizen != nil {
		dto.Citizen = (&CountryDto{}).FromModel(m.Citizen)
	}

	if m.Company != nil {
		dto.Company = (&CompanyDto{}).FromModel(m.Company)
	}
	return dto
}

type CustomerPasswordDto struct {
	UserID   string `json:"userID"`
	Password string `json:"password"`
}

type CustomerFindAllRequest struct {
	FindAllRequest
	CompanyID        string
	CustomerIDParent string
	UserID           string
	SIMStatus        string
}

func (dto *CustomerFindAllRequest) GenerateFilter() {
	if dto.CustomerIDParent != "" {
		dto.AddFilter(
			pagination.FilterItem{
				Field: "customer_id_parent",
				Op:    "eq",
				Val:   dto.CustomerIDParent,
			},
		)
	}
	if dto.CompanyID != "" {
		dto.AddFilter(
			pagination.FilterItem{
				Field: "company_id",
				Op:    "eq",
				Val:   dto.CompanyID,
			},
		)
	}
	if dto.SIMStatus != "" {
		dto.AddFilter(
			pagination.FilterItem{
				Field: "sim_status",
				Op:    "eq",
				Val:   dto.SIMStatus,
			},
		)
	}
}

func (dto *CustomerDto) ToApprovalSuspendDto(uid string) *ApprovalDto {
	return &ApprovalDto{
		OrganizationID: dto.OrganizationID,
		UserIDRequest:  uid,
		RefID:          dto.ID,
		RefTable:       "customer",
		Detail:         "Perubahan Status Pelanggan [" + dto.FirstName + "]",
		Type:           "CUSTOMER",
		Action:         model.ApprovalAction(dto.SIMStatus),
		Status:         string(dto.ApprovalStatus),
		Reason:         "Update Customer Status",
	}
}

func (dto *CustomerDto) ToApprovalKYCDto(uid string) *ApprovalDto {
	return &ApprovalDto{
		OrganizationID: dto.OrganizationID,
		UserIDRequest:  uid,
		RefID:          dto.ID,
		RefTable:       "customer",
		Detail:         "Perubahan Data KYC [" + dto.FirstName + "]",
		Type:           "CUSTOMER",
		Action:         "UPDATE",
		Status:         string(dto.ApprovalStatus),
		Reason:         "Update Customer Information",
	}
}

func (dto *CustomerDto) ToApprovalSubmitDto(uid string) *ApprovalDto {
	return &ApprovalDto{
		OrganizationID: dto.OrganizationID,
		UserIDRequest:  uid,
		RefID:          dto.ID,
		RefTable:       "customer",
		Detail:         "Pendaftaran Pelanggan Baru [" + dto.FirstName + "]",
		Type:           "CUSTOMER",
		Action:         "ADD",
		Status:         "SUBMIT",
		Reason:         "New Participant Registration",
	}
}

func (dto *CustomerDto) ToParticipantCompanyDto() *ParticipantDto {
	return &ParticipantDto{
		OrganizationID: dto.OrganizationID,
		CustomerID:     dto.ID,
		Type:           model.InvestmentTypeDKP,
	}
}

func (dto *CustomerDto) ToCustomerPointDto() *CustomerPointDto {
	return &CustomerPointDto{
		CustomerID:  dto.CustomerIDParent,
		Point:       POINT_REFERRAL,
		Description: "Referral from " + dto.FirstName + " " + dto.LastName,
		RefID:       dto.ID,
		RefCode:     dto.ReferralCode,
		RefModule:   "REFERRAL",
	}
}

type CustomerCount struct {
	DKP           int64 `json:"dkp" gorm:"column:dkp"`
	PPIPCorporate int64 `json:"ppipCorporate" gorm:"column:ppip_corporate"`
	PPIPMandiri   int64 `json:"ppipMandiri" gorm:"column:ppip_mandiri"`
	EmployeeTotal int64 `json:"employeeTotal"`
}

type ExcelImportResult struct {
	Data   []*CustomerInputDto
	Errors []error
}

func (i *CustomerDto) GetInfo() RejectEmail {
	return RejectEmail{
		Email:       i.Email,
		Name:        i.FirstName,
		Description: "Pengajuan Pendaftaran " + i.FirstName,
	}
}
