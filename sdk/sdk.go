package sdk

import (
	"fmt"
	"regexp"
	"runtime"

	log "github.com/sirupsen/logrus"

	hooktypes "github.com/flant/addon-operator/pkg/hook/types"
	"github.com/flant/addon-operator/pkg/module_manager/go_hook"

	"github.com/flant/shell-operator/pkg/hook/binding_context"
	shhooktypes "github.com/flant/shell-operator/pkg/hook/types"
	metricoperation "github.com/flant/shell-operator/pkg/metric_storage/operation"

	"github.com/flant/addon-operator/pkg/utils"
)

// Register is a method to define go hooks.
// return value is for trick with
//   var _ =
var Register = func(_ go_hook.GoHook) bool { return false }

var _ go_hook.GoHook = (*CommonGoHook)(nil)

// /path/.../global-hooks/a/b/c/hook-name.go
// $1 - hook path
// $3 - hook name
var globalRe = regexp.MustCompile(`/global-hooks/(([^/]+/)*([^/]+))$`)

// /path/.../modules/module-name/hooks/a/b/c/hook-name.go
// $1 - hook path
// $2 - module name
// $4 - hook name
var moduleRe = regexp.MustCompile(`/modules/(([^/]+)/hooks/([^/]+/)*([^/]+))$`)

var moduleNameRe = regexp.MustCompile(`^[0-9][0-9][0-9]-(.*)$`)

type CommonGoHook struct {
	config *go_hook.HookConfig
}

func NewCommonGoHook(config *go_hook.HookConfig) *CommonGoHook {
	return &CommonGoHook{config: config}
}

func (h CommonGoHook) Config() *go_hook.HookConfig {
	return h.config
}

func (*CommonGoHook) Metadata() *go_hook.HookMetadata {
	hookMeta := &go_hook.HookMetadata{
		Name:       "",
		Path:       "",
		Global:     false,
		Module:     false,
		ModuleName: "",
	}

	_, f, _, _ := runtime.Caller(1)

	matches := globalRe.FindStringSubmatch(f)
	if matches != nil {
		hookMeta.Global = true
		hookMeta.Name = matches[3]
		hookMeta.Path = matches[1]
	} else {
		matches = moduleRe.FindStringSubmatch(f)
		if matches != nil {
			hookMeta.Module = true
			hookMeta.Name = matches[4]
			hookMeta.Path = matches[1]
			modNameMatches := moduleNameRe.FindStringSubmatch(matches[2])
			if modNameMatches != nil {
				hookMeta.ModuleName = modNameMatches[1]
			}
		}
	}

	return hookMeta
}

// Run executes a handler like in BindingContext.MapV1 or in framework/shell.
func (h *CommonGoHook) Run(bindingContexts []binding_context.BindingContext, values, configValues utils.Values, logLabels map[string]string) (*go_hook.HookOutput, error) {
	logEntry := log.WithFields(utils.LabelsToLogFields(logLabels)).
		WithField("output", "golang")

	out := &go_hook.HookOutput{
		ConfigValuesPatches: utils.NewValuesPatch(),
		MemoryValuesPatches: utils.NewValuesPatch(),
		Metrics:             make([]metricoperation.MetricOperation, 0),
	}

	for _, bc := range bindingContexts {
		patchableValues, err := go_hook.NewPatchableValues(values)
		if err != nil {
			return nil, err
		}

		patchableConfigValues, err := go_hook.NewPatchableValues(configValues)
		if err != nil {
			return nil, err
		}

		bindingInput := &go_hook.HookInput{
			BindingContext: bc,
			Values:         patchableValues,
			ConfigValues:   patchableConfigValues,
			LogEntry:       logEntry,
			LogLabels:      logLabels,
		}

		handler := func() go_hook.BindingHandler {
			if bc.Metadata.Group != "" {
				if h.config.GroupHandlers != nil {
					h := h.config.GroupHandlers[bc.Metadata.Group]
					if h != nil {
						return h
					}
				}
				return h.config.MainHandler
			}

			handler := h.config.MainHandler
			switch bc.Metadata.BindingType {
			case shhooktypes.OnStartup:
				if h.config.OnStartup != nil {
					handler = h.config.OnStartup.Handler
				}
			case hooktypes.BeforeAll:
				if h.config.OnBeforeAll != nil {
					handler = h.config.OnBeforeAll.Handler
				}
			case hooktypes.AfterAll:
				if h.config.OnAfterAll != nil {
					handler = h.config.OnAfterAll.Handler
				}
			case hooktypes.BeforeHelm:
				if h.config.OnBeforeHelm != nil {
					handler = h.config.OnBeforeHelm.Handler
				}
			case hooktypes.AfterHelm:
				if h.config.OnAfterHelm != nil {
					handler = h.config.OnAfterHelm.Handler
				}
			case hooktypes.AfterDeleteHelm:
				if h.config.OnAfterDeleteHelm != nil {
					handler = h.config.OnAfterDeleteHelm.Handler
				}
			case shhooktypes.Schedule:
				for _, sc := range h.config.Schedule {
					if sc.Name == bc.Binding {
						handler = sc.Handler
						break
					}
				}
			case shhooktypes.OnKubernetesEvent:
				// Find handler by name.
				// TODO split to Synchronization and Event handlers?
				for _, sc := range h.config.Kubernetes {
					if sc.Name == bc.Binding {
						handler = sc.Handler
						break
					}
				}
			}

			return handler
		}()

		if handler == nil {
			return nil, fmt.Errorf("no handler defined for binding context type=%s binding=%s group=%s", bc.Metadata.BindingType, bc.Binding, bc.Metadata.Group)
		}

		bindingOut, err := handler(bindingInput)
		if err != nil {
			return nil, err
		}
		if bindingOut != nil {
			if bindingOut.ConfigValuesPatches != nil {
				out.ConfigValuesPatches.MergeOperations(bindingOut.ConfigValuesPatches)
			}
			if bindingOut.MemoryValuesPatches != nil {
				out.MemoryValuesPatches.MergeOperations(bindingOut.MemoryValuesPatches)
			}
			if bindingOut.Metrics != nil {
				out.Metrics = append(out.Metrics, bindingOut.Metrics...)
			}
		}
	}

	return out, nil
}
