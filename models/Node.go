package models

import (
	"errors"
)

// Node represents a node in the graph
type Node struct {
	Text string
	ID   int64
}

// NodeSet represents a set of nodes
type NodeSet struct {
	set map[Node]bool
}

// NewNodeSet constructs an empty node set
func NewNodeSet() NodeSet {
	return NodeSet{
		set: make(map[Node]bool),
	}
}

// Add adds a node to the set
func (ns *NodeSet) Add(node Node) {
	ns.set[node] = true
}

// GetByText gets a node from the set by text,
// returns error if node does not exist
func (ns *NodeSet) GetByText(text string) (Node, error) {
	for k := range ns.set {
		if k.Text == text {
			return k, nil
		}
	}
	return Node{}, errors.New("No Node Found")
}

// GetByID gets a node from the set by ID,
// returns error if node does not exist
func (ns *NodeSet) GetByID(id int64) (Node, error) {
	for k := range ns.set {
		if k.ID == id {
			return k, nil
		}
	}
	return Node{}, errors.New("No Node Found")
}

// RemoveID removes a node by the Id
func (ns *NodeSet) RemoveID(id int64) {
	node, err := ns.GetByID(id)
	if err != nil {
		return
	}
	delete(ns.set, node)
}

// RemoveText removes a node by text
func (ns *NodeSet) RemoveText(text string) {
	node, err := ns.GetByText(text)
	if err != nil {
		return
	}
	delete(ns.set, node)
}

// ContainsID returns true if the set contains a node with the ID
func (ns *NodeSet) ContainsID(id int64) bool {
	_, err := ns.GetByID(id)
	if err != nil {
		return false
	}
	return true
}

// ContainsText returns true if the set contains a node with the text
func (ns *NodeSet) ContainsText(text string) bool {
	_, err := ns.GetByText(text)
	if err != nil {
		return false
	}
	return true
}
