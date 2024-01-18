package types

// ExtraArgs extra arguments to add to each command
type ExtraArgs struct {
	Install    []string `yaml:"install,omitempty" toml:"install,omitempty"`
	Upgrade    []string `yaml:"upgrade,omitempty" toml:"upgrade,omitempty"`
	Template   []string `yaml:"template,omitempty" toml:"template,omitempty"`
	Diff       []string `yaml:"diff,omitempty" toml:"diff,omitempty"`
	ShowValues []string `yaml:"show_values,omitempty" toml:"show_values,omitempty"`
}

type TargetOptions struct {
	AtomicUpdate bool      `yaml:"atomic" toml:"atomic"`
	PassContext  bool      `yaml:"pass_context" toml:"pass_context"` // Pass the context with --kube-context, if false just check the current context
	AskForDryRun bool      `yaml:"ask_dry_run" toml:"ask_dry_run"`   // Ask to Run with --dry-run before
	ExtraArgs    ExtraArgs `yaml:"extra_args,omitempty" toml:"extra_args,omitempty"`
}

type Target struct {
	TargetOptions
	Release    string           `yaml:"release,omitempty" toml:"release,omitempty"`
	Chart      string           `yaml:"chart" toml:"chart"`
	ValueFiles []string         `yaml:"files" toml:"files"` // Common value files regardless env
	Stages     map[string]Stage `yaml:"stages" toml:"stages"`
}

type Stage struct {
	StageOptions
	ValueFiles []string `yaml:"files" toml:"files"` // List of Values to include
}

type StageOptions struct {
	KubeContext  string `yaml:"kube_context,omitempty" toml:"kube_context,omitempty"` // Expected kubecontext, will pass --kube-context if PassContext
	AskForDryRun bool   `yaml:"ask_dry_run" toml:"ask_dry_run"`
}
