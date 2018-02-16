package dto

// Info -
// No version || version 0
// Wake    Wake  `json:"wake"`
// Prefs   Prefs `json:"prefs"`
//
// Version 1
// Previous +
// Version int
// Ups []Sample
//
// Version 2
// + Available
type Info struct {
	Version  int             `json:"version"`
	Wake     Wake            `json:"wake"`
	Prefs    Prefs           `json:"prefs"`
	Samples  []Sample        `json:"samples"`
	Features map[string]bool `json:"features"`
}
