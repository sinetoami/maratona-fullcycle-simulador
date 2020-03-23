package entity

// Order type
type Order struct {
	UUID        string `json:"order"`
	Destination string `json:"destination"`
}

// Destination type
type Destination struct {
	Order string `json:"order"`
	Lat   string `json:"lat"`
	Lng   string `json:"lng"`
}
