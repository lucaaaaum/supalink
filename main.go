package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"supalink/src/filesystem"
	"supalink/src/output"
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

		createSymlinks(matchingPathsAndDestinations, settings)

		return nil
	},
}

func createSymlinks(matchingPathsAndDestinations map[string]string, settings settings.Settings) {
	printSymlinks(matchingPathsAndDestinations, settings)

	if settings.DryRun {
		fmt.Println("Dry run enabled, no symlinks will be created.")
		return
	}

	if settings.Confirm {
		var response string
		fmt.Print("Are you sure you want to create these symlinks? (y/n): ")
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" {
			fmt.Println("Operation cancelled by user.")
			return
		}
	}

	for source, destination := range matchingPathsAndDestinations {
		os.MkdirAll(filepath.Dir(destination), os.ModePerm)
		err := os.Symlink(source, destination)
		if err != nil {
			fmt.Printf("Failed to create symlink: %s -> %s. Error: %v\n", source, destination, err)
		} else {
		}
	}
}

func printSymlinks(matchingPathsAndDestinations map[string]string, settings settings.Settings) {
	switch settings.Format {
	case output.TreeFormat:
		output.PrintAsTree(matchingPathsAndDestinations)
	case output.TableFormat:
		output.PrintAsTable(matchingPathsAndDestinations)
	}
}

func main() {
	flags := rootCmd.Flags()
	settings.SetupFlags(flags)
	rootCmd.Execute()
}
