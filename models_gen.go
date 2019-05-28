// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package storm

type Channel struct {
	UID      string `json:"uid"`
	Node1    *Node  `json:"node1"`
	Node2    *Node  `json:"node2"`
	Capacity int    `json:"capacity"`
}

type NewChannel struct {
	Node1    string `json:"node1"`
	Node2    string `json:"node2"`
	Capacity int    `json:"capacity"`
}

type NewNode struct {
	PubKey string `json:"pubKey"`
	Name   string `json:"name"`
}

type Node struct {
	UID      string     `json:"uid"`
	PubKey   string     `json:"pubKey"`
	Name     string     `json:"name"`
	Channels []*Channel `json:"channels"`
}