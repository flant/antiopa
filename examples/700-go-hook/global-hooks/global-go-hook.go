package global_hooks

import (
	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
)

func init() {
	h := GoHook{}
	sdk.Register(sdk.NewCommonGoHook(&go_hook.HookConfig{
		OnStartup:   &go_hook.OrderedConfig{Order: 10},
		MainHandler: h.Main,
	}))
}

type GoHook struct {
	sdk.CommonGoHook
}

func (h *GoHook) Main(input *go_hook.HookInput) (*go_hook.HookOutput, error) {
	input.LogEntry.Infof("Start Global Go hook")
	return nil, nil
}
