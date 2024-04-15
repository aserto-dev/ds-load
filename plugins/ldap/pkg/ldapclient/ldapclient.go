package ldapclient

import (
	"crypto/tls"
	"log"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/rs/zerolog"
)

type LDAPClient struct {
	ldapConn    *ldap.Conn
	credentials *Credentials
	conOptions  *ConnectionOptions
	logger      *zerolog.Logger
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
	IDField     string
}

func NewLDAPClient(credentials *Credentials, conOptions *ConnectionOptions, logger *zerolog.Logger) (*LDAPClient, error) {
	ldapClient := &LDAPClient{}

	ldapConn, err := ldapClient.initLDAPConnection(credentials.User, credentials.Password, conOptions.Host, conOptions.Insecure)
	if err != nil {
		return nil, err
	}
	ldapClient.ldapConn = ldapConn
	ldapClient.credentials = credentials
	ldapClient.conOptions = conOptions
	ldapClient.logger = logger

	return ldapClient, nil
}

func (l *LDAPClient) initLDAPConnection(username, password, host string, insecure bool) (*ldap.Conn, error) {
	var dialOptions []ldap.DialOpt

	// Set default timeout for init connection
	ldap.DefaultTimeout = 10 * time.Second

	// Disable the security check if insecure is true
	if insecure { // #nosec G402
		dialOptions = append(dialOptions, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: insecure}))
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
	err := l.ldapConn.Close()
	if err != nil {
		l.logger.Error().Err(err)
	}
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
