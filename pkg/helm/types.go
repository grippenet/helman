package helm

import "fmt"

const (
	CommandInstall    = "install"
	CommandUpgrade    = "upgrade"
	CommandDiff       = "diff"
	CommandTemplate   = "template"
	CommandShowValues = "show-values"
)

type Resolved struct {
	Command      []string
	Release      string
	Chart        string
	Files        []ValueFile
	ExtraArgs    []ExtraArg
	KubeContext  string
	PassContext  bool
	AtomicUpdate bool
	AskForDryRun bool
	Dir          string
}

func (r *Resolved) Print() {
	fmt.Printf("Dir: %s\n", r.Dir)
	fmt.Printf("Helm Commands: %v\n", r.Command)
	fmt.Printf("Release : %s %s\n", r.Release, r.Chart)
	var pass string
	if r.PassContext {
		pass = "passing --kube-context"
	} else {
		pass = "checking kube context"
	}
	fmt.Printf("KubeContext : %s (%s)\n", r.KubeContext, pass)
	fmt.Printf("Ask for --dry-run : %t\n", r.AskForDryRun)

	fmt.Println("Files:")
	for _, file := range r.Files {
		file.Print()
	}
	fmt.Println("Extra args:")
	for _, arg := range r.ExtraArgs {
		arg.Print()
	}
}

type ValueFile struct {
	Resolved string
	Template string
	From     string
}

func (vf *ValueFile) Print() {
	fmt.Printf("- %s [%s] (%s)\n", vf.Resolved, vf.Template, vf.From)
}

type ExtraArg struct {
	Arg  string
	From string
}

func (arg *ExtraArg) Print() {
	fmt.Printf(" %s (%s)\n", arg.Arg, arg.From)
}
