# Helman helm companion

Helman is an helm command wrapper to manage common configuration to deploy chart in given environment.

Helman expect a configuration to define for each chart, the list of value files to uses. Value files can include common files and stage specific files.
It's also possible to define extra argument to pass to helm.

The goal of helman is to provide a simple way to reproduce command, and share a way to organize value files by enabling naming convention.

Helman aims at cover simple use cases. If you want to cover more complex cases or a complete tool you can have a look to [Helmfile](https://github.com/helmfile/helmfile).

## Runnning helman

Helman provides several commands to run an helm command with arguments built from the target configuration.

The general form is:
helman command [-show|-config] target stage [...extra args]

target : the target name (as define in helman config)
stage  : name of the stage (it must be defined in target or at the global `stages` level)

Available commands:
- install
- upgrade
- diff : will run `diff upgrade`
- template
- show-value will run `show values`

If sub commands '-show' or '-config' is defined, the the command is not run with helm.
-show   : just print the to console the command
-config : show the config after all the files has been resolved for the target with the given stage. It will shown the real file named after variables
have been replaced by their values and the different global/target/stage specific options merged. 

Example

Install a target named `grippenet` with the stage `dev`
```bash
helman install grippenet dev
````

If you want to see the command only
```bash
helman install -show grippenet dev
````

If you want to see how target config is resolved 
```bash
helman install -config grippenet dev
```

## Configuration

Helman configuration is a yaml file (toml is also possible).

Global config structure is

```yaml
# Globale options for each stage (will be applied for all stage with this name, in all targets)
stages: <stages_config>
# Globals definition of target options (can be overriden in each target)
globals: <target_options>
# Variable you can use in files path specifications (${stage} is defined internally by the name of the requested stage)
vars: <vars_config>
# Targets definition
targets: <targets_config>
```

Some options are resolved by merging some sections of configurations

Command options & arguments for a target named 'example' and stage 'prod' will be resolved by using
.globals, targets.example, .stages.prod, .targets.example.stages.prod

Files passed to helm will be resolved with target files + stage files


### stages_config
defines global stage options (to be applied to all stages of all targets)

```yaml
 <stage_name>:
    # Kube context to apply for this stage name
    kube_context: my-context-name
    # Ask for running the command with --dry-run before, only for install|upgrade
    ask_dry_run: true
```

For example, this defines the kube context to apply for 'prod' stage, for all targets
```yaml
stages:
  prod:
    kube_context: prod-kube
```

### vars_config
Configure variable to use in value files path.
It's a simple dictionary

```yaml
vars:
    config: "/path/to/yamls"
    secrets: "/path/to/secrets"
```

Vars are useable in value files path using the `${name}` syntax. The value will be resolved before the command is run.
Helman only defines an automatic variable, named 'stage' (it's not possible to use this variable name), for others names you have to define it

If variable value starts with 'env:', the value will be taken from the environment variable defined after the colon character.
For example "env:MY_ENV" will get the variable value from environment variable 'MY_ENV'

### targets_config
Define the targets. The key of each entry is used as the target name

```yaml
targets:
    good-chart: 
        <target_config>
    other-chart:
        <target_config>
```

### target_config: Define a target

Each target has a common structure:

```yaml
my-target:
    # Specific target options (described in target_options)
    <target_options>

    # Chart name or path to local chart 
    chart: /path/to/chart
    # Release name to use, if not defined, the target name is used
    release: myrelease

    # List of files to include for all stages. It's up to you to define this list (helman doesnt force any organisation)
    # It's possible to use a special variable ${stage}, it will be replaced by the stage name
    files:
        - "/path/to/value/file.yaml"
        # Path using a variable named 'config' (defining the path for base config yaml files)
        # To use a variable you need to define `vars` at the global level (helman only provides ${stage})
        - "${config}/base.yaml"
        # Using 
        - "${config}/${stage}.yaml" 
        - "${secrets}/base.yaml"
        # By using ${stage} variable you dont need to define it in each stage if you follow a naming convention. But it's up to you.
        - "${secrets}/${stage}.yaml"
    stages:
        # Stages allows to defined stage specific options or files
        # Options can be defined for each stage or defined globally in stages.
        # Global stage options are then used as default value.
        prod:
            # Kube context to use. It's also possible to defined it at the global level so all stage with this name will use the same kube context.
            kube_context: my-prod-context
            # If true Ask for running command with --dry-run before to run for good, only for install|upgrade
            ask_dry_run: true
            # List of values files to add for this stage.
            # Stage-specific value files are added after the target's ones.
            files:
                - /path/for/yaml/to/include/in/prod/only.yaml
                # It's also possible to use variables in stage files
                - ${config}/prod.yaml 

````

### target_options (globals or in targets.[target])
```yaml
  # if true kube_context will passed to helm using --kube-context option, if false, the current context will be checked to be this one 
  # before running the command
  pass_context: false
  # Use Atomic, if true install/upgrade command will use --atomic option 
  atomic: true
  # If true Ask for running command with --dry-run before to run for good, only for install|upgrade
  ask_dry_run: true
  extra_args: <extra_args>
```