package go_hooks

import (
	"fmt"

	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
)

const moduleName = "node_manager"

func init() {
	s := JqFilterFunc{}

	sdk.Register(sdk.NewCommonGoHook(&go_hook.HookConfig{
		MainHandler: func(input *go_hook.HookInput) (*go_hook.HookOutput, error) {
			return s.HandleEverything(input)
		},

		Schedule: []go_hook.ScheduleConfig{
			{
				Crontab:              "*/10 * * * * *",
				AllowFailure:         true,
				IncludeSnapshotsFrom: []string{"name1", "name2"},
				Queue:                fmt.Sprintf("/modules/%s/handle_node_templates", moduleName),
				Group:                "",
				Handler: func(input *go_hook.HookInput) (*go_hook.HookOutput, error) {
					return s.HandleNgs(input)
				},
			},
		},

		Kubernetes: []go_hook.KubernetesConfig{
			{
				Name:                         "ngs",
				ApiVersion:                   "",
				Kind:                         "",
				NameSelector:                 nil,
				NamespaceSelector:            nil,
				LabelSelector:                nil,
				FieldSelector:                nil,
				JqFilter:                     "",
				IncludeSnapshotsFrom:         nil,
				Queue:                        "",
				Group:                        "",
				ExecuteHookOnEvents:          nil,
				ExecuteHookOnSynchronization: false,
				WaitForSynchronization:       false,
				KeepFullObjectsInMemory:      false,
				AllowFailure:                 false,
				Handler: func(input *go_hook.HookInput) (*go_hook.HookOutput, error) {
					return nil, nil
				},
				FilterFunc: nil,
			},
		},

		OnStartup: &go_hook.OrderedConfig{
			Order: 20,
			Handler: func(input *go_hook.HookInput) (*go_hook.HookOutput, error) {
				return nil, nil
			},
		},
		OnBeforeHelm: &go_hook.OrderedConfig{
			Order: 20,
			Handler: func(input *go_hook.HookInput) (*go_hook.HookOutput, error) {
				return nil, nil
			},
		},
		OnAfterHelm: &go_hook.OrderedConfig{
			Order: 20,
			Handler: func(input *go_hook.HookInput) (*go_hook.HookOutput, error) {
				return nil, nil
			},
		},
		OnAfterDeleteHelm: &go_hook.OrderedConfig{
			Order: 20,
			Handler: func(input *go_hook.HookInput) (*go_hook.HookOutput, error) {
				return nil, nil
			},
		},

		OnBeforeAll: &go_hook.OrderedConfig{
			Order: 20,
			Handler: func(input *go_hook.HookInput) (*go_hook.HookOutput, error) {
				return nil, nil
			},
		},
		OnAfterAll: &go_hook.OrderedConfig{
			Order: 20,
			Handler: func(input *go_hook.HookInput) (*go_hook.HookOutput, error) {
				return nil, nil
			},
		},

		GroupHandlers: map[string]go_hook.BindingHandler{
			"ngs": s.HandleGroupAzaza,
			"pods": func(input *go_hook.HookInput) (*go_hook.HookOutput, error) {
				return nil, nil
			},
			//	func(input *go_hook.HookInput) (out *go_hook.HookOutput, err error) {
			//	// ...
			//	filterRes, err := ToTaints(input.BindingContext.Snapshots["ngs"][0].FilterResult)
			//
			//	return nil
			//},
		},
	}))
}

type NgsFilterResult struct {
	Annotations map[string]string `json:"annotations"`
	Labels      map[string]string `json:"labels"`
	Taints      []Taint           `json:"taints"`
}

type Taint struct {
	Some   string `json: ...`
	Fields string `json: ...`
	For    string `json: ...`
	Taint  string `json: ...`
}

type JqFilterFunc struct {
	sdk.CommonGoHook
}

func (s *JqFilterFunc) HandleNgs(input *go_hook.HookInput) (output *go_hook.HookOutput, err error) {
	return nil, nil
}

func (s *JqFilterFunc) HandleEverything(input *go_hook.HookInput) (output *go_hook.HookOutput, err error) {
	return nil, nil
}

func (s *JqFilterFunc) HandleGroupAzaza(input *go_hook.HookInput) (output *go_hook.HookOutput, err error) {
	return nil, nil
}
