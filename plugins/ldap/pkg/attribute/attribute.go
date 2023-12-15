package attribute

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/go-objectsid"
	"github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"
)

func Transform(ldapEntry *ldap.Entry, userDnTOKey, groupDnTOKey map[string]string) map[string][]string {
	if _, ok := groupDnTOKey[ldapEntry.DN]; ok {
		transformedAttributes := decodeAttributes(ldapEntry)
		return addMembersByType(transformedAttributes, userDnTOKey, groupDnTOKey)
	} else {
		return decodeAttributes(ldapEntry)
	}
}

func addMembersByType(attributes map[string][]string, userDnTOKey, groupDnTOKey map[string]string) map[string][]string {
	if attributes["member"] == nil {
		return attributes
	}

	for _, member := range attributes["member"] {
		if groupKey, ok := groupDnTOKey[member]; ok {
			attributes["memberGroup"] = append(attributes["memberGroup"], groupKey)
			continue
		}

		if userKey, ok := userDnTOKey[member]; ok {
			attributes["memberUser"] = append(attributes["memberUser"], userKey)
		}
	}

	return attributes
}

func decodeAttributes(ldapEntry *ldap.Entry) map[string][]string {
	var data = make(map[string][]string)
	for _, attribute := range ldapEntry.Attributes {
		if attribute.Name == "objectSid" {
			data[attribute.Name] = []string{ObjectSid(ldapEntry)}
			continue
		}

		if attribute.Name == "objectGUID" {
			data[attribute.Name] = []string{ObjectGUID(ldapEntry)}
			continue
		}

		data[attribute.Name] = attribute.Values
	}
	return data
}

/*
/* Extract the ObjectSid. This is specific to Active Directory.
*/
func ObjectSid(entry *ldap.Entry) string {
	rawObjectSid := entry.GetRawAttributeValue("objectSid")
	if len(rawObjectSid) > 0 {
		return objectsid.Decode(rawObjectSid).String()
	}
	return ""
}

/*
/* Extract the ObjectGUID. This is specific to Active Directory.
*/
func ObjectGUID(entry *ldap.Entry) string {
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
/* The GUID is formed by splitting the octal string in groups of 2 symbols and then re-arrange the groups.
/* {3}{2}{1}{0}-{5}{4}-{7}{6}-{8}{9}-{10}{11}{12}{13}{14}{15} .
*/
func uuidToComStyle(token string) string {
	token = strings.ReplaceAll(token, "-", "")
	return fmt.Sprintf("%s%s%s%s-%s%s-%s%s-%s%s-%s%s%s%s%s%s",
		token[6:8], token[4:6], token[2:4], token[0:2], token[10:12], token[8:10], token[14:16], token[12:14],
		token[16:18], token[18:20], token[20:22], token[22:24], token[24:26], token[26:28], token[28:30], token[30:32])
}
