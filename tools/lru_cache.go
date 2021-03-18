package tools

import "container/ring"

// LRUCache not implemented
type LRUCache struct {
	maxSize int
	used    int
	store   map[string]ring.Ring
}
