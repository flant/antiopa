package hooks

import (
	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
)

var ModuleOneHook = sdk.NewCommonGoHook(&go_hook.HookConfig{})
var _ = sdk.Register(ModuleOneHook)
