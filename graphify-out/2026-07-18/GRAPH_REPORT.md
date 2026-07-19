# Graph Report - coffee-api  (2026-07-18)

## Corpus Check
- 172 files · ~42,083 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1309 nodes · 2181 edges · 112 communities (102 shown, 10 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 205 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `bf90c990`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- [[_COMMUNITY_Customer Service & Reports|Customer Service & Reports]]
- [[_COMMUNITY_Bootstrap & Config Loader|Bootstrap & Config Loader]]
- [[_COMMUNITY_Admin Module Handlers|Admin Module Handlers]]
- [[_COMMUNITY_Customer & UserLog Repositories|Customer & UserLog Repositories]]
- [[_COMMUNITY_UnitLink Portfolio Repository|UnitLink Portfolio Repository]]
- [[_COMMUNITY_UnitLink Portfolio Service|UnitLink Portfolio Service]]
- [[_COMMUNITY_Order Entity & Repository|Order Entity & Repository]]
- [[_COMMUNITY_Customer Handlers & Routes|Customer Handlers & Routes]]
- [[_COMMUNITY_Admin Repository & Pagination|Admin Repository & Pagination]]
- [[_COMMUNITY_Response & Error Middleware|Response & Error Middleware]]
- [[_COMMUNITY_Item Repository & Entity|Item Repository & Entity]]
- [[_COMMUNITY_User Entity & DTOs|User Entity & DTOs]]
- [[_COMMUNITY_Participant Service|Participant Service]]
- [[_COMMUNITY_Bank Customer Service & Routes|Bank Customer Service & Routes]]
- [[_COMMUNITY_Admin Service|Admin Service]]
- [[_COMMUNITY_OJK Excel Report Handlers|OJK Excel Report Handlers]]
- [[_COMMUNITY_Customer Point Repository|Customer Point Repository]]
- [[_COMMUNITY_Order Service|Order Service]]
- [[_COMMUNITY_Participant Repository|Participant Repository]]
- [[_COMMUNITY_Bank Customer Repository|Bank Customer Repository]]
- [[_COMMUNITY_Customer Point Service|Customer Point Service]]
- [[_COMMUNITY_User Log Repository|User Log Repository]]
- [[_COMMUNITY_Company Handlers|Company Handlers]]
- [[_COMMUNITY_Company Repository & Entity|Company Repository & Entity]]
- [[_COMMUNITY_User Repository|User Repository]]
- [[_COMMUNITY_Item Service|Item Service]]
- [[_COMMUNITY_Portfolio Repository|Portfolio Repository]]
- [[_COMMUNITY_Role Repository|Role Repository]]
- [[_COMMUNITY_Role Service|Role Service]]
- [[_COMMUNITY_Pagination Filter Logic|Pagination Filter Logic]]
- [[_COMMUNITY_AWS S3 Service|AWS S3 Service]]
- [[_COMMUNITY_GORM Common Mixins|GORM Common Mixins]]
- [[_COMMUNITY_Config YAML Examples|Config YAML Examples]]
- [[_COMMUNITY_Order Item Entity|Order Item Entity]]
- [[_COMMUNITY_User Credential Service|User Credential Service]]
- [[_COMMUNITY_User Service|User Service]]
- [[_COMMUNITY_Upload Handlers & Routes|Upload Handlers & Routes]]
- [[_COMMUNITY_Admin Entity|Admin Entity]]
- [[_COMMUNITY_User Identity Verification Repo|User Identity Verification Repo]]
- [[_COMMUNITY_User Identity Verification Service|User Identity Verification Service]]
- [[_COMMUNITY_Role Entity|Role Entity]]
- [[_COMMUNITY_User Model|User Model]]
- [[_COMMUNITY_Order Model|Order Model]]
- [[_COMMUNITY_Organization Config|Organization Config]]
- [[_COMMUNITY_App Status Codes|App Status Codes]]
- [[_COMMUNITY_Company Service|Company Service]]
- [[_COMMUNITY_Order Item Repository|Order Item Repository]]
- [[_COMMUNITY_Order Item Service|Order Item Service]]
- [[_COMMUNITY_Authentication Service|Authentication Service]]
- [[_COMMUNITY_UserHasRole Entity|UserHasRole Entity]]
- [[_COMMUNITY_Storage Handlers & Routes|Storage Handlers & Routes]]
- [[_COMMUNITY_ItemCompany Model|ItemCompany Model]]
- [[_COMMUNITY_Order Payment Entity|Order Payment Entity]]
- [[_COMMUNITY_Organization Entity|Organization Entity]]
- [[_COMMUNITY_Order Item Model|Order Item Model]]
- [[_COMMUNITY_Role Permission Model|Role Permission Model]]
- [[_COMMUNITY_User Has Role Model|User Has Role Model]]
- [[_COMMUNITY_Unit Link Admin Router|Unit Link Admin Router]]
- [[_COMMUNITY_Auth Handler & Sign In|Auth Handler & Sign In]]
- [[_COMMUNITY_Mail Config (Brevo)|Mail Config (Brevo)]]
- [[_COMMUNITY_Server Config|Server Config]]
- [[_COMMUNITY_WhatsApp Config (Saungwa)|WhatsApp Config (Saungwa)]]
- [[_COMMUNITY_FindAllRequest & Pagination|FindAllRequest & Pagination]]
- [[_COMMUNITY_Item Model|Item Model]]
- [[_COMMUNITY_Admin Model|Admin Model]]
- [[_COMMUNITY_Order Payment Model|Order Payment Model]]
- [[_COMMUNITY_Role Model|Role Model]]
- [[_COMMUNITY_Organization Helpers|Organization Helpers]]
- [[_COMMUNITY_Pagination Request DTO|Pagination Request DTO]]
- [[_COMMUNITY_Sign In DTO|Sign In DTO]]
- [[_COMMUNITY_Database Config|Database Config]]
- [[_COMMUNITY_Logger Config|Logger Config]]
- [[_COMMUNITY_Company Model|Company Model]]
- [[_COMMUNITY_Permission Model|Permission Model]]
- [[_COMMUNITY_AWS Config|AWS Config]]
- [[_COMMUNITY_Role Config|Role Config]]
- [[_COMMUNITY_Organization Model|Organization Model]]
- [[_COMMUNITY_Response Helper|Response Helper]]
- [[_COMMUNITY_Cron Config|Cron Config]]
- [[_COMMUNITY_Go Module Manifest|Go Module Manifest]]
- [[_COMMUNITY_Community 87|Community 87]]
- [[_COMMUNITY_Community 88|Community 88]]
- [[_COMMUNITY_Community 89|Community 89]]
- [[_COMMUNITY_Community 90|Community 90]]
- [[_COMMUNITY_Community 91|Community 91]]
- [[_COMMUNITY_Community 92|Community 92]]
- [[_COMMUNITY_Community 93|Community 93]]
- [[_COMMUNITY_Community 94|Community 94]]
- [[_COMMUNITY_Community 96|Community 96]]
- [[_COMMUNITY_Community 97|Community 97]]
- [[_COMMUNITY_Community 99|Community 99]]
- [[_COMMUNITY_Community 103|Community 103]]
- [[_COMMUNITY_Community 104|Community 104]]
- [[_COMMUNITY_Community 105|Community 105]]
- [[_COMMUNITY_Community 106|Community 106]]
- [[_COMMUNITY_Community 107|Community 107]]
- [[_COMMUNITY_Community 108|Community 108]]
- [[_COMMUNITY_Community 109|Community 109]]
- [[_COMMUNITY_Community 110|Community 110]]
- [[_COMMUNITY_Community 112|Community 112]]
- [[_COMMUNITY_Community 113|Community 113]]
- [[_COMMUNITY_Community 114|Community 114]]
- [[_COMMUNITY_Community 115|Community 115]]
- [[_COMMUNITY_Community 116|Community 116]]

## God Nodes (most connected - your core abstractions)
1. `GetOrganization()` - 38 edges
2. `Service` - 30 edges
3. `GetUserCredential()` - 26 edges
4. `Context` - 25 edges
5. `Repository` - 23 edges
6. `Service` - 22 edges
7. `NewTable()` - 22 edges
8. `Context` - 19 edges
9. `Context` - 18 edges
10. `InitRouter()` - 17 edges

## Surprising Connections (you probably didn't know these)
- `Connect()` --calls--> `GetConfig()`  [INFERRED]
  cmd/db/migrate/db_migrate.go → config/loader.go
- `Connect()` --calls--> `GetDatabaseDriverMigration()`  [INFERRED]
  cmd/db/migrate/db_migrate.go → pkg/infrastructure/database/sql_connection.go
- `MigrateUpAll()` --calls--> `GetConfig()`  [INFERRED]
  cmd/db/migrate/db_migrate.go → config/loader.go
- `NewClient()` --calls--> `GetConfig()`  [INFERRED]
  pkg/infrastructure/brevo/client.go → config/loader.go
- `Protected()` --calls--> `GetConfig()`  [INFERRED]
  pkg/infrastructure/middleware/auth.go → config/loader.go

## Import Cycles
- None detected.

## Communities (112 total, 10 thin omitted)

### Community 0 - "Customer Service & Reports"
Cohesion: 0.07
Nodes (33): GenerateOJKCompanyReportExcel(), GenerateOJKCustomerReportExcel(), GetCompanyParticipantReport(), FindAllCompanyUserLog(), GetCompanyID(), GetOrigin(), GetOriginTypeKey(), NewBackgroundContext() (+25 more)

### Community 1 - "Bootstrap & Config Loader"
Cohesion: 0.36
Nodes (10): connect(), GetDatabaseDriverMigration(), getGormDialect(), getSqlDB(), NewSQLConnection(), SQLConnection, Dialector, Driver (+2 more)

### Community 2 - "Admin Module Handlers"
Cohesion: 0.10
Nodes (44): FindAllCompanies(), CreateItemCategory(), DeleteItemCategory(), FindAllItemCategories(), FindOneItemCategory(), UpdateItemCategory(), CreateItem(), DeleteItem() (+36 more)

### Community 3 - "Customer & UserLog Repositories"
Cohesion: 0.06
Nodes (38): GetOrganization(), Repository, NewRepository(), Service, NewService(), Service, NewService(), Context (+30 more)

### Community 4 - "UnitLink Portfolio Repository"
Cohesion: 0.12
Nodes (15): Context, DB, PortfolioFindAllRequest, ResultPagination, SumUnitLinkPortfolioDto, Time, UnitLinkDto, UnitLinkFindAllRequest (+7 more)

### Community 5 - "UnitLink Portfolio Service"
Cohesion: 0.11
Nodes (16): Context, DB, PortfolioFindAllRequest, Repository, ResultPagination, SumUnitLinkPortfolioDto, Time, UnitLinkDto (+8 more)

### Community 6 - "Order Entity & Repository"
Cohesion: 0.09
Nodes (23): NewOrderDtoFromModel(), OrderCountDto, OrderDto, OrderFindAllRequest, OrderInputDto, Repository, NewRepository(), OrderItemInputDto (+15 more)

### Community 7 - "Customer Handlers & Routes"
Cohesion: 0.14
Nodes (25): FindAllCompanyCustomer(), FindCompanyCustomerByID(), CompanyCustomerRouter(), ChangePassword(), CreateCustomer(), customerAttachment(), DeleteCustomer(), FindAllByCompany() (+17 more)

### Community 8 - "Admin Repository & Pagination"
Cohesion: 0.12
Nodes (26): Repository, NewRepository(), Service, NewService(), DriverFindAllRequest, Table, calculateTotalPages(), getMappingField() (+18 more)

### Community 9 - "Response & Error Middleware"
Cohesion: 0.12
Nodes (21): Error, DefaultErrorHandler(), DefaultResponseHandler(), Ctx, AppStatus, Time, ErrorResponse, Response (+13 more)

### Community 10 - "Item Repository & Entity"
Cohesion: 0.06
Nodes (32): NewItemCategoryDtoFromModel(), NewItemDtoFromModel(), ItemCategoryDto, ItemCategoryFindAllRequest, ItemDto, ItemFindAllRequest, Repository, NewRepository() (+24 more)

### Community 11 - "User Entity & DTOs"
Cohesion: 0.16
Nodes (16): CreateUser, CustomerAccount, CustomerAccountListItem, MyAccountProfile, StringToUserType(), UserCredential, UserDto, UserType (+8 more)

### Community 12 - "Participant Service"
Cohesion: 0.20
Nodes (10): Service, NewService(), ParticipantStatus, Context, CustomerDto, ParticipantDto, ParticipantFindAllRequest, Repository (+2 more)

### Community 13 - "Bank Customer Service & Routes"
Cohesion: 0.22
Nodes (10): Service, NewService(), CustomerBankCustomerRouter(), Router, Service, BankCustomerDto, BankCustomerFindAllRequest, Context (+2 more)

### Community 14 - "Admin Service"
Cohesion: 0.25
Nodes (9): Service, NewService(), AdminDto, Context, CreateAdminCompany, DB, FindAllRequest, Repository (+1 more)

### Community 15 - "OJK Excel Report Handlers"
Cohesion: 0.28
Nodes (14): formatIndonesianMonthForFilename(), GenerateOJKCompanyReportExcel(), GenerateOJKCustomerReportExcel(), ReportAUM(), ReportChannel(), ReportContributionSummary(), ReportParticipantSummary(), ReportPortfolio() (+6 more)

### Community 16 - "Customer Point Repository"
Cohesion: 0.28
Nodes (7): Repository, NewRepository(), Context, CustomerPointDto, DB, FindAllRequest, ResultPagination

### Community 17 - "Order Service"
Cohesion: 0.25
Nodes (8): Service, NewService(), Context, OrderCountDto, OrderDto, OrderFindAllRequest, Repository, ResultPagination

### Community 18 - "Participant Repository"
Cohesion: 0.29
Nodes (7): Repository, NewRepository(), Context, DB, ParticipantDto, ParticipantFindAllRequest, ResultPagination

### Community 19 - "Bank Customer Repository"
Cohesion: 0.28
Nodes (7): Repository, NewRepository(), BankCustomerDto, BankCustomerFindAllRequest, Context, DB, ResultPagination

### Community 20 - "Customer Point Service"
Cohesion: 0.28
Nodes (7): Service, NewService(), Context, CustomerPointDto, FindAllRequest, Repository, ResultPagination

### Community 21 - "User Log Repository"
Cohesion: 0.31
Nodes (7): Context, DB, ResultPagination, UserLogDto, UserLogFindAllRequest, Repository, NewRepository()

### Community 22 - "Company Handlers"
Cohesion: 0.35
Nodes (11): companyAttachment(), CompanyGetMyProfile(), CreateCompany(), DeleteCompany(), FindAllCompany(), FindCompanyByID(), UpdateCompany(), CompanyDto (+3 more)

### Community 23 - "Company Repository & Entity"
Cohesion: 0.16
Nodes (13): Repository, NewRepository(), NewCompanyDtoFromModel(), CompanyDto, CompanyFindAllRequest, Company, FindAllRequest, OrganizationDto (+5 more)

### Community 24 - "User Repository"
Cohesion: 0.27
Nodes (7): CreateUser, Context, DB, IdentityStatus, UserDto, Repository, NewRepository()

### Community 25 - "Item Service"
Cohesion: 0.05
Nodes (45): CashAdjustmentDto, EmployeeReportRowDto, CashAdjustmentDto, CashAdjustmentInputDto, CloseStockSessionInputDto, CloseStockSessionItemInputDto, DailyReportDto, DashboardSummaryDto (+37 more)

### Community 26 - "Portfolio Repository"
Cohesion: 0.29
Nodes (7): Context, DB, FindAllRequest, ResultPagination, Repository, NewRepository(), PortfolioDto

### Community 27 - "Role Repository"
Cohesion: 0.29
Nodes (7): Context, DB, FindAllRequest, ResultPagination, RoleDto, Repository, NewRepository()

### Community 28 - "Role Service"
Cohesion: 0.29
Nodes (7): Context, Repository, ResultPagination, UserLogDto, UserLogFindAllRequest, Service, NewService()

### Community 30 - "AWS S3 Service"
Cohesion: 0.29
Nodes (7): NewS3Service(), S3Service, UploadFileInput, DeleteObjectOutput, Reader, PutObjectOutput, S3Config

### Community 31 - "GORM Common Mixins"
Cohesion: 0.31
Nodes (5): CommonWithID, CommonWithIDs, DeletedAt, DB, Time

### Community 32 - "Config YAML Examples"
Cohesion: 0.27
Nodes (10): AWS, Cron, Database (Postgres), Redis, Logger, Mail (Brevo), Role, Server (REST) (+2 more)

### Community 33 - "Order Item Entity"
Cohesion: 0.27
Nodes (7): NewOrderItemDtoFromModel(), OrderItemDto, OrderItemInputDto, OrderItemPerItemCountDto, ItemDto, OrderItem, OrderItemDto

### Community 34 - "User Credential Service"
Cohesion: 0.36
Nodes (6): ChangePasswordDto, Context, Repository, UserCredentialDto, Service, NewService()

### Community 35 - "User Service"
Cohesion: 0.36
Nodes (6): Context, IdentityStatus, Repository, UserDto, Service, NewService()

### Community 36 - "Upload Handlers & Routes"
Cohesion: 0.31
Nodes (7): UploadCustomer(), AdminUploadRouter(), CompanyUploadRouter(), Handler, Service, Router, Service

### Community 37 - "Admin Entity"
Cohesion: 0.25
Nodes (6): AdminDto, CreateAdminCompany, DriverFindAllRequest, Admin, GetListRequest, UserDto

### Community 38 - "User Identity Verification Repo"
Cohesion: 0.44
Nodes (5): Context, DB, UserIdentityVerificationDto, Repository, NewRepository()

### Community 39 - "User Identity Verification Service"
Cohesion: 0.39
Nodes (6): Context, Repository, UserIdentityVerificationDto, Callback, Service, NewService()

### Community 40 - "Role Entity"
Cohesion: 0.29
Nodes (4): RoleDto, RoleInputDto, Role, RoleDto

### Community 41 - "User Model"
Cohesion: 0.32
Nodes (8): IdentityStatus, User, StringToUserType(), UserType, CommonWithIDs, Organization, UserType, UserCredential

### Community 42 - "Order Model"
Cohesion: 0.29
Nodes (8): Order, Admin, CommonWithIDs, Company, OrderItem, OrderPayment, Organization, Time

### Community 43 - "Organization Config"
Cohesion: 0.32
Nodes (7): Config, configDefault(), getOrganizationByOrigin(), New(), Ctx, Handler, Organization

### Community 44 - "App Status Codes"
Cohesion: 0.43
Nodes (7): AppStatus, ClientErrorCase, ServerErrorCase, NewClientErrorAppStatus(), NewServerErrorAppStatus(), NewSuccessAppStatus(), SuccessCase

### Community 45 - "Company Service"
Cohesion: 0.31
Nodes (7): Service, NewService(), CompanyDto, CompanyFindAllRequest, Context, Repository, ResultPagination

### Community 46 - "Order Item Repository"
Cohesion: 0.38
Nodes (6): Repository, NewRepository(), Context, DB, OrderFindAllRequest, OrderItemPerItemCountDto

### Community 47 - "Order Item Service"
Cohesion: 0.38
Nodes (6): Service, NewService(), Context, OrderFindAllRequest, OrderItemPerItemCountDto, Repository

### Community 48 - "Authentication Service"
Cohesion: 0.40
Nodes (4): Service, LoginRequestDto, Context, LoginDto

### Community 49 - "UserHasRole Entity"
Cohesion: 0.33
Nodes (4): UserHasRoleDto, RoleDto, UserDto, UserHasRole

### Community 50 - "Storage Handlers & Routes"
Cohesion: 0.33
Nodes (4): GetStorageFile(), Handler, Router, StorageFileRouter()

### Community 51 - "ItemCompany Model"
Cohesion: 0.33
Nodes (5): ItemCompany, CommonWithIDs, Company, Item, Organization

### Community 52 - "Order Payment Entity"
Cohesion: 0.60
Nodes (3): NewOrderPaymentDtoFromModel(), OrderPaymentDto, OrderPayment

### Community 53 - "Organization Entity"
Cohesion: 0.40
Nodes (4): OrganizationData, OrganizationDto, UserCredentialData, CommonWithIDs

### Community 54 - "Order Item Model"
Cohesion: 0.40
Nodes (4): OrderItem, CommonWithIDs, Item, Order

### Community 55 - "Role Permission Model"
Cohesion: 0.40
Nodes (4): RolePermission, Permission, CommonWithIDs, Role

### Community 56 - "User Has Role Model"
Cohesion: 0.40
Nodes (4): UserHasRole, CommonWithIDs, Role, User

### Community 57 - "Unit Link Admin Router"
Cohesion: 0.50
Nodes (3): AdminUnitLinkRouter(), Router, Service

### Community 58 - "Auth Handler & Sign In"
Cohesion: 0.50
Nodes (3): SignIn(), Handler, Service

### Community 59 - "Mail Config (Brevo)"
Cohesion: 0.50
Nodes (3): Brevo, Brevo, MailConfig

### Community 60 - "Server Config"
Cohesion: 0.67
Nodes (4): MessageBroker, Server, ServerList, Server

### Community 61 - "WhatsApp Config (Saungwa)"
Cohesion: 0.50
Nodes (3): Saungwa, WhatsappConfig, Saungwa

### Community 62 - "FindAllRequest & Pagination"
Cohesion: 0.50
Nodes (3): FindAllRequest, GetListRequest, OrganizationData

### Community 63 - "Item Model"
Cohesion: 0.47
Nodes (6): ItemCompany, Item, ItemCategory, CommonWithIDs, ItemCategory, Organization

### Community 65 - "Admin Model"
Cohesion: 0.67
Nodes (4): Admin, CommonWithIDs, Company, User

### Community 66 - "Order Payment Model"
Cohesion: 0.50
Nodes (3): OrderPayment, CommonWithIDs, Order

### Community 67 - "Role Model"
Cohesion: 0.67
Nodes (4): Role, CommonWithIDs, Organization, RolePermission

### Community 68 - "Organization Helpers"
Cohesion: 0.50
Nodes (3): AddFilterOrganizationData(), Ctx, PaginationRequestDto

### Community 69 - "Pagination Request DTO"
Cohesion: 0.67
Nodes (4): FilterItem, Pagination, PaginationRequestDto, ResultPagination

### Community 71 - "Database Config"
Cohesion: 1.00
Nodes (3): Database, Database, DatabaseList

### Community 72 - "Logger Config"
Cohesion: 1.00
Nodes (3): Logger, LoggerConfig, Logger

### Community 73 - "Company Model"
Cohesion: 1.00
Nodes (3): Company, CommonWithIDs, Organization

### Community 74 - "Permission Model"
Cohesion: 1.00
Nodes (3): Permission, CommonWithIDs, RolePermission

### Community 87 - "Community 87"
Cohesion: 0.12
Nodes (21): EmployeeSalary, EmployeeSalaryComponent, NewEmployeeSalaryComponentDtoFromModel(), NewEmployeeSalaryDtoFromModel(), EmployeeSalaryComponentDto, EmployeeSalaryDto, EmployeeSalaryFindAllRequest, SavePayrollRequest (+13 more)

### Community 88 - "Community 88"
Cohesion: 0.16
Nodes (12): NewSalaryComponentDtoFromModel(), SalaryComponentDto, SalaryComponentFindAllRequest, FindAllRequest, Context, DB, ResultPagination, SalaryComponentDto (+4 more)

### Community 89 - "Community 89"
Cohesion: 0.34
Nodes (8): Context, DB, Repository, ResultPagination, Service, StockSessionDto, StockSessionFindAllRequest, NewService()

### Community 90 - "Community 90"
Cohesion: 0.21
Nodes (10): Service, NewService(), Context, EmployeeSalaryDto, EmployeeSalaryFindAllRequest, Repository, ResultPagination, SimulatePayrollResultDto (+2 more)

### Community 91 - "Community 91"
Cohesion: 0.36
Nodes (6): ChangePasswordDto, Context, DB, UserCredentialDto, Repository, NewRepository()

### Community 92 - "Community 92"
Cohesion: 0.23
Nodes (7): DailyReportDto, DashboardSummaryDto, EmployeePerformanceRowDto, MonthlyReportDto, Context, service, TopProductRowDto

### Community 93 - "Community 93"
Cohesion: 0.08
Nodes (46): AdminGetMyProfile(), CreateAdminCompany(), FindAllAdmin(), FindAllAdminByCompanyID(), UpdateAdminName(), UpdateAdminProfileImage(), GetUserCredential(), CustomerGetMyProfile() (+38 more)

### Community 94 - "Community 94"
Cohesion: 0.18
Nodes (9): AwsConfig, Config, Cron, DatabaseList, LoggerConfig, MailConfig, RoleConfig, ServerList (+1 more)

### Community 96 - "Community 96"
Cohesion: 0.42
Nodes (7): Command, GetConfig(), Rest, initDatabase(), NewRest(), startRest(), startRestProduction()

### Community 97 - "Community 97"
Cohesion: 0.22
Nodes (8): StockSession, Admin, CashAdjustment, CommonWithIDs, PaymentDetail, SessionLog, StockSessionItem, Time

### Community 99 - "Community 99"
Cohesion: 0.57
Nodes (6): Migrate, Connect(), MigrateUpAll(), Migration(), migrationDown(), migrationUp()

### Community 103 - "Community 103"
Cohesion: 0.40
Nodes (5): jwtError(), Protected(), SuccessHandler(), Ctx, Handler

### Community 104 - "Community 104"
Cohesion: 0.33
Nodes (4): ErrorResponse, SetupValidator(), XValidator, Validate

### Community 105 - "Community 105"
Cohesion: 0.40
Nodes (5): EmployeeSalary, EmployeeSalaryComponent, Admin, CommonWithIDs, Time

### Community 106 - "Community 106"
Cohesion: 0.40
Nodes (4): Context, LoginDto, UserCredentialData, Service

### Community 107 - "Community 107"
Cohesion: 0.50
Nodes (3): NewS3Config(), S3Config, Client

### Community 108 - "Community 108"
Cohesion: 0.40
Nodes (4): SessionLog, CommonWithIDs, StockSession, Time

### Community 109 - "Community 109"
Cohesion: 0.40
Nodes (4): StockSessionItem, CommonWithIDs, Item, StockSession

### Community 110 - "Community 110"
Cohesion: 0.50
Nodes (3): FindAllDrivers(), Handler, Service

### Community 112 - "Community 112"
Cohesion: 0.50
Nodes (3): CashAdjustment, CommonWithIDs, StockSession

### Community 113 - "Community 113"
Cohesion: 0.50
Nodes (3): PaymentDetail, CommonWithIDs, StockSession

### Community 114 - "Community 114"
Cohesion: 0.50
Nodes (3): SalaryComponent, CommonWithIDs, Company

## Knowledge Gaps
- **292 isolated node(s):** `AwsConfig`, `Database`, `ServerList`, `DatabaseList`, `LoggerConfig` (+287 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **10 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `NewTable()` connect `Admin Repository & Pagination` to `Customer & UserLog Repositories`, `UnitLink Portfolio Repository`, `Order Entity & Repository`, `Item Repository & Entity`, `Customer Point Repository`, `Participant Repository`, `Bank Customer Repository`, `User Log Repository`, `Community 87`, `Company Repository & Entity`, `Community 88`, `Item Service`, `Portfolio Repository`, `Role Repository`?**
  _High betweenness centrality (0.212) - this node is a cross-community bridge._
- **Why does `GetOrganization()` connect `Customer & UserLog Repositories` to `Customer Service & Reports`, `User Identity Verification Repo`, `Company Service`, `Bank Customer Service & Routes`, `Admin Service`, `Order Service`, `Community 89`, `Community 90`, `Community 91`, `Role Service`?**
  _High betweenness centrality (0.165) - this node is a cross-community bridge._
- **Why does `GetUserCredential()` connect `Community 93` to `Customer Service & Reports`, `Admin Module Handlers`, `Customer Handlers & Routes`, `Bank Customer Service & Routes`, `Company Handlers`?**
  _High betweenness centrality (0.094) - this node is a cross-community bridge._
- **Are the 35 inferred relationships involving `GetOrganization()` (e.g. with `.CreateAdminCompany()` and `.Create()`) actually correct?**
  _`GetOrganization()` has 35 INFERRED edges - model-reasoned connections that need verification._
- **Are the 23 inferred relationships involving `GetUserCredential()` (e.g. with `AdminGetMyProfile()` and `UpdateAdminName()`) actually correct?**
  _`GetUserCredential()` has 23 INFERRED edges - model-reasoned connections that need verification._
- **What connects `AwsConfig`, `Database`, `ServerList` to the rest of the system?**
  _292 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Customer Service & Reports` be split into smaller, more focused modules?**
  _Cohesion score 0.07033315705975675 - nodes in this community are weakly interconnected._