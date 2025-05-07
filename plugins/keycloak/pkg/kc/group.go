package kc

//nolint:tagliatelle // maintain json user formatting
type Group struct {
	Type          string `json:"type"`
	ID            string `json:"id"`
	Name          string `json:"name"`
	Path          string `json:"path"`
	SubGroupCount int    `json:"subGroupCount"`
	SubGroups     []any  `json:"subGroups"`
	Access        struct {
		View             bool `json:"view"`
		ViewMembers      bool `json:"viewMembers"`
		ManageMembers    bool `json:"manageMembers"`
		Manage           bool `json:"manage"`
		ManageMembership bool `json:"manageMembership"`
	} `json:"access"`
}
