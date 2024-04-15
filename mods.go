package mod

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
)

// Modular is an interface that defines the methods that a module must implement
type Modular interface {
	Name() string
	AutoMigrate(ctx context.Context) error
	Init(ctx context.Context) error
	Release(ctx context.Context) error
	RegisterRouters(ctx context.Context, version string, router *gin.RouterGroup) error
}

var (
	mods = map[string]Modular{}
	mux  = sync.Mutex{}
)

func Register(modular Modular) {
	mux.Lock()
	defer mux.Unlock()
	if _, ok := mods[modular.Name()]; ok {
		panic("modular already registered: " + modular.Name())
	}
	mods[modular.Name()] = modular
}

func LoadModules() []Modular {
	var modules []Modular
	mux.Lock()
	defer mux.Unlock()
	for _, mod := range mods {
		modules = append(modules, mod)
	}
	return modules
}
