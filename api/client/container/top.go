package container

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/client"
	"github.com/docker/docker/cli"
	"github.com/spf13/cobra"
)

type topOptions struct {
	container string

	args []string
}

// NewTopCommand creates a new cobra.Command for `docker top`
func NewTopCommand(dockerCli *client.DockerCli) *cobra.Command {
	var opts topOptions

	cmd := &cobra.Command{
		Use:   "top CONTAINER [ps OPTIONS]",
		Short: "显示一个运行容器中运行的所有进程信息",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.container = args[0]
			opts.args = args[1:]
			return runTop(dockerCli, &opts)
		},
	}
	cmd.SetFlagErrorFunc(flagErrorFunc)

	flags := cmd.Flags()
	flags.SetInterspersed(false)

	return cmd
}

func runTop(dockerCli *client.DockerCli, opts *topOptions) error {
	ctx := context.Background()

	procList, err := dockerCli.Client().ContainerTop(ctx, opts.container, opts.args)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(dockerCli.Out(), 20, 1, 3, ' ', 0)
	fmt.Fprintln(w, strings.Join(procList.Titles, "\t"))

	for _, proc := range procList.Processes {
		fmt.Fprintln(w, strings.Join(proc, "\t"))
	}
	w.Flush()
	return nil
}
