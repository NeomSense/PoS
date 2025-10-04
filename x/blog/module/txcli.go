package module

import (
    "github.com/spf13/cobra"
    blogcli "<MODULE_PATH>/x/blog/client/cli"
)

func (AppModuleBasic) GetTxCmd() *cobra.Command {
    return blogcli.NewTxCmd()
}
