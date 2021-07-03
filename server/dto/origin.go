package dto

// Origin -
type Origin struct {
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Address  string `json:"address"`
}
