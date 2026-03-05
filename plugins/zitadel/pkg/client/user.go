package client

//nolint:tagliatelle // maintain json user formatting
type User struct {
	Type                       string `json:"type"`
	ID                         string `json:"id"`
	Username                   string `json:"username"`
	FirstName                  string `json:"firstName"`
	LastName                   string `json:"lastName"`
	Email                      string `json:"email"`
	EmailVerified              bool   `json:"emailVerified"`
	CreatedTimestamp           int64  `json:"createdTimestamp"`
	Enabled                    bool   `json:"enabled"`
	Totp                       bool   `json:"totp"`
	DisableableCredentialTypes []any  `json:"disableableCredentialTypes"`
	RequiredActions            []any  `json:"requiredActions"`
	NotBefore                  int    `json:"notBefore"`
	Access                     struct {
		Manage bool `json:"manage"`
	} `json:"access"`
}
