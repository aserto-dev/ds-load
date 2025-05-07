package kc

import "time"

const TypeUser string = "user"

//nolint:tagliatelle // maintain json user formatting
type BaseUser struct {
	IID         string `json:"_id"`
	ID          string `json:"id"`
	Type        string `json:"type"`
	DisplayName string `json:"displayname"`
	Email       string `json:"email"`
	Username    string `json:"username"`
}

//nolint:tagliatelle // maintain json user formatting
type User struct {
	BaseUser
	SystemUsername                string        `json:"systemUsername"`
	Activated                     bool          `json:"activated"`
	FirstName                     string        `json:"firstname"`
	MiddleName                    string        `json:"middlename"`
	LastName                      string        `json:"lastname"`
	Organization                  string        `json:"organization"`
	JobTitle                      string        `json:"jobTitle"`
	Description                   string        `json:"description"`
	AccountLocked                 bool          `json:"account_locked"`
	AccountLockedDate             interface{}   `json:"account_locked_date"`
	Addresses                     []interface{} `json:"addresses"`
	AllowPublicKey                bool          `json:"allow_public_key"`
	AlternateEmail                interface{}   `json:"alternateEmail"`
	Attributes                    []interface{} `json:"attributes"`
	Company                       string        `json:"company"`
	CostCenter                    string        `json:"costCenter"`
	Created                       time.Time     `json:"created"`
	Department                    string        `json:"department"`
	DisableDeviceMaxLoginAttempts bool          `json:"disableDeviceMaxLoginAttempts"`
	EmployeeIdentifier            interface{}   `json:"employeeIdentifier"`
	EmployeeType                  string        `json:"employeeType"`
	EnableManagedUID              bool          `json:"enable_managed_uid"`
	EnableUserPortalMultiFactor   bool          `json:"enable_user_portal_multifactor"`
	ExternalDn                    string        `json:"external_dn"`
	ExternalSourceType            string        `json:"external_source_type"`
	ExternallyManaged             bool          `json:"externally_managed"`
	LdapBindingUser               bool          `json:"ldap_binding_user"`
	Location                      string        `json:"location"`
	ManagedAppleID                string        `json:"managedAppleId"`
	Manager                       interface{}   `json:"manager"`
	Mfa                           struct {
		Configured bool `json:"configured,omitempty"`
		Exclusion  bool `json:"exclusion,omitempty"`
	} `json:"mfa,omitempty"`
	MfaEnrollment struct {
		OverallStatus  string `json:"overallStatus,omitempty"`
		PushStatus     string `json:"pushStatus,omitempty"`
		TotpStatus     string `json:"totpStatus,omitempty"`
		WebAuthnStatus string `json:"webAuthnStatus,omitempty"`
	} `json:"mfaEnrollment,omitempty"`
	PasswordDate         time.Time     `json:"password_date"`
	PasswordExpired      bool          `json:"password_expired"`
	PasswordNeverExpires bool          `json:"password_never_expires"`
	PasswordlessSudo     bool          `json:"passwordless_sudo"`
	PhoneNumbers         []interface{} `json:"phoneNumbers"`
	RestrictedFields     []interface{} `json:"restrictedFields,omitempty"`
	SambaServiceUser     bool          `json:"samba_service_user"`
	SSHKeys              []interface{} `json:"ssh_keys"`
	State                string        `json:"state"`
	Sudo                 bool          `json:"sudo"`
	Suspended            bool          `json:"suspended"`
	TotpEnabled          bool          `json:"totp_enabled"`
	UnixGUID             int           `json:"unix_guid"`
	UnixUID              int           `json:"unix_uid"`
}
