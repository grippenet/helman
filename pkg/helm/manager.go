package helm

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/grippenet/helman/pkg/types"
)

type HelmManager struct {
	config *types.Config
}

func NewHelmManager(config *types.Config) *HelmManager {
	return &HelmManager{config: config}
}

// resolveStage resolve values, applying chart's options & globals if necessary
func (h *HelmManager) resolveStage(target types.Target, stageName string) (types.Stage, error) {
	defOpts, hasDefault := h.config.KnownStages[stageName]
	o, ok := target.Stages[stageName]

	if hasDefault {
		if o.KubeContext == "" {
			o.KubeContext = defOpts.KubeContext
		}
	}

	if !ok {
		if hasDefault {
			return o, nil
		} else {
			return o, fmt.Errorf("unknown stage '%s'", stageName)
		}
	}
	return o, nil
}

func (h *HelmManager) resolveChartOptions(target types.Target) types.TargetOptions {
	o := types.TargetOptions{}
	o.PassContext = target.PassContext || h.config.Globals.PassContext
	o.AtomicUpdate = target.AtomicUpdate || h.config.Globals.AtomicUpdate
	return o
}

func (h *HelmManager) ChartCommand(commandName string, name string, stageName string, extraArgs []string) (*Command, error) {
	target, ok := h.config.Targets[name]
	if !ok {
		return nil, fmt.Errorf("unknown target named '%s' in helman targets config", name)
	}

	opts := h.resolveChartOptions(target)

	var release string
	if target.Release == "" {
		release = name
	} else {
		release = target.Release
	}

	stage, err := h.resolveStage(target, stageName)
	if err != nil {
		return nil, err
	}

	fileTemplates := make([]string, 0, len(target.ValueFiles)+len(stage.ValueFiles))

	if len(target.ValueFiles) > 0 {
		fileTemplates = append(fileTemplates, target.ValueFiles...)
	}
	if len(stage.ValueFiles) > 0 {
		fileTemplates = append(fileTemplates, stage.ValueFiles...)
	}

	args := make([]string, 0)

	args = append(args, release, target.Chart)
	valueFiles, err := h.ResolveValueFiles(fileTemplates, stageName)
	if err != nil {
		return nil, err
	}
	for _, file := range valueFiles {
		args = append(args, "-f", file)
	}

	cmd := NewCommand(commandName, args)

	if opts.AtomicUpdate {
		cmd.AddArg("--atomic")
	}
	if stage.KubeContext != "" {
		if opts.PassContext {
			cmd.AddArg("--kube-context", stage.KubeContext)
		} else {
			cmd.CheckKubeContext = stage.KubeContext
		}
	}

	cmd.AddArg(extraArgs...)

	return cmd, nil
}

func (h *HelmManager) ResolveValueFiles(fileTemplates []string, stage string) ([]string, error) {
	vars := copyVars(h.config.Vars)
	vars["stage"] = stage
	files := make([]string, 0, len(fileTemplates))
	for _, file := range fileTemplates {
		o, err := bindVars(file, vars)
		if err != nil {
			return files, err
		}
		files = append(files, o)
	}
	return files, nil
}

var VarRegexp = regexp.MustCompile(`\$\{(\w+)\}`)

func bindVars(s string, vars map[string]string) (string, error) {
	vv := VarRegexp.FindAllStringSubmatch(s, -1)
	if len(vv) == 0 {
		return s, nil
	}
	out := s
	for _, v := range vv {
		ref := v[0]
		name := v[1]
		value, ok := vars[name]
		if !ok {
			return s, fmt.Errorf("unknown var %s", name)
		}
		out = strings.Replace(out, ref, value, -1)
	}
	return out, nil
}

// Create new map
func copyVars(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}
