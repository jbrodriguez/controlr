package model

import "plugin/dto"

// State -
type State struct {
	Name         string
	Timezone     string
	Version      string
	CsrfToken    string
	Host         string
	Origin       dto.Origin
	Secure       bool
	Cert         string
	UseSelfCerts bool
}
