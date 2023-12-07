package ldapclient

import (
	"crypto/tls"
	"fmt"
	"log"

	ldap "github.com/go-ldap/ldap/v3"
)

type Attribute struct {
	name   string
	values []string
}

type LDAPClient struct {
	ldapConn *ldap.Conn
}

func NewLDAPClient(user, password, host string) (*LDAPClient, error) {
	ldapClient := &LDAPClient{}

	ldapConn, err := ldapClient.initLDAPConnection(user, password, host)
	if err != nil {
		return nil, err
	}
	ldapClient.ldapConn = ldapConn

	return ldapClient, nil
}

func (l *LDAPClient) initLDAPConnection(username, password, host string) (*ldap.Conn, error) {
	//ldapConn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", host, 1636))

	//ldapConn, err := ldap.DialURL(fmt.Sprintf("ldap://%s:%d", "127.0.0.1", 1389), ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
	ldapConn, err := ldap.DialURL(fmt.Sprintf("ldap://%s:%d", "127.0.0.1", 1389), ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
	if err != nil {
		return nil, err
	}

	err = ldapConn.Bind(username, password)
	if err != nil {
		return nil, err
	}

	return ldapConn, nil
}

func (l *LDAPClient) Close() {
	l.ldapConn.Close()
}

func (l *LDAPClient) ListUsers() []*ldap.Entry {
	//baseDN := "dc=trial-1033352, dc=okta, dc=com"
	baseDN := "dc=example,dc=org"
	filter := "(&(objectClass=organizationalPerson))"
	attributes := []string{}

	// Search for all users
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		attributes,
		nil,
	)
	sr, err := l.ldapConn.SearchWithPaging(searchRequest, 1000)
	if err != nil {
		log.Fatal(err)
	}

	return sr.Entries
}

type Entry struct {
	DN         string
	Attributes map[string][]string
}

func (l *LDAPClient) GetUserByID(id string) {

}

func (l *LDAPClient) GetUserByEmail(email string) {

}
