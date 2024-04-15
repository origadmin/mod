package mod

import (
	"context"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
)

// Modular is an interface that defines the methods that a module must implement
type Modular interface {
	AutoMigrate(ctx context.Context) error
	Init(ctx context.Context, config Config) error
	Release(ctx context.Context) error
	RegisterRouters(ctx context.Context, version string, router *gin.RouterGroup) error
}

var (
	mods   = map[string]Modular{}
	modsMu = sync.Mutex{}
)

func Register(name string, modular Modular) {
	modsMu.Lock()
	defer modsMu.Unlock()
	if modular == nil {
		panic("mods: Register module is nil")
	}
	if _, ok := mods[name]; ok {
		panic("mods: Register called twice for module " + name)
	}
	mods[name] = modular
}

// LoadModules loads all modules
func LoadModules() []Modular {
	var modules []Modular
	modsMu.Lock()
	for _, mod := range mods {
		modules = append(modules, mod)
	}
	modsMu.Unlock()
	return modules
}

// Load loads a module by name
func Load(name string) (Modular, error) {
	modsMu.Lock()
	m, ok := mods[name]
	modsMu.Unlock()
	if !ok {
		return nil, fmt.Errorf("sql: unknown module %q (forgotten import?)", name)
	}
	return m, nil
}
