package command

import (
	"fmt"

	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/moduledeps"
	"github.com/xlab/treeprint"
)

// ProvidersCommand is a Command implementation that prints out information
// about the providers used in the current configuration/state.
type ProvidersCommand struct {
	Meta
}

func (c *ProvidersCommand) Help() string {
	return providersCommandHelp
}

func (c *ProvidersCommand) Synopsis() string {
	return "Prints a tree of the providers used in the configuration"
}

func (c *ProvidersCommand) Run(args []string) int {
	args = c.Meta.process(args)
	cmdFlags := c.Meta.defaultFlagSet("providers")
	cmdFlags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := cmdFlags.Parse(args); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing command-line flags: %s\n", err.Error()))
		return 1
	}

	/*
		configPath, err := ModulePath(cmdFlags.Args())
		if err != nil {
			c.Ui.Error(err.Error())
			return 1
		}

		var diags tfdiags.Diagnostics

		empty, err := configs.IsEmptyDir(configPath)
		if err != nil {
			diags = diags.Append(tfdiags.Sourceless(
				tfdiags.Error,
				"Error validating configuration directory",
				fmt.Sprintf("Terraform encountered an unexpected error while verifying that the given configuration directory is valid: %s.", err),
			))
			c.showDiagnostics(diags)
			return 1
		}
		if empty {
			absPath, err := filepath.Abs(configPath)
			if err != nil {
				absPath = configPath
			}
			diags = diags.Append(tfdiags.Sourceless(
				tfdiags.Error,
				"No configuration files",
				fmt.Sprintf("The directory %s contains no Terraform configuration files.", absPath),
			))
			c.showDiagnostics(diags)
			return 1
		}

		config, configDiags := c.loadConfig(configPath)
		diags = diags.Append(configDiags)
		if configDiags.HasErrors() {
			c.showDiagnostics(diags)
			return 1
		}

		// Load the backend
		b, backendDiags := c.Backend(&BackendOpts{
			Config: config.Module.Backend,
		})
		diags = diags.Append(backendDiags)
		if backendDiags.HasErrors() {
			c.showDiagnostics(diags)
			return 1
		}

		// Get the state
		env := c.Workspace()
		state, err := b.StateMgr(env)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Failed to load state: %s", err))
			return 1
		}
		if err := state.RefreshState(); err != nil {
			c.Ui.Error(fmt.Sprintf("Failed to load state: %s", err))
			return 1
		}

		s := state.State()
		depTree := terraform.ConfigTreeDependencies(config, s)
		depTree.SortDescendents()

		printRoot := treeprint.New()
		providersCommandPopulateTreeNode(printRoot, depTree)

		c.Ui.Output(printRoot.String())

		c.showDiagnostics(diags)
		if diags.HasErrors() {
			return 1
		}
	*/

	c.Ui.Output(fmt.Sprintf("terraform providers is temporarily disabled"))
	return 0
}

func providersCommandPopulateTreeNode(node treeprint.Tree, deps *moduledeps.Module) {
	fqns := make([]addrs.Provider, 0, len(deps.Providers))
	for fqn := range deps.Providers {
		fqns = append(fqns, fqn)
	}

	for _, fqn := range fqns {
		dep := deps.Providers[fqn]
		versionsStr := dep.Constraints.String()
		if versionsStr != "" {
			versionsStr = " " + versionsStr
		}
		var reasonStr string
		switch dep.Reason {
		case moduledeps.ProviderDependencyInherited:
			reasonStr = " (inherited)"
		case moduledeps.ProviderDependencyFromState:
			reasonStr = " (from state)"
		}
		node.AddNode(fmt.Sprintf("provider.%s%s%s", fqn.LegacyString(), versionsStr, reasonStr))
	}

	for _, child := range deps.Children {
		childNode := node.AddBranch(fmt.Sprintf("module.%s", child.Name))
		providersCommandPopulateTreeNode(childNode, child)
	}
}

const providersCommandHelp = `
Usage: terraform providers [dir]

  Prints out a tree of modules in the referenced configuration annotated with
  their provider requirements.

  This provides an overview of all of the provider requirements across all
  referenced modules, as an aid to understanding why particular provider
  plugins are needed and why particular versions are selected.

`
