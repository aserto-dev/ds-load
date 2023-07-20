package app

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/ds-load/plugins/hubspot/pkg/hubspotclient"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/plugin"
)

type FetchCmd struct {
	ClientID           string `short:"i" help:"Hubspot Client ID" env:"HUBSPOT_CLIENT_ID"`
	ClientSecret       string `short:"s" help:"Hubspot Client Secret" env:"HUBSPOT_CLIENT_SECRET"`
	RefreshToken       string `short:"r" help:"Hubspot Refresh Token" env:"HUBSPOT_REFRESH_TOKEN"`
	PrivateAccessToken string `short:"p" help:"Hubspot Private Access Token" env:"HUBSPOT_PAT"`
	Contacts           bool   `help:"Retrieve Hubspot contacts" env:"HUBSPOT_CONTACTS" default:"false"`
	Companies          bool   `help:"Retrieve Hubspot companies" env:"HUBSPOT_COMPANIES" default:"false"`
}

func (cmd *FetchCmd) Run(ctx *kong.Context) error {
	var hubspotClient *hubspotclient.HubspotClient
	var err error
	if cmd.PrivateAccessToken != "" {
		hubspotClient, _ = createHubspotClient(cmd.PrivateAccessToken)
	} else {
		hubspotClient, err = createHubspotOAuth2Client(cmd.ClientID, cmd.ClientSecret, cmd.RefreshToken)
		if err != nil {
			return err
		}
	}

	results := make(chan map[string]interface{}, 1)
	errors := make(chan error, 1)
	go func() {
		Fetch(hubspotClient, cmd.Contacts, cmd.Companies, results, errors)
		close(results)
		close(errors)
	}()
	return plugin.NewDSPlugin().WriteFetchOutput(results, errors, false)
}

func Fetch(hubspotClient *hubspotclient.HubspotClient, fetchContacts, fetchCompanies bool, results chan map[string]interface{}, errors chan error) {
	users, err := hubspotClient.ListUsers()
	if err != nil {
		errors <- err
		common.SetExitCode(1)
		return
	}

	for _, user := range users {
		userBytes, err := json.Marshal(user)
		if err != nil {
			errors <- err
			common.SetExitCode(1)
			continue
		}
		var obj map[string]interface{}
		err = json.Unmarshal(userBytes, &obj)
		if err != nil {
			errors <- err
			common.SetExitCode(1)
			continue
		}
		obj["type"] = "user"

		results <- obj
	}

	if fetchCompanies {
		companies, err := hubspotClient.ListCompanies()
		if err != nil {
			errors <- err
			common.SetExitCode(1)
			return
		}

		for _, company := range companies {
			companyBytes, err := json.Marshal(company)
			if err != nil {
				errors <- err
				common.SetExitCode(1)
				continue
			}
			var obj map[string]interface{}
			err = json.Unmarshal(companyBytes, &obj)
			if err != nil {
				errors <- err
				common.SetExitCode(1)
				continue
			}
			obj["type"] = "company"

			results <- obj
		}
	}

	if fetchContacts {
		contacts, err := hubspotClient.ListContacts()
		if err != nil {
			errors <- err
			common.SetExitCode(1)
			return
		}

		for _, contact := range contacts {
			contactBytes, err := json.Marshal(contact)
			if err != nil {
				errors <- err
				common.SetExitCode(1)
				continue
			}
			var obj map[string]interface{}
			err = json.Unmarshal(contactBytes, &obj)
			if err != nil {
				errors <- err
				common.SetExitCode(1)
				continue
			}
			obj["key"] = contact.Properties.Email
			obj["displayName"] = createDisplayName(contact.Properties.FirstName, contact.Properties.LastName, contact.Properties.Email)
			obj["type"] = "contact"
			if contact.Properties.Owner != "" {
				user := hubspotClient.LookupUser(contact.Properties.Owner)
				if user != "" {
					obj["owner"] = user
				}
			}

			if contact.Properties.Company != "" {
				company := hubspotClient.LookupCompany(contact.Properties.Company)
				if company.ID != "" {
					obj["companyId"] = company.ID
				}
			}

			results <- obj
		}
	}
}

func createDisplayName(firstName, lastName, email string) string {
	returnVal := firstName
	if returnVal != "" {
		if lastName != "" {
			returnVal = fmt.Sprintf("%s %s", firstName, lastName)
		}
	} else {
		if lastName != "" {
			returnVal = lastName
		} else {
			returnVal = email
		}
	}
	return returnVal
}

func createHubspotClient(privateAccessToken string) (hubspotClient *hubspotclient.HubspotClient, err error) {
	return hubspotclient.NewHubspotClient(
		context.Background(),
		privateAccessToken)
}

func createHubspotOAuth2Client(clientID, clientSecret, refreshToken string) (hubspotClient *hubspotclient.HubspotClient, err error) {
	return hubspotclient.NewHubspotOAuth2Client(
		context.Background(),
		clientID,
		clientSecret,
		refreshToken)
}
