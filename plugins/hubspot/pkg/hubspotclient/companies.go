package hubspotclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Company struct {
	ID         string            `json:"id"`
	CreatedAt  time.Time         `json:"createdAt"`
	UpdatedAt  time.Time         `json:"updatedAt"`
	Properties CompanyProperties `json:"properties"`
}

type CompanyProperties struct {
	City     string `json:"city"`
	Domain   string `json:"domain"`
	Industry string `json:"industry"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	State    string `json:"state"`
}

type CompaniesResponse struct {
	Results []Company `json:"results"`
	Paging  Paging    `json:"paging"`
}

func (c *HubspotClient) ListCompanies() ([]Company, error) {
	companies := make([]Company, 0)

	params := "&properties=domain&properties=name&properties=city&properties=industry&properties=phone&properties=state&archived=false"
	after := ""
	for {
		req, err := http.NewRequest("GET", fmt.Sprintf("https://api.hubapi.com/crm/v3/objects/companies?limit=100%s%s", params, after), nil)
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

		var response CompaniesResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}

		companies = append(companies, response.Results...)

		if response.Paging.Next.After == "" {
			break
		}
		after = fmt.Sprintf("&after=%s", response.Paging.Next.After)
	}

	// populate the map from company names to companies
	for _, company := range companies {
		c.companies[company.Properties.Name] = company
	}

	return companies, nil
}

func (c *HubspotClient) LookupCompany(name string) Company {
	return c.companies[name]
}
