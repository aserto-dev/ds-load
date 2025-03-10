package jc

import "time"

type User struct {
	JCID                          string        `json:"_id"`
	ID                            string        `json:"id"`
	Username                      string        `json:"username"`
	Displayname                   string        `json:"displayname"`
	Email                         string        `json:"email"`
	SystemUsername                string        `json:"systemUsername"`
	Activated                     bool          `json:"activated"`
	Firstname                     string        `json:"firstname"`
	Middlename                    string        `json:"middlename"`
	Lastname                      string        `json:"lastname"`
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
	EnableUserPortalMultifactor   bool          `json:"enable_user_portal_multifactor"`
	ExternalDn                    string        `json:"external_dn"`
	ExternalSourceType            string        `json:"external_source_type"`
	ExternallyManaged             bool          `json:"externally_managed"`
	LdapBindingUser               bool          `json:"ldap_binding_user"`
	Location                      string        `json:"location"`
	ManagedAppleID                string        `json:"managedAppleId"`
	Manager                       interface{}   `json:"manager"`
	Mfa                           struct {
		Configured bool `json:"configured"`
		Exclusion  bool `json:"exclusion"`
	} `json:"mfa"`
	MfaEnrollment struct {
		OverallStatus  string `json:"overallStatus"`
		PushStatus     string `json:"pushStatus"`
		TotpStatus     string `json:"totpStatus"`
		WebAuthnStatus string `json:"webAuthnStatus"`
	} `json:"mfaEnrollment"`
	PasswordDate         time.Time     `json:"password_date"`
	PasswordExpired      bool          `json:"password_expired"`
	PasswordNeverExpires bool          `json:"password_never_expires"`
	PasswordlessSudo     bool          `json:"passwordless_sudo"`
	PhoneNumbers         []interface{} `json:"phoneNumbers"`
	RestrictedFields     []interface{} `json:"restrictedFields"`
	SambaServiceUser     bool          `json:"samba_service_user"`
	SSHKeys              []interface{} `json:"ssh_keys"`
	State                string        `json:"state"`
	Sudo                 bool          `json:"sudo"`
	Suspended            bool          `json:"suspended"`
	TotpEnabled          bool          `json:"totp_enabled"`
	UnixGUID             int           `json:"unix_guid"`
	UnixUID              int           `json:"unix_uid"`
}
