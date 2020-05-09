package model

// State -
type State struct {
	Name         string
	Timezone     string
	Version      string
	CsrfToken    string
	Host         string
	Origin       string
	Secure       bool
	Cert         string
	UseSelfCerts bool
}
