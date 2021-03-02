package hooks

import (
	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
)

func init() {
	sdk.Register(&ModuleOneHook{})
}

type ModuleOneHook struct {
	sdk.CommonGoHook
}

func (h *ModuleOneHook) Metadata() go_hook.HookMetadata {
	return h.CommonGoHook.Metadata()
}
