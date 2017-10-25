package dto

// Info -
// No version || version 0
// Wake    Wake  `json:"wake"`
// Prefs   Prefs `json:"prefs"`

// Version 1
// Previous +
// Version int
// Ups []Sample
type Info struct {
	Version int      `json:"version"`
	Wake    Wake     `json:"wake"`
	Prefs   Prefs    `json:"prefs"`
	Samples []Sample `json:"samples"`
}
