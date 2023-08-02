package hubspotclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Contact struct {
	ID         string     `json:"id"`
	Properties Properties `json:"properties"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

type Properties struct {
	FirstName             string `json:"firstname"`
	LastName              string `json:"lastname"`
	Email                 string `json:"email"`
	Company               string `json:"company"`
	Owner                 string `json:"hubspot_owner_id"`
	Phone                 string `json:"phone"`
	DeveloperSetup        string `json:"developer_setup"`
	Domain                string `json:"hs_email_domain"`
	Title                 string `json:"jobtitle"`
	OriginalSource        string `json:"hs_analytics_source"`
	OriginalSourceData2   string `json:"hs_analytics_source_data_2"`
	CreatedAccount        string `json:"created_account"`
	CreatedOrganization   string `json:"created_organization"`
	InvitedToOrganization string `json:"invited_to"`
	FirstConversion       string `json:"first_conversion_event_name"`
	RecentConversion      string `json:"recent_conversion_event_name"`
}

//	Event source- if apply- hs_analytics_source_data_2

type Paging struct {
	Next struct {
		After string `json:"after"`
	} `json:"next"`
}

type ContactResponse struct {
	Results []Contact `json:"results"`
	Paging  Paging    `json:"paging"`
}

func (c *HubspotClient) ListContacts() ([]Contact, error) {
	contacts := make([]Contact, 0)

	properties := []string{"firstname",
		"lastname",
		"email",
		"company",
		"hubspot_owner_id",
		"phone",
		"developer_setup",
		"hs_email_domain",
		"jobtitle",
		"hs_analytics_source",
		"hs_analytics_source_data_2",
		"created_account",
		"created_organization",
		"invited_to",
		"first_conversion_event_name",
		"recent_conversion_event_name"}

	params := strings.Join(properties, "&properties=")
	after := ""
	for {
		req, err := http.NewRequest("GET", fmt.Sprintf("https://api.hubapi.com/crm/v3/objects/contacts?limit=100&properties=%s%s", params, after), nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
		req.Header.Set("Accept", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var contactResponse ContactResponse
		err = json.Unmarshal(body, &contactResponse)
		if err != nil {
			return nil, err
		}

		contacts = append(contacts, contactResponse.Results...)

		if contactResponse.Paging.Next.After == "" {
			break
		}
		after = fmt.Sprintf("&after=%s", contactResponse.Paging.Next.After)
	}

	return contacts, nil
}
