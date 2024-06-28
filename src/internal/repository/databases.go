package repository

import (
	"social-network/src/internal/repository/cache"
	"social-network/src/internal/repository/relation"
	"social-network/src/internal/repository/writeOptimized"
)

type Databases struct {
	cache.Cache
	relation.Relation
	writeOptimized.WriteOptimized
}

func NewDatabases(cache cache.Cache, relation relation.Relation, writeOptimized writeOptimized.WriteOptimized) *Databases {
	return &Databases{
		Cache:          cache,
		Relation:       relation,
		WriteOptimized: writeOptimized,
	}
}
