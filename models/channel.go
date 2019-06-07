package models

type Channel struct {
	UID  string `json:"uid,omitempty"`
	Node []struct {
		UID string `json:"uid,omitempty"`
	} `json:"node,omitempty"`
	Capacity int `json:"capacity,omitempty"`
}
