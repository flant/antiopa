package module_manager

import (
	"testing"

	"github.com/flant/shell-operator/pkg/kube/object_patch"
	. "github.com/onsi/gomega"

	. "github.com/flant/shell-operator/pkg/hook/binding_context"

	metric_operation "github.com/flant/shell-operator/pkg/metric_storage/operation"

	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/pkg/utils"
	"github.com/flant/addon-operator/sdk/registry"
)

type SimpleHook struct {
}

func (s *SimpleHook) Metadata() *go_hook.HookMetadata {
	return &go_hook.HookMetadata{
		Name:   "simple",
		Path:   "simple",
		Global: true,
	}
}

func (s *SimpleHook) Config() (config *go_hook.HookConfig) {
	return &go_hook.HookConfig{
		OnStartup: &go_hook.OrderedConfig{Order: 10},
	}
}

func (s *SimpleHook) Run(bindingContexts []BindingContext, values, configValues utils.Values, objectPatcher *object_patch.ObjectPatcher,
	logLabels map[string]string) (output *go_hook.HookOutput, err error) {

	return &go_hook.HookOutput{
		MemoryValuesPatches: new(utils.ValuesPatch),
		Metrics: []metric_operation.MetricOperation{
			metric_operation.MetricOperation{},
		},
		Error: nil,
	}, nil
}

func Test_Config_GoHook(t *testing.T) {
	g := NewWithT(t)

	goHook := &SimpleHook{}

	goHookRegistry := registry.Registry()
	goHookRegistry.Add(goHook)

	moduleManager := NewMainModuleManager()

	gh := NewGlobalHook("simple", "simple")
	gh.WithGoHook(goHook)
	err := gh.WithGoConfig(goHook.Config())
	g.Expect(err).ShouldNot(HaveOccurred())
	gh.WithModuleManager(moduleManager)

	bc := []BindingContext{}

	e := NewHookExecutor(gh, bc, "v1", nil)
	res, err := e.Run()
	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(res.Patches).ShouldNot(BeEmpty())
	g.Expect(res.Metrics).ShouldNot(BeEmpty())
}
