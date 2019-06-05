package models

type Channel struct {
	UID      string  `json:"uid,omitempty"`
	Node     []*Node `json:"node,omitempty"`
	Capacity int     `json:"capacity,omitempty"`
}
