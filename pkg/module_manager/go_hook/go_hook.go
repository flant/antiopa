package go_hook

import (
	"github.com/Jeffail/gabs"
	"github.com/flant/shell-operator/pkg/hook/binding_context"
	"github.com/flant/shell-operator/pkg/kube_events_manager/types"
	"github.com/flant/shell-operator/pkg/metric_storage/operation"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/flant/addon-operator/pkg/utils"
)

type PatchableValues struct {
	Values      *gabs.Container
	valuesPatch *utils.ValuesPatch
}

func NewPatchableValues(values map[string]interface{}) (*PatchableValues, error) {
	gabsContainer, err := gabs.Consume(values)
	if err != nil {
		return nil, err
	}

	return &PatchableValues{Values: gabsContainer, valuesPatch: utils.NewValuesPatch()}, nil
}

func (p *PatchableValues) Set(path string, value interface{}) {
	op := &utils.ValuesPatchOperation{
		Op:    "add",
		Path:  path,
		Value: value,
	}

	p.valuesPatch.Operations = append(p.valuesPatch.Operations, op)
}

func (p *PatchableValues) Remove(path string) {
	op := &utils.ValuesPatchOperation{
		Op:   "remove",
		Path: path,
	}

	p.valuesPatch.Operations = append(p.valuesPatch.Operations, op)
}

func (p *PatchableValues) GetPatches() *utils.ValuesPatch {
	return p.valuesPatch
}

type HookMetadata struct {
	Name       string
	Path       string
	Global     bool
	Module     bool
	ModuleName string
}

type HookInput struct {
	BindingContext binding_context.BindingContext
	Values         *PatchableValues
	ConfigValues   *PatchableValues
	LogLabels      map[string]string
	LogEntry       *log.Entry
	Envs           map[string]string
}

type HookOutput struct {
	ConfigValuesPatches *utils.ValuesPatch
	MemoryValuesPatches *utils.ValuesPatch
	Metrics             []operation.MetricOperation
	Error               error
}

type BindingHandler func(input *HookInput) (*HookOutput, error)

type HookConfig struct {
	Schedule          []ScheduleConfig
	Kubernetes        []KubernetesConfig
	OnStartup         *OrderedConfig
	OnBeforeHelm      *OrderedConfig
	OnAfterHelm       *OrderedConfig
	OnAfterDeleteHelm *OrderedConfig
	OnBeforeAll       *OrderedConfig
	OnAfterAll        *OrderedConfig
	MainHandler       BindingHandler
	GroupHandlers     map[string]BindingHandler
}

type ScheduleConfig struct {
	Name                 string
	Crontab              string
	AllowFailure         bool
	IncludeSnapshotsFrom []string
	Queue                string
	Group                string
	Handler              BindingHandler
}

type KubernetesConfig struct {
	Name                         string
	ApiVersion                   string
	Kind                         string
	NameSelector                 *types.NameSelector
	NamespaceSelector            *types.NamespaceSelector
	LabelSelector                *v1.LabelSelector
	FieldSelector                *types.FieldSelector
	JqFilter                     string
	IncludeSnapshotsFrom         []string
	Queue                        string
	Group                        string
	ExecuteHookOnEvents          []types.WatchEventType
	ExecuteHookOnSynchronization bool
	WaitForSynchronization       bool
	KeepFullObjectsInMemory      bool
	AllowFailure                 bool
	Handler                      BindingHandler
	FilterFunc                   func(obj *unstructured.Unstructured) (string, error)
}

type OrderedConfig struct {
	Order   float64
	Handler BindingHandler
}

type JqFilterHelper struct {
	Name            string
	JqFilterFn      func(obj *unstructured.Unstructured) (result string, err error)
	ResultConverter func(string, interface{}) error
}

type HookBindingContext struct {
	Type       string // type: Event Synchronization Group Schedule
	Binding    string // binding name
	Snapshots  map[string][]types.ObjectAndFilterResult
	WatchEvent string // Added/Modified/Deleted
	Objects    []types.ObjectAndFilterResult
	Object     types.ObjectAndFilterResult
}

type Handlers struct {
	Main              func()
	Group             map[string]func()
	Kubernetes        map[string]func()
	Schedule          map[string]func()
	OnStartup         func()
	OnBeforeAll       func()
	OnAfterAll        func()
	OnBeforeHelm      func()
	OnAfterHelm       func()
	OnAfterDeleteHelm func()
}

type GoHook interface {
	Metadata() *HookMetadata
	Config() *HookConfig
	Run(bindingContexts []binding_context.BindingContext, values, configValues utils.Values, logLabels map[string]string) (output *HookOutput, err error)
}
