package models

type Channel struct {
	UID      string `json:"uid,omitempty"`
	Node1    string `json:"node1,omitempty"`
	Node2    string `json:"node2,omitempty"`
	Capacity int    `json:"capacity,omitempty"`
}
