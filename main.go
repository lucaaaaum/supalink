package main

import (
	"fmt"
	"supalink/src/filesystem"
	"supalink/src/settings"
	"supalink/src/utils"

	"github.com/spf13/cobra"
)

// objective: supalink src/path/.*S([0-9]{2})E([0-9]{2}).*.mkv -r destination/path/Season\ $STEP/Name (Year) S$1E$2.mkv

var rootCmd = &cobra.Command{
	Use:  "supalink <source path regex> <destination path template>",
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		srcPath := args[0]
		destPath := args[1]

		settings, err := settings.GetSettings(cmd.Flags())
		if err != nil {
			return err
		}

		utils.AddStopSuffixToPattern(&srcPath)

		matchingPathsAndDestinations := filesystem.GetMatchingPathsAndDestinations(srcPath, destPath, settings)

		if len(matchingPathsAndDestinations) == 0 {
			fmt.Println("No matching paths found.")
			return nil
		}

		filesystem.PrintSymlinks(matchingPathsAndDestinations, settings)

		filesystem.CreateSymlinks(matchingPathsAndDestinations, settings)

		return nil
	},
}

func main() {
	flags := rootCmd.Flags()
	settings.SetupFlags(flags)
	rootCmd.Execute()
}
