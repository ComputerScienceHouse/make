package ldap

import (
	"errors"
	"fmt"
	"log"
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

	if err := conn.Bind(bindDN, password); err != nil {
		log.Fatalf("bind failed: %v", err)
	}

	return conn, nil
}

type UserWUUID struct {
	Username string `json:"username"`
	UUID     string `json:"uuid"`
}

func GetGroupMembers(group string) ([]UserWUUID, error) {
	conn, err := createLDAPConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	groupSearch := ldap.NewSearchRequest(
		"cn=groups,cn=accounts,dc=csh,dc=rit,dc=edu",
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		1,
		0,
		false,
		fmt.Sprintf("(cn=%s)", ldap.EscapeFilter(group)),
		[]string{"member"},
		nil,
	)

	res, err := conn.Search(groupSearch)
	if err != nil {
		return nil, err
	}

	if len(res.Entries) == 0 {
		return nil, errors.New("group not found")
	}

	var users []UserWUUID

	for _, memberDN := range res.Entries[0].GetAttributeValues("member") {

		userSearch := ldap.NewSearchRequest(
			memberDN,
			ldap.ScopeBaseObject,
			ldap.NeverDerefAliases,
			1,
			0,
			false,
			"(objectClass=*)",
			[]string{"uid", "ipaUniqueID"},
			nil,
		)

		uRes, err := conn.Search(userSearch)
		if err != nil || len(uRes.Entries) == 0 {
			continue
		}

		e := uRes.Entries[0]

		users = append(users, UserWUUID{
			Username: e.GetAttributeValue("uid"),
			UUID:     e.GetAttributeValue("ipaUniqueID"),
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
