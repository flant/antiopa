package global_hooks

import (
	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
)

func init() {
	sdk.Register(sdk.NewCommonGoHook(&go_hook.HookConfig{}))
}

type GoHook struct {
	sdk.CommonGoHook
}
