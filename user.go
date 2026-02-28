package redcap

type User struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	FirstName       string `json:"firstname"`
	LastName        string `json:"lastname"`
	RoleID          int    `json:"role_id"`
	RoleLabel       string `json:"role_label"`
	DataAccessGroup string `json:"data_access_group"`
	Expiration      string `json:"expiration"`
	LastLogin       string `json:"last_login"`
	APIExport       bool   `json:"api_export"`
	APIImport       bool   `json:"api_import"`
	APIProject      bool   `json:"api_project"`
	MobileApp       bool   `json:"mobile_app"`
}
