package types

type TargetOptions struct {
	AtomicUpdate bool `yaml:"atomic" toml:"atomic"`
	PassContext  bool `yaml:"pass_context" toml:"pass_context"` // Pass the context with --kube-context, if false just check the current context
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
	KubeContext string `yaml:"kube_context,omitempty" toml:"kube_context,omitempty"` // Expected kubecontext, will pass --kube-context if PassContext
}
