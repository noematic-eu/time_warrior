package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var feeCmd = &cobra.Command{
	Use:   "fee [project] [amount]",
	Short: "Set or show the fee per hour for a project",
	Long: `Set or show the fee per hour for a project.
If no amount is provided, the current fee will be displayed.
If an amount is provided, it will be set as the fee per hour for the project.

Examples:
  tw fee MyProject          # Show current fee for MyProject
  tw fee MyProject 50       # Set fee to $50/hour for MyProject`,
	Aliases: []string{"f"},
	Args:    cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		config := initializeConfig()
		project := args[0]

		if len(args) == 1 {
			// Show current fee
			fee, err := config.GetProjectFee(project)
			if err != nil {
				fmt.Println(err)
				return
			}
			if fee == 0 {
				fmt.Printf("No fee set for project: %s\n", project)
			} else {
				fmt.Printf("Fee per hour for %s: $%.2f\n", project, fee)
			}
			return
		}

		// Set new fee
		var fee float64
		if _, err := fmt.Sscanf(args[1], "%f", &fee); err != nil {
			fmt.Printf("Invalid fee format: %v\n", err)
			return
		}
		if err := config.SetProjectFee(project, fee); err != nil {
			fmt.Printf("Error setting fee: %v\n", err)
			return
		}
		fmt.Printf("Fee per hour set to $%.2f for project: %s\n", fee, project)
	},
}

func init() {
	rootCmd.AddCommand(feeCmd)
}
