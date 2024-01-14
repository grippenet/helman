package types

type Config struct {
	Vars        map[string]string       `yaml:"vars" toml:"vars"`
	Globals     TargetOptions           `yaml:"globals" toml:"globals"`
	Targets     map[string]Target       `yaml:"targets" toml:"targets"`
	KnownStages map[string]StageOptions `yaml:"stages"`
	File        string
}
