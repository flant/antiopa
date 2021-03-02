package hooks

import (
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/flant/addon-operator/pkg/module_manager"
	"github.com/flant/addon-operator/pkg/utils"
	"github.com/flant/addon-operator/pkg/utils/values_store"
	"github.com/flant/addon-operator/sdk"

	"github.com/flant/shell-operator/pkg/app"
	"github.com/flant/shell-operator/pkg/kube"
	"github.com/flant/shell-operator/pkg/kube_events_manager/types"
	"github.com/flant/shell-operator/pkg/metric_storage/operation"
)

var _ = sdk.Register(&GoHook{})

type GoHook struct {
	sdk.CommonGoHook
	kubeClient kube.KubernetesClient
}

type podSpecFilteredObj v1.PodSpec

func (ps *podSpecFilteredObj) FilterSelf(obj *unstructured.Unstructured) (interface{}, error) {
	runtimeObject, err := sdk.Convert(obj)
	if err != nil {
		return "", err
	}

	return runtimeObject.(*v1.Pod).Spec, nil
}

func (h *GoHook) Metadata() module_manager.HookMetadata {
	return module_manager.HookMetadata{
		Name:       "go_hook.go",
		Path:       "001-module-go-hooks/hooks/go_hook.go",
		Module:     true,
		ModuleName: "module-go-hooks",
	}
}

func (h *GoHook) Config() *module_manager.HookConfig {
	return h.CommonGoHook.Config(&module_manager.HookConfig{
		OnStartup: &module_manager.OrderedConfig{
			Order: 10,
			Handler: func(input *module_manager.BindingInput) (*module_manager.BindingOutput, error) {
				input.LogEntry.Infof("Hello from module 'hooks-only' golang hook 'go_hooks'!\n")
				return nil, nil
			},
		},

		OnBeforeHelm: &module_manager.OrderedConfig{
			Order: 10,
			Handler: func(input *module_manager.BindingInput) (*module_manager.BindingOutput, error) {
				input.LogEntry.Infof("Hello from module 'hooks-only' golang hook 'go_hooks' beforeHelm!\n")
				vs := values_store.NewValuesStoreFromValues(input.Values)

				input.LogEntry.Infof("go_hooks beforeHelm hook got values: %s", vs.GetAsYaml())

				return nil, nil
			},
		},

		Kubernetes: []module_manager.KubernetesConfig{
			{
				Name:       "pods-for-hooks-only",
				ApiVersion: "v1",
				Kind:       "Pods",
				Group:      "pods",
				//JqFilter:             ".spec",
				FilterFunc:                   sdk.WrapFilterable(&podSpecFilteredObj{}),
				ExecuteHookOnEvents:          []types.WatchEventType{types.WatchEventAdded, types.WatchEventModified, types.WatchEventDeleted},
				ExecuteHookOnSynchronization: true,
				Handler: func(input *module_manager.BindingInput) (*module_manager.BindingOutput, error) {
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

					vs := values_store.NewValuesStoreFromValues(input.Values)

					input.LogEntry.Infof("go_hooks kube hook got values: %s", vs.GetAsYaml())

					return nil, nil
				},
			},
		},

		Schedule: []module_manager.ScheduleConfig{
			{
				Name:    "metrics",
				Crontab: "*/5 * * * * *",
				Handler: h.SendMetrics,
			},
		},
	})
}

// Go hook has no kubectl, but it can initialize its own kubernetes client!
func (h *GoHook) initKubeClient() error {
	if h.kubeClient != nil {
		return nil
	}

	h.kubeClient = kube.NewKubernetesClient()
	h.kubeClient.WithContextName(app.KubeContext)
	h.kubeClient.WithConfigPath(app.KubeConfig)
	h.kubeClient.WithRateLimiterSettings(app.KubeClientQps, app.KubeClientBurst)
	// Initialize kube client for kube events hooks.
	err := h.kubeClient.Init()
	if err != nil {
		return err
	}
	return nil
}

func (h *GoHook) SendMetrics(input *module_manager.BindingInput) (*module_manager.BindingOutput, error) {
	err := h.initKubeClient()
	if err != nil {
		input.LogEntry.Errorf("Fatal: initialize kube client: %s", err)
		return nil, err
	}

	podList, err := h.kubeClient.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		input.LogEntry.Errorf("Fatal: cannot list pods: %s", err)
		return nil, err
	}

	input.LogEntry.Infof("Get %d pods:", len(podList.Items))
	for i, pod := range podList.Items {
		input.LogEntry.Infof("%02d. Pod/%s in ns/%s", i, pod.Name, pod.Namespace)
	}

	input.LogEntry.Infof("Hello from on_kube.pods2! I have %d snapshots for '%s' event\n",
		len(input.BindingContext.Snapshots),
		input.BindingContext.WatchEvent)

	vs := values_store.NewValuesStoreFromValues(input.Values)

	input.LogEntry.Infof("go_hooks schedule hook got values: %s", vs.GetAsYaml())

	out := &module_manager.BindingOutput{
		Metrics: []operation.MetricOperation{},
	}

	v := 1.0
	out.Metrics = append(out.Metrics, operation.MetricOperation{
		Name: "addon_go_hooks_total",
		Add:  &v,
	})

	out.ConfigValuesPatches = &utils.ValuesPatch{
		[]*utils.ValuesPatchOperation{
			{
				Op:    "add",
				Path:  "/moduleGoHooks/time",
				Value: fmt.Sprintf("%d", time.Now().Unix()),
			},
		},
	}

	out.MemoryValuesPatches = &utils.ValuesPatch{
		[]*utils.ValuesPatchOperation{
			{
				Op:    "add",
				Path:  "/moduleGoHooks/time_temp",
				Value: fmt.Sprintf("%d", time.Now().Unix()),
			},
		},
	}
	return out, nil
}
