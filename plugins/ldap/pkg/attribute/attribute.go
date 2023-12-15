package attribute

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/go-objectsid"
	"github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

func Transform(ldapEntry *ldap.Entry, distinguishedNames []string, entryType string) map[string][]string {
	if entryType == "group" {
		transformedAttributes := decodeAttributes(ldapEntry)
		return addMembersByType(transformedAttributes, distinguishedNames)
	} else {
		return decodeAttributes(ldapEntry)
	}
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

func normalizeDN(dn string) string {
	return strings.ReplaceAll(dn, " ", "")
}

func decodeAttributes(ldapEntry *ldap.Entry) map[string][]string {
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

/*
/* Extract the ObjectSid. This is specific to Active Directory
*/
func getObjectSid(entry *ldap.Entry) string {
	rawObjectSid := entry.GetRawAttributeValue("objectSid")
	if len(rawObjectSid) > 0 {
		return objectsid.Decode(rawObjectSid).String()
	}
	return ""
}

/*
/* Extract the ObjectGUID. This is specific to Active Directory
*/
func getObjectGUID(entry *ldap.Entry) string {
	rawObjectGUID := entry.GetRawAttributeValue("objectGUID")
	if len(rawObjectGUID) > 0 {
		objectGUID := entry.GetRawAttributeValue("objectGUID")
		uuidString, err := uuid.FromBytes(objectGUID)
		if err != nil {
			return ""
		}
		return uuidToComStyle(uuidString.String())
	}
	return ""
}

/*
/* Transform the octal string in the com style GUID from Active Directory.
/* The GUID is formed by splitting the octal string in groups of 2 symbols
/* {3}{2}{1}{0}-{5}{4}-{7}{6}-{8}{9}-{10}{11}{12}{13}{14}{15}
*/
func uuidToComStyle(token string) string {
	token = strings.ReplaceAll(token, "-", "")
	return fmt.Sprintf("%s%s%s%s-%s%s-%s%s-%s%s-%s%s%s%s%s%s",
		token[6:8], token[4:6], token[2:4], token[0:2], token[10:12], token[8:10], token[14:16], token[12:14],
		token[16:18], token[18:20], token[20:22], token[22:24], token[24:26], token[26:28], token[28:30], token[30:32])
}
