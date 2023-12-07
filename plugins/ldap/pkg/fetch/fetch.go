package fetch

import (
	"context"
	"io"

	"github.com/aserto-dev/ds-load/plugins/ldap/pkg/ldapclient"
	"github.com/aserto-dev/ds-load/sdk/common/js"
	"github.com/go-ldap/ldap/v3"
)

type Fetcher struct {
	ldapClient *ldapclient.LDAPClient
}

type Entry struct {
	DN         string
	Attributes map[string][]string
}

type Attribute struct {
	Name   string
	Values []string
}

func New(ldapClient *ldapclient.LDAPClient) (*Fetcher, error) {
	return &Fetcher{
		ldapClient: ldapClient,
	}, nil
}

func (f *Fetcher) Fetch(ctx context.Context, outputWriter, errorWriter io.Writer) error {
	jsonWriter, err := js.NewJSONArrayWriter(outputWriter)
	if err != nil {
		return err
	}
	defer jsonWriter.Close()

	ldapEntries := f.ldapClient.ListUsers()

	for _, ldapEntry := range ldapEntries {
		entry := Entry{
			DN:         ldapEntry.DN,
			Attributes: transformAttributes(ldapEntry.Attributes),
		}

		err = jsonWriter.Write(entry)
		if err != nil {
			_, _ = errorWriter.Write([]byte(err.Error()))
		}
	}

	return nil
}

func transformToAtTributesArray(attributes []*ldap.EntryAttribute) []*Attribute {
	var data = make([]*Attribute, 0)
	for _, attribute := range attributes {
		data = append(data, &Attribute{
			Name:   attribute.Name,
			Values: attribute.Values,
		})
	}
	return data
}

func transformAttributes(attributes []*ldap.EntryAttribute) map[string][]string {
	var data = make(map[string][]string)
	for _, attribute := range attributes {
		data[attribute.Name] = attribute.Values
	}
	return data
}
