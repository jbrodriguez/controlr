package dto

// Wake -
type Wake struct {
	Mac       string `json:"mac"`
	Broadcast string `json:"broadcast"`
}

// Prefs -
type Prefs struct {
	Number string `json:"number"`
	Unit   string `json:"unit"`
}

// Info -
type Info struct {
	Wake  Wake  `json:"wake"`
	Prefs Prefs `json:"prefs"`
}
