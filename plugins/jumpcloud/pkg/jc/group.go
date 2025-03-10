package jc

type Group struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Email       string      `json:"email"`
	Attributes  interface{} `json:"attributes"`
}
