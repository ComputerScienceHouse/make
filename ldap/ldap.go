package ldap

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

func createLDAPConnection() (*ldap.Conn, error) {
	conn, err := ldap.DialURL("ldap://ldap.csh.rit.edu")
	if err != nil {
		return nil, err
	}

	bindDN := os.Getenv("MAKE_LDAP_BIND_DN")
	password := os.Getenv("MAKE_LDAP_PASS")

	if bindDN == "" || password == "" {
		return nil, errors.New("LDAP credentials not configured")
	}

	if err := conn.Bind(bindDN, password); err != nil {
		return nil, err
	}

	return conn, nil
}

type UserWUUID struct {
	Username string `json:"username"`
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
}

func GetGroupMembers(group string) ([]UserWUUID, error) {
	conn, err := createLDAPConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	groupDN := fmt.Sprintf(
		"cn=%s,cn=groups,cn=accounts,dc=csh,dc=rit,dc=edu",
		ldap.EscapeDN(group),
	)

	search := ldap.NewSearchRequest(
		"cn=users,cn=accounts,dc=csh,dc=rit,dc=edu",
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(memberOf=%s)", ldap.EscapeFilter(groupDN)),
		[]string{"uid", "ipaUniqueID", "cn"},
		nil,
	)

	res, err := conn.Search(search)
	if err != nil {
		return nil, err
	}

	var users []UserWUUID

	for _, entry := range res.Entries {
		users = append(users, UserWUUID{
			Username: entry.GetAttributeValue("uid"),
			UUID:     entry.GetAttributeValue("ipaUniqueID"),
			Name:     entry.GetAttributeValue("cn"),
		})
	}

	return users, nil
}

func GetUserLDAPGroups(uuid string) ([]string, error) {
	conn, err := createLDAPConnection()
	if err != nil {
		return []string{}, err
	}
	defer conn.Close()

	search := ldap.NewSearchRequest(
		"cn=users,cn=accounts,dc=csh,dc=rit,dc=edu",
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		1,
		0,
		false,
		fmt.Sprintf("(ipaUniqueID=%s)", ldap.EscapeFilter(uuid)),
		[]string{"memberOf"},
		nil,
	)

	result, err := conn.Search(search)
	if err != nil {
		return nil, err
	}

	if len(result.Entries) == 0 {
		return nil, fmt.Errorf("user with uuid %s not found", uuid)
	}

	entry := result.Entries[0]

	var groups []string

	for _, dnStr := range entry.GetAttributeValues("memberOf") {
		if !strings.Contains(dnStr, ",cn=groups,cn=accounts,") {
			continue
		}

		dn, err := ldap.ParseDN(dnStr)
		if err != nil {
			continue
		}

		groups = append(groups, dn.RDNs[0].Attributes[0].Value)
	}

	return groups, nil
}
