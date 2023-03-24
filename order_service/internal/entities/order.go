package entities

type Order struct {
	ID        uint64     `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	FleetTask *FleetTask `json:"fleet_task,omitempty"`
}
