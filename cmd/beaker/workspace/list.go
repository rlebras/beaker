package workspace

import (
	"os"
	"encoding/json"
	"context"
	"fmt"

	"github.com/beaker/client/api"
	beaker "github.com/beaker/client/client"
	"github.com/fatih/color"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/allenai/beaker/config"
)

type listOptions struct {
	org         string
}

func newListCmd(
	parent *kingpin.CmdClause,
	parentOpts *workspaceOptions,
	config *config.Config,
) {
	o := &listOptions{}
	cmd := parent.Command("list", "List accessible workspaces")
	cmd.Action(func(c *kingpin.ParseContext) error {
		beaker, err := beaker.NewClient(parentOpts.addr, config.UserToken)
		if err != nil {
			return err
		}
		if o.org == "" {
			o.org = config.DefaultOrg
		}
		return o.run(beaker)
	})

	cmd.Flag("org", "Organization to search for accessible workspaces").Short('o').StringVar(&o.org)
}

func (o *listOptions) run(beakerClient *beaker.Client) error {
	fmt.Println(color.RedString("Workspace commands are still under development and should be considered experimental."))

	ctx := context.TODO()

	var workspaces []api.Workspace
	cursor := ""
	for {
		var results []api.Workspace
		var err error
		results, cursor, err = beakerClient.ListWorkspaces(ctx, o.org, &beaker.ListWorkspaceOptions{Cursor: cursor})
		if err != nil {
			return err
		}
		workspaces = append(workspaces, results...)
		if cursor == "" {
			break
		}
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "    ")
	return encoder.Encode(workspaces)
}
