package global_go_hook

import (
	"github.com/flant/shell-operator/pkg/hook/binding_context"
	"github.com/flant/shell-operator/pkg/kube/object_patch"

	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/pkg/utils"
	"github.com/flant/addon-operator/sdk"
)

func init() {
	sdk.Register(&Simple{})
}

type Simple struct {
}

func (s *Simple) Metadata() *go_hook.HookMetadata {
	return &go_hook.HookMetadata{
		Name:   "simple",
		Path:   "global-hooks/simple",
		Global: true,
	}
}

func (s *Simple) Config() *go_hook.HookConfig {
	return &go_hook.HookConfig{
		OnStartup: &go_hook.OrderedConfig{Order: 1},
	}
}

func (s *Simple) Run(bindingContexts []binding_context.BindingContext, values, configValues utils.Values, objectPatcher *object_patch.ObjectPatcher,
	logLabels map[string]string) (output *go_hook.HookOutput, err error) {

	return &go_hook.HookOutput{}, nil
}
