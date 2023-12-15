package fetch

import (
	"context"
	"io"
	"strings"

	"github.com/aserto-dev/ds-load/plugins/ldap/pkg/attribute"
	"github.com/aserto-dev/ds-load/plugins/ldap/pkg/ldapclient"
	"github.com/aserto-dev/ds-load/sdk/common/js"
	"github.com/go-ldap/ldap/v3"
)

type Fetcher struct {
	ldapClient *ldapclient.LDAPClient
}

type Entry struct {
	EntryType  string
	DN         string
	Attributes map[string][]string
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

	groups := f.ldapClient.ListGroups()
	err = writeEntries(groups, jsonWriter, "group")
	if err != nil {
		return err
	}

	users := f.ldapClient.ListUsers()
	return writeEntries(users, jsonWriter, "user")
}

func writeEntries(ldapEntries []*ldap.Entry, jsonWriter *js.JSONArrayWriter, entryType string) error {
	distinguishedNames := extractDNs(ldapEntries)

	for _, ldapEntry := range ldapEntries {
		entry := Entry{
			EntryType:  entryType,
			DN:         normalizeDN(ldapEntry.DN),
			Attributes: attribute.Transform(ldapEntry, distinguishedNames, entryType),
		}

		err := jsonWriter.Write(entry)
		if err != nil {
			return err
		}
	}

	return nil
}

func normalizeDN(dn string) string {
	return strings.ReplaceAll(dn, " ", "")
}

func extractDNs(ldapEntries []*ldap.Entry) []string {
	var distinguishedNames []string
	for _, entry := range ldapEntries {
		distinguishedNames = append(distinguishedNames, normalizeDN(entry.DN))
	}
	return distinguishedNames
}
