package utils

const (
	ErrorTokenExpired = "Token is expired"

	RequestIDKey                  = "request_id"
	LoggerField                   = "logger_field"
	KeycloakInvalidGrant          = "invalid_grant"
	MessageInternalServerError    = "Internal server error"
	ErrorInputRequired            = "MSGCM001"
	ErrorInputFail                = "MSGCM002"
	ErrorPasswordFail             = "MSGCM003"
	ErrorEmailFail                = "MSGCM004"
	ErrorInputCharacterLimit      = "MSGCM005"
	ErrorTeamDuplicate            = "MSGCM006"
	ErrorPermissionsNameDuplicate = "MSGCM007"
	ErrorDuplicate                = "MSGCM008"

	ErrorTotalAssignPercentLagerThan100 = "MSGCM009"

	DefaultPerPage = 30

	MinimumAssignPercent = 30

	TimeFormatDate = "2006-01-02"

	OrganizationPrefixName = "MonstarLab"
	AdminUserName          = "admin"
	AdminPassword          = "password"
	MasterRealm            = "master"
	EmployeeManageRealm    = "employee_manage"
)
