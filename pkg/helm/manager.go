package helm

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/grippenet/helman/pkg/types"
)

func resolveArgs(commandName string, extra types.ExtraArgs) []string {
	if commandName == CommandInstall {
		return extra.Install
	}
	if commandName == CommandUpgrade {
		return extra.Upgrade
	}
	if commandName == CommandTemplate {
		return extra.Template
	}
	if commandName == CommandShowValues {
		return extra.ShowValues
	}
	if commandName == CommandDiff {
		return extra.Diff
	}
	return nil
}

func createExtraArgs(commandName string, extra types.ExtraArgs, from string) []ExtraArg {
	e := make([]ExtraArg, 0)
	args := resolveArgs(commandName, extra)
	if len(args) == 0 {
		return e
	}
	for _, arg := range args {
		e = append(e, ExtraArg{Arg: arg, From: from})
	}
	return e
}

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

// Create list of Helm commands (command & eventual sub command) to use for a given helman command
func (h *HelmManager) resolveHelmCommand(commandName string) []string {
	if commandName == CommandShowValues {
		return []string{"show", "values"}
	}
	if commandName == CommandDiff {
		return []string{"diff", "upgrade"}
	}
	return []string{commandName}
}

// ResolveCommand resolves configuration for a given helm command
func (h *HelmManager) ResolveCommand(commandName string, name string, stageName string, extraArgs []string) (*Resolved, error) {
	target, ok := h.config.Targets[name]
	if !ok {
		return nil, fmt.Errorf("unknown target named '%s' in helman targets config", name)
	}

	var release string
	if target.Release == "" {
		release = name
	} else {
		release = target.Release
	}
	if commandName == CommandShowValues {
		// Dont pass release for show values
		release = ""
	}

	resolved := &Resolved{
		Command:      h.resolveHelmCommand(commandName),
		Release:      release,
		Chart:        target.Chart,
		PassContext:  target.PassContext || h.config.Globals.PassContext,
		AtomicUpdate: target.AtomicUpdate || h.config.Globals.AtomicUpdate,
		KubeContext:  "",
	}

	stage, err := h.resolveStage(target, stageName)
	if err != nil {
		return nil, err
	}

	files := make([]ValueFile, 0, len(target.ValueFiles)+len(stage.ValueFiles))

	vars := copyVars(h.config.Vars)
	vars["stage"] = stageName

	if len(target.ValueFiles) > 0 {
		for i, template := range target.ValueFiles {
			o, err := bindVars(template, vars)
			if err != nil {
				return nil, errors.Join(fmt.Errorf("unable to parse target file template %d ", i), err)
			}
			vf := ValueFile{Resolved: o, Template: template, From: "target"}
			files = append(files, vf)
		}
	}
	if len(stage.ValueFiles) > 0 {
		for i, template := range stage.ValueFiles {
			o, err := bindVars(template, vars)
			if err != nil {
				return nil, errors.Join(fmt.Errorf("unable to parse stage file template %d ", i), err)
			}
			vf := ValueFile{Resolved: o, Template: template, From: "stage"}
			files = append(files, vf)
		}
	}
	resolved.Files = files

	extra := make([]ExtraArg, 0)
	var e []ExtraArg
	e = createExtraArgs(commandName, h.config.Globals.ExtraArgs, "globals")
	if len(e) > 0 {
		extra = append(extra, e...)
	}
	e = createExtraArgs(commandName, target.ExtraArgs, "target")
	if len(e) > 0 {
		extra = append(extra, e...)
	}
	if len(extraArgs) > 0 {
		for _, arg := range extraArgs {
			extra = append(extra, ExtraArg{Arg: arg, From: "command-line"})
		}
	}
	resolved.ExtraArgs = extra
	return resolved, nil
}

// CreateHelmCommand transforms resolved config for helm command to real command to pass to Helm
func (h *HelmManager) CreateHelmCommand(resolved *Resolved) (*Command, error) {

	args := make([]string, 0)

	args = append(args, resolved.Command...)
	if resolved.Release != "" {
		args = append(args, resolved.Release)
	}
	if resolved.Chart != "" {
		args = append(args, resolved.Chart)
	}

	for _, vf := range resolved.Files {
		args = append(args, "-f")
		args = append(args, vf.Resolved)
	}

	cmd := NewCommand(args)

	if resolved.AtomicUpdate {
		cmd.AddArg("--atomic")
	}
	if resolved.KubeContext != "" {
		if resolved.PassContext {
			cmd.AddArg("--kube-context", resolved.KubeContext)
		} else {
			cmd.CheckKubeContext = resolved.KubeContext
		}
	}

	for _, arg := range resolved.ExtraArgs {
		cmd.AddArg(arg.Arg)
	}
	return cmd, nil
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
