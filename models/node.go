package models

type Node struct {
	UID     string `json:"uid,omitempty"`
	PubKey  string `json:"pubKey,omitempty"`
	Name    string `json:"name,omitempty"`
	Channel []struct {
		UID string `json:"uid,omitempty"`
	} `json:"channel,omitempty"`
}
