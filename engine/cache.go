package engine

import (
	"sync"

	"github.com/CloudyKit/jet/v6"
)

// cache is the cache used by default in a new Set.
type ECache struct {
	m sync.Map
}

func (c *ECache) Get(templatePath string) *jet.Template {
	_t, ok := c.m.Load(templatePath)
	if !ok {
		return nil
	}
	return _t.(*jet.Template)
}

func (c *ECache) Put(templatePath string, t *jet.Template) {
	c.m.Store(templatePath, t)
}

func (c *ECache) Del(templatePath string) {
	c.m.Delete(templatePath)
}
