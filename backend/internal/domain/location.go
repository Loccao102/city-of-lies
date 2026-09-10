package domain

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// Location represents a gameplay district anchor point.
type Location struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Description        string   `json:"description"`
	Position           Position `json:"position"`
	ConnectedLocations []string `json:"connected_locations"`
}
