package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"supalink/src/output"
	"supalink/src/settings"
)

func CreateSymlinks(matchingPathsAndDestinations map[string]string, settings settings.Settings) {
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

func PrintSymlinks(matchingPathsAndDestinations map[string]string, settings settings.Settings) {
	switch settings.Format {
	case output.TreeFormat:
		output.PrintAsTree(matchingPathsAndDestinations)
	case output.TableFormat:
		output.PrintAsTable(matchingPathsAndDestinations)
	}
}
