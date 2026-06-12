package models

type Area struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	LdapGroup   *string `json:"ldapGroup"`
	PhotoURL    string  `json:"photourl"`
}

type CreateAreaRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	PhotoURL    string  `json:"photourl"`
	LdapGroup   *string `json:"ldapGroup"`
}

type AreaWithAccess struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PhotoURL    string `json:"photourl"`
	LdapGroup   string `json:"ldapGroup"`
	HasAccess   bool   `json:"hasAccess"`
}
