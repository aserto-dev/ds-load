package fetch

import (
	"context"
	"io"
	"strings"

	"github.com/aserto-dev/ds-load/plugins/ldap/pkg/ldapclient"
	"github.com/aserto-dev/ds-load/sdk/common/js"
	"github.com/bwmarrin/go-objectsid"
	"github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
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
			Attributes: attributes(ldapEntry, distinguishedNames, entryType),
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

func attributes(ldapEntry *ldap.Entry, distinguishedNames []string, entryType string) map[string][]string {
	if entryType == "group" {
		transformedAttributes := transformAttributes(ldapEntry)
		return addMembersByType(transformedAttributes, distinguishedNames)
	} else {
		return transformAttributes(ldapEntry)
	}
}

func transformAttributes(ldapEntry *ldap.Entry) map[string][]string {
	var data = make(map[string][]string)
	for _, attribute := range ldapEntry.Attributes {
		if attribute.Name == "objectSid" {
			data[attribute.Name] = []string{getObjectSid(ldapEntry)}
			continue
		}

		if attribute.Name == "objectGUID" {
			data[attribute.Name] = []string{getObjectGUID(ldapEntry)}
			continue
		}

		data[attribute.Name] = attribute.Values
	}
	return data
}

func getObjectSid(entry *ldap.Entry) string {
	rawObjectSid := entry.GetRawAttributeValue("objectSid")
	if len(rawObjectSid) > 0 {
		return objectsid.Decode(rawObjectSid).String()
	}
	return ""
}

func getObjectGUID(entry *ldap.Entry) string {
	rawObjectGUID := entry.GetRawAttributeValue("objectGUID")
	if len(rawObjectGUID) > 0 {
		objectGUID := entry.GetRawAttributeValue("objectGUID")
		u, err := uuid.FromBytes(objectGUID)
		if err != nil {
			return ""
		}
		return u.String()
	}
	return ""
}

func addMembersByType(attributes map[string][]string, groupDNs []string) map[string][]string {
	if attributes["member"] != nil {
		for _, member := range attributes["member"] {
			if slices.Contains(groupDNs, normalizeDN(member)) {
				attributes["memberGroup"] = append(attributes["memberGroup"], normalizeDN(member))
			} else {
				attributes["memberUser"] = append(attributes["memberUser"], normalizeDN(member))
			}
		}
	}

	return attributes
}
