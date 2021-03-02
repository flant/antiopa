package global_go_hook

import (
	"github.com/flant/addon-operator/pkg/module_manager"
	"github.com/flant/addon-operator/pkg/utils"
	"github.com/flant/addon-operator/sdk"
)

func init() {
	sdk.Register(&Simple{})
}

type Simple struct {
}

func (s *Simple) Metadata() module_manager.HookMetadata {
	return module_manager.HookMetadata{
		Name:   "simple",
		Path:   "global-hooks/simple",
		Global: true,
	}
}

func (s *Simple) Config() *module_manager.HookConfig {
	return &module_manager.HookConfig{
		YamlConfig: `
configVersion: v1
onStartup: 10
`,
	}
}

func (s *Simple) Run(bindingContexts []..BindingContext, values, configValues utils.Values, logLabels, envs map[string]string) (output *module_manager.HookOutput, err error) {
	return &module_manager.HookOutput{}, nil
}
