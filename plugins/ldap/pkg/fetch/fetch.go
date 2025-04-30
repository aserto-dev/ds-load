package fetch

import (
	"context"
	"io"

	"github.com/aserto-dev/ds-load/plugins/ldap/pkg/attribute"
	"github.com/aserto-dev/ds-load/plugins/ldap/pkg/ldapclient"
	"github.com/aserto-dev/ds-load/sdk/common"
	"github.com/aserto-dev/ds-load/sdk/common/js"
	"github.com/go-ldap/ldap/v3"
)

type Fetcher struct {
	ldapClient *ldapclient.LDAPClient
	idField    string
}

type Entry struct {
	EntryType  string
	Key        string
	Attributes map[string][]string
}

func New(ldapClient *ldapclient.LDAPClient, idField string) *Fetcher {
	return &Fetcher{
		ldapClient: ldapClient,
		idField:    idField,
	}
}

func (f *Fetcher) Fetch(ctx context.Context, outputWriter io.Writer, errorWriter common.ErrorWriter) error {
	jsonWriter := js.NewJSONArrayWriter(outputWriter)
	defer jsonWriter.Close()

	groups := f.ldapClient.ListGroups()
	users := f.ldapClient.ListUsers()
	userDnTOKey := buildMapFromDNToKey(users, f.idField)
	groupDnTOKey := buildMapFromDNToKey(groups, f.idField)

	err := writeEntries(groups, jsonWriter, userDnTOKey, groupDnTOKey)
	if err != nil {
		return err
	}

	return writeEntries(users, jsonWriter, userDnTOKey, groupDnTOKey)
}

func writeEntries(ldapEntries []*ldap.Entry, jsonWriter *js.JSONArrayWriter, userDnTOKey, groupDnTOKey map[string]string) error {
	dnToKey := userDnTOKey
	for k, v := range groupDnTOKey {
		dnToKey[k] = v
	}

	for _, ldapEntry := range ldapEntries {
		entry := Entry{
			EntryType:  entryType(ldapEntry, groupDnTOKey),
			Key:        dnToKey[ldapEntry.DN],
			Attributes: attribute.Transform(ldapEntry, userDnTOKey, groupDnTOKey),
		}

		err := jsonWriter.Write(entry)
		if err != nil {
			return err
		}
	}

	return nil
}

func entryType(ldapEntry *ldap.Entry, groupDnToKey map[string]string) string {
	if _, ok := groupDnToKey[ldapEntry.DN]; ok {
		return "group"
	} else {
		return "user"
	}
}

func buildMapFromDNToKey(ldapEntries []*ldap.Entry, key string) map[string]string {
	mapDNToKey := make(map[string]string)
	for _, entry := range ldapEntries {
		mapDNToKey[entry.DN] = extractKey(key, entry)
	}

	return mapDNToKey
}

func extractKey(key string, entry *ldap.Entry) string {
	if key == "objectGUID" {
		return attribute.ObjectGUID(entry)
	}

	if key == "objectSid" {
		return attribute.ObjectSid(entry)
	}

	return entry.GetAttributeValue(key)
}
