package list

import (
	"fmt"
	"github.com/ekristen/azure-nuke/pkg/nuke"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"sort"

	"github.com/ekristen/azure-nuke/pkg/commands/global"
	"github.com/ekristen/azure-nuke/pkg/common"

	"github.com/ekristen/libnuke/pkg/resource"

	_ "github.com/ekristen/azure-nuke/resources"
)

func execute(c *cli.Context) error {
	ls := resource.GetNames()

	sort.Strings(ls)

	for _, name := range ls {
		reg := resource.GetRegistration(name)

		if reg.AlternativeResource != "" {
			color.New(color.Bold).Printf("%-55s\n", name)
			color.New(color.Bold, color.FgYellow).Printf("  > %-55s", reg.AlternativeResource)
			color.New(color.FgCyan).Printf("alternative resource\n")
		} else {
			color.New(color.Bold).Printf("%-55s", name)
			c := color.FgGreen
			if reg.Scope == nuke.Tenant {
				c = color.FgHiGreen
			} else if reg.Scope == nuke.Subscription {
				c = color.FgHiBlue
			} else if reg.Scope == nuke.ResourceGroup {
				c = color.FgHiMagenta
			} else {

			}
			color.New(c).Printf(fmt.Sprintf("%s\n", string(reg.Scope)))
		}
	}

	return nil
}

func init() {
	cmd := &cli.Command{
		Name:    "resource-types",
		Aliases: []string{"list-resources"},
		Usage:   "list available resources to nuke",
		Flags:   global.Flags(),
		Before:  global.Before,
		Action:  execute,
	}

	common.RegisterCommand(cmd)
}