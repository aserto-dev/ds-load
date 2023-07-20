package hubspotclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type UsersResponse struct {
	Results []User `json:"results"`
	Paging  Paging `json:"paging"`
}

func (c *HubspotClient) ListUsers() ([]User, error) {
	users := make([]User, 0)

	after := ""
	for {
		req, err := http.NewRequest("GET", fmt.Sprintf("https://api.hubapi.com/crm/v3/owners?limit=100%s", after), nil)
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

		var response UsersResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}

		users = append(users, response.Results...)

		if response.Paging.Next.After == "" {
			break
		}
		after = fmt.Sprintf("&after=%s", response.Paging.Next.After)
	}

	// populate the map from user IDs to emails
	for _, user := range users {
		c.users[user.ID] = user.Email
	}

	return users, nil
}

func (c *HubspotClient) LookupUser(id string) string {
	return c.users[id]
}
