package hooks

import (
	"time"

	"github.com/flant/shell-operator/pkg/hook/binding_context"
	"github.com/flant/shell-operator/pkg/kube/object_patch"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/pkg/utils"
	"github.com/flant/addon-operator/sdk"

	"github.com/flant/shell-operator/pkg/kube_events_manager/types"
	"github.com/flant/shell-operator/pkg/metric_storage/operation"
)

type podSpecFilteredObj v1.PodSpec

func (ps *podSpecFilteredObj) FilterSelf(obj *unstructured.Unstructured) (interface{}, error) {
	spec := obj.DeepCopyObject().(*v1.Pod).Spec

	return spec, nil
}

type GoHook struct{}

var _ = sdk.Register(&GoHook{})

func (h *GoHook) Config() *go_hook.HookConfig {
	return &go_hook.HookConfig{
		OnStartup: &go_hook.OrderedConfig{
			Order: 10,
			Handler: func(input *go_hook.HookInput) (*go_hook.HookOutput, error) {
				input.LogEntry.Infof("Hello from module 'hooks-only' golang hook 'go_hooks'!\n")
				return nil, nil
			},
		},

		OnBeforeHelm: &go_hook.OrderedConfig{
			Order: 10,
			Handler: func(input *go_hook.HookInput) (*go_hook.HookOutput, error) {
				input.LogEntry.Infof("Hello from module 'hooks-only' golang hook 'go_hooks' beforeHelm!\n")
				input.LogEntry.Infof("go_hooks beforeHelm hook got values: %s", input.Values.Values.String())
				return nil, nil
			},
		},

		Kubernetes: []go_hook.KubernetesConfig{
			{
				Name:       "pods-for-hooks-only",
				ApiVersion: "v1",
				Kind:       "Pods",
				Group:      "pods",
				//JqFilter:             ".spec",
				FilterFunc:                   sdk.WrapFilterable(&podSpecFilteredObj{}),
				ExecuteHookOnEvents:          []types.WatchEventType{types.WatchEventAdded, types.WatchEventModified, types.WatchEventDeleted},
				ExecuteHookOnSynchronization: true,
				Handler: func(input *go_hook.HookInput) (*go_hook.HookOutput, error) {
					for _, o := range input.BindingContext.Snapshots["pods"] {
						var podSpec podSpecFilteredObj
						err := sdk.UnmarshalFilteredObject(o.FilterResult, &podSpec)
						if err != nil {
							return nil, err
						}
					}

					input.LogEntry.Infof("Hello from on_kube.pods2! I have %d snapshots for '%s' event\n",
						len(input.BindingContext.Snapshots),
						input.BindingContext.WatchEvent)

					input.LogEntry.Infof("go_hooks kube hook got values: %s", input.Values.Values.String())

					return nil, nil
				},
			},
		},

		Schedule: []go_hook.ScheduleConfig{
			{
				Name:    "metrics",
				Crontab: "*/5 * * * * *",
				Group:   "pods",
				Handler: h.SendMetrics,
			},
		},
	}
}

func (h *GoHook) Metadata() *go_hook.HookMetadata {
	return &go_hook.HookMetadata{
		Name:       "go_hook.go",
		Path:       "001-module-go-hooks/hooks/go_hook.go",
		Module:     true,
		ModuleName: "module-go-hooks",
	}
}

func (h *GoHook) Run(bindingContexts []binding_context.BindingContext, values, configValues utils.Values,
	objectPatcher *object_patch.ObjectPatcher, logLabels map[string]string) (*go_hook.HookOutput, error) {
	return nil, nil
}

func (h *GoHook) SendMetrics(input *go_hook.HookInput) (*go_hook.HookOutput, error) {
	input.LogEntry.Infof("Hello from on_kube.pods2! I have %d snapshots for '%s' event\n",
		len(input.BindingContext.Snapshots),
		input.BindingContext.WatchEvent)
	input.LogEntry.Infof("go_hooks schedule hook got values: %s", input.Values.Values.String())

	out := &go_hook.HookOutput{
		Metrics: []operation.MetricOperation{},
	}

	v := 1.0
	out.Metrics = append(out.Metrics, operation.MetricOperation{
		Name: "addon_go_hooks_total",
		Add:  &v,
	})

	input.ConfigValues.Set("moduleGoHooks.time", time.Now().Unix())
	input.Values.Set("moduleGoHooks.time_temp", time.Now().Unix())

	out.ConfigValuesPatches.Operations = input.ConfigValues.GetPatches()
	out.ConfigValuesPatches.Operations = input.ConfigValues.GetPatches()

	return out, nil
}
