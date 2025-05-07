package kc

const TypeGroup string = "group"

type Group struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Description string      `json:"description,omitempty"`
	Email       string      `json:"email,omitempty"`
	Attributes  interface{} `json:"attributes,omitempty"`
}
