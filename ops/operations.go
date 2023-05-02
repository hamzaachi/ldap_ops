package ops

import (
	"fmt"
	"log"

	"github.com/go-ldap/ldap"
)

const (
	ldapServer   = "ldap://1.1.1.1:389"
	ldapBaseDN   = "ou=personnel,dc=blah,dc=dz"
	BindUsername = "cn=ldapadmin,dc=blah,dc=dz"
	BindPassword = "balhablah"
)

func Connect() (*ldap.Conn, error) {

	conn, err := ldap.DialURL(ldapServer)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func Search(Conn *ldap.Conn, Filter string) (*ldap.SearchResult, error) {
	err := Conn.Bind(BindUsername, BindPassword)
	if err != nil {
		return nil, err
	}
	searchRequest := ldap.NewSearchRequest(
		ldapBaseDN, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		Filter,     // The filter to apply
		[]string{}, // A list attributes to retrieve
		nil,
	)

	result, err := Conn.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("Search Error: %s", err)
	}

	if len(result.Entries) > 0 {
		return result, nil
	} else {
		return nil, fmt.Errorf("Couldn't fetch search entries")
	}
}

func PasswordModifyAdmin(Conn *ldap.Conn, username, password string) error {
	//_, err := Search(Conn, "(uid=hamza.achi)")
	//if err != nil {
	//	return fmt.Errorf(err.Error())
	//}

	err := Conn.Bind(BindUsername, BindPassword)
	if err != nil {
		return err
	}

	username = "uid=" + username + "," + ldapBaseDN
	passwordModifyRequest := ldap.NewPasswordModifyRequest(username, "", password)
	_, err = Conn.PasswordModify(passwordModifyRequest)

	defer Conn.Close()
	if err != nil {
		return fmt.Errorf("Password could not be changed: %s", err.Error())
	}
	return nil
}

func AddUser(Conn *ldap.Conn, username, password string) error {
	err := Conn.Bind(BindUsername, BindPassword)
	if err != nil {
		return err
	}

	dn := fmt.Sprintf("uid=%s,%s", username, ldapBaseDN)
	//username = "uid=" + username + "," + ldapBaseDN
	request := ldap.NewAddRequest(dn, []ldap.Control{})
	request.Attribute("objectClass", []string{"inetOrgPerson", "shadowAccount", "posixAccount", "top"})
	request.Attribute("cn", []string{username})
	request.Attribute("gidNumber", []string{"100"})
	request.Attribute("homeDirectory", []string{fmt.Sprintf("/home/%s", username)})
	request.Attribute("sn", []string{username})
	request.Attribute("uid", []string{username})
	request.Attribute("loginShell", []string{"/bin/bash"})
	request.Attribute("uidNumber", []string{"60811"})
	request.Attribute("mail", []string{fmt.Sprintf("%s@eadn.dz", username)})
	request.Attribute("shadowLastChange", []string{"17058"})
	request.Attribute("shadowMax", []string{"99999"})
	request.Attribute("shadowMin", []string{"0"})
	request.Attribute("userPassword", []string{password})

	if err := Conn.Add(request); err != nil {
		log.Fatal("error adding service:", request, err)
	}

	defer Conn.Close()
	return nil
}

func DelUser(Conn *ldap.Conn, username string) error {
	err := Conn.Bind(BindUsername, BindPassword)
	if err != nil {
		return err
	}

	dn := fmt.Sprintf("uid=%s,%s", username, ldapBaseDN)
	delReq := ldap.NewDelRequest(dn, []ldap.Control{})

	if err := Conn.Del(delReq); err != nil {
		log.Fatalf("Error deleting service: %v", err)
	}
	return nil
}
