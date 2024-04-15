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
	mods  = map[string]Modular{}
	mutex = sync.Mutex{}
)

func Register(modular Modular) {
	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := mods[modular.Name()]; ok {
		panic("modular already registered: " + modular.Name())
	}
	mods[modular.Name()] = modular
}

func LoadModules() []Modular {
	var modules []Modular
	mutex.Lock()
	defer mutex.Unlock()
	for _, mod := range mods {
		modules = append(modules, mod)
	}
	return modules
}
