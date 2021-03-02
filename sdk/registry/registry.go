package registry

import (
	"sync"

	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	. "github.com/flant/addon-operator/sdk"
)

var _ = initRegistry()

func initRegistry() bool {
	Register = func(h go_hook.GoHook) bool {
		Registry().Add(h)
		return true
	}
	return true
}

type HookRegistry struct {
	hooks []go_hook.GoHook
	m     sync.Mutex
}

var instance *HookRegistry
var once sync.Once

func Registry() *HookRegistry {
	once.Do(func() {
		instance = new(HookRegistry)
	})
	return instance
}

func (h *HookRegistry) Hooks() []go_hook.GoHook {
	return h.hooks
}

func (h *HookRegistry) Add(hook go_hook.GoHook) {
	h.m.Lock()
	defer h.m.Unlock()
	h.hooks = append(h.hooks, hook)
}
