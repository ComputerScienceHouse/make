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

	if bindDN == "" || password == "" {
		return nil, errors.New("LDAP credentials not configured")
	}

	if err := conn.Bind(bindDN, password); err != nil {
		return nil, err
	}

	return conn, nil
}

func AddUserToGroup(username string, group string) error {
	conn, err := createLDAPConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	userDN := fmt.Sprintf(
		"uid=%s,cn=users,cn=accounts,dc=csh,dc=rit,dc=edu",
		username,
	)

	groupDN := fmt.Sprintf(
		"cn=%s,cn=groups,cn=accounts,dc=csh,dc=rit,dc=edu",
		ldap.EscapeDN(group),
	)

	req := ldap.NewModifyRequest(groupDN, nil)
	req.Add("member", []string{userDN})

	return conn.Modify(req)

}

func RemoveUserFromGroup(username string, group string) error {
	conn, err := createLDAPConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	userDN := fmt.Sprintf(
		"uid=%s,cn=users,cn=accounts,dc=csh,dc=rit,dc=edu",
		username,
	)

	groupDN := fmt.Sprintf(
		"cn=%s,cn=groups,cn=accounts,dc=csh,dc=rit,dc=edu",
		ldap.EscapeDN(group),
	)

	req := ldap.NewModifyRequest(groupDN, nil)
	req.Delete("member", []string{userDN})

	return conn.Modify(req)

}

type UserWUUID struct {
	Username string `json:"username"`
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
}

func getUserFromDN(conn *ldap.Conn, dn string) (UserWUUID, error) {
	search := ldap.NewSearchRequest(
		dn,
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		1,
		0,
		false,
		"(objectClass=*)",
		[]string{"uid", "cn", "ipaUniqueID"},
		nil,
	)

	res, err := conn.Search(search)
	if err != nil {
		return UserWUUID{}, err
	}

	if len(res.Entries) == 0 {
		return UserWUUID{}, fmt.Errorf("user not found: %s", dn)
	}

	e := res.Entries[0]

	return UserWUUID{
		Username: e.GetAttributeValue("uid"),
		UUID:     e.GetAttributeValue("ipaUniqueID"),
		Name:     e.GetAttributeValue("cn"),
	}, nil
}

// only used by worker
func GetDirectGroupMembers(group string) ([]UserWUUID, error) {
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
		groupDN,
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		1,
		0,
		false,
		"(objectClass=*)",
		[]string{"member"},
		nil,
	)

	res, err := conn.Search(search)
	if err != nil {
		return []UserWUUID{}, err
	}

	if len(res.Entries) == 0 {
		return nil, fmt.Errorf("group not found: %s", groupDN)
	}

	users := []UserWUUID{}

	members := res.Entries[0].GetAttributeValues("member")
	for _, dn := range members {
		if !strings.HasPrefix(dn, "uid=") {
			log.Printf("[WARN] [WORKER] There are external groups assigned as members of the %s group, make is not single source of truth!", group)
			continue
		}

		u, err := getUserFromDN(conn, dn)
		if err != nil {
			return []UserWUUID{}, err
		}

		users = append(users, u)
	}

	return users, nil
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

func UUIDtoUsername(uuid string) (string, error) {
	conn, err := createLDAPConnection()
	if err != nil {
		return "", err
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
		[]string{"uid"},
		nil,
	)

	result, err := conn.Search(search)
	if err != nil {
		return "", err
	}

	if len(result.Entries) == 0 {
		return "", fmt.Errorf("user with uuid %s not found", uuid)
	}

	entry := result.Entries[0]

	return entry.GetAttributeValue("uid"), nil
}
