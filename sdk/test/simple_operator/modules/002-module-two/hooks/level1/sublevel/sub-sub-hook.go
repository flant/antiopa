package sublevel

import (
	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
)

var ModuleTwoHook = sdk.NewCommonGoHook(&go_hook.HookConfig{})
var _ = sdk.Register(ModuleTwoHook)
