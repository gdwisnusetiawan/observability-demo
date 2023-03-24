package entities

type FleetTask struct {
	ID      uint64  `json:"id,omitempty"`
	Name    string  `json:"name,omitempty"`
	Vehicle Vehicle `json:"vehicle,omitempty"`
}

type Vehicle struct {
	ID           uint64 `json:"id,omitempty"`
	LisencePlate string `json:"lisence_plate,omitempty"`
}
