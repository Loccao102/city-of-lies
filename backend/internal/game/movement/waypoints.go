package movement

import (
	"city-of-lies/backend/internal/domain"
)

type WaypointGraph struct {
	locations map[string]domain.Location
}

func NewWaypointGraph(locations []domain.Location) *WaypointGraph {
	g := &WaypointGraph{
		locations: make(map[string]domain.Location),
	}
	for _, loc := range locations {
		g.locations[loc.ID] = loc
	}
	return g
}

func (g *WaypointGraph) GetLocation(id string) (domain.Location, bool) {
	loc, ok := g.locations[id]
	return loc, ok
}

func (g *WaypointGraph) AreConnected(fromID, toID string) bool {
	from, ok := g.locations[fromID]
	if !ok {
		return false
	}
	for _, neighbor := range from.ConnectedLocations {
		if neighbor == toID {
			return true
		}
	}
	return false
}
