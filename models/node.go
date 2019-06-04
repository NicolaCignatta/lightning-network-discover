package models

type Node struct {
	UID      string   `json:"uid,omitempty"`
	PubKey   string   `json:"pubKey,omitempty"`
	Name     string   `json:"name,omitempty"`
	Channels []string `json:"channels,omitempty"`
}
