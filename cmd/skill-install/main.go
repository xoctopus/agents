package main

import (
	"github.com/spf13/cobra"
	_ "github.com/xoctopus/concx/pkg/chanx"
	"github.com/xoctopus/confx/pkg/cmdx"
	_ "github.com/xoctopus/confx/pkg/types"
	_ "github.com/xoctopus/logx"
	_ "github.com/xoctopus/sqlx/pkg/builder"
	_ "github.com/xoctopus/x/testx"

	"github.com/xoctopus/agents/internal/installer"
)

var root *cobra.Command

func init() {
	root = &cobra.Command{}

	root.AddCommand(cmdx.NewCommand("install", &installer.Installer{}).Cmd())
	root.AddCommand(cmdx.NewCommand("version", &installer.Version{}).Cmd())
}

func main() {
	if err := root.Execute(); err != nil {
		root.Println(err)
		return
	}
}
