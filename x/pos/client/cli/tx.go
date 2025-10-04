package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/NeomSense/PoS/x/pos/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdSubmitRecord(),
		CmdVerifyRecord(),
	)

	return cmd
}

// CmdSubmitRecord implements the submit-record command
func CmdSubmitRecord() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-record [data-file] [merkle-root]",
		Short: "Submit a new record as a validator",
		Long: `Submit a new proof-of-record as a validator.
The data-file should contain the record data, and merkle-root is the merkle root hash of the data.

Example:
  posd tx pos submit-record ./my-record.json abc123def456... --from validator1`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Read data from file
			dataFile := args[0]
			data, err := os.ReadFile(dataFile)
			if err != nil {
				return fmt.Errorf("failed to read data file: %w", err)
			}

			merkleRoot := args[1]

			// Get validator address from the from flag
			validatorAddr := clientCtx.GetFromAddress().String()

			msg := &types.MsgSubmitRecord{
				ValidatorAddress: validatorAddr,
				Data:             data,
				MerkleRoot:       merkleRoot,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdVerifyRecord implements the verify-record command
func CmdVerifyRecord() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify-record [record-id] [approved]",
		Short: "Verify a record (true/false for approval)",
		Long: `Verify a submitted record by record ID.
The approved parameter should be 'true' to approve or 'false' to reject.

Example:
  posd tx pos verify-record abc123 true --from verifier1`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			recordID := args[0]
			approved := args[1] == "true"

			verifierAddr := clientCtx.GetFromAddress().String()

			msg := &types.MsgVerifyRecord{
				Verifier: verifierAddr,
				RecordId: recordID,
				Approved: approved,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
