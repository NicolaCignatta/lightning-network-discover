package models

type Channel struct {
	UID  string `json:"uid,omitempty"`
	Node []struct {
		UID string `json:"uid,omitempty"`
	} `json:"node,omitempty"`
	Capacity int `json:"capacity,omitempty"`
}

type ChannelDG struct {
	Channel Channel
	Node    []*Node `json:"node,omitempty"`
}
