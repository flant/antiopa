package sublevel

import (
	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
)

func init() {
	sdk.Register(&SubSubHook{})
}

type SubSubHook struct {
	sdk.CommonGoHook
}

func (h *SubSubHook) Metadata() go_hook.HookMetadata {
	return h.Metadata()
}

func (h *SubSubHook) Config() *go_hook.HookConfig {
	return h.CommonGoHook.HookConfig
}
