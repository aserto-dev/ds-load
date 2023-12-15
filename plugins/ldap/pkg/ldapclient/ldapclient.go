package ldapclient

import (
	"crypto/tls"
	"log"

	"github.com/go-ldap/ldap/v3"
)

type LDAPClient struct {
	ldapConn    *ldap.Conn
	credentials *Credentials
	conOptions  *ConnectionOptions
}

type Credentials struct {
	User     string
	Password string
}

type ConnectionOptions struct {
	Host        string
	Insecure    bool
	BaseDN      string
	UserFilter  string
	GroupFilter string
	UuidField   string
}

func NewLDAPClient(credentials *Credentials, conOptions *ConnectionOptions) (*LDAPClient, error) {
	ldapClient := &LDAPClient{}

	ldapConn, err := ldapClient.initLDAPConnection(credentials.User, credentials.Password, conOptions.Host, conOptions.Insecure)
	if err != nil {
		return nil, err
	}
	ldapClient.ldapConn = ldapConn
	ldapClient.credentials = credentials
	ldapClient.conOptions = conOptions

	return ldapClient, nil
}

func (l *LDAPClient) initLDAPConnection(username, password, host string, insecure bool) (*ldap.Conn, error) {
	var dialOptions []ldap.DialOpt

	if insecure {
		dialOptions = append(dialOptions, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
	}

	ldapConn, err := ldap.DialURL(host, dialOptions...)
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
	return l.search(l.conOptions.UserFilter)
}

func (l *LDAPClient) ListGroups() []*ldap.Entry {
	return l.search(l.conOptions.GroupFilter)
}

func (l *LDAPClient) search(filter string) []*ldap.Entry {
	attributes := []string{}

	searchRequest := ldap.NewSearchRequest(
		l.conOptions.BaseDN,
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
