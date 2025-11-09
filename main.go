package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// objective: supalink src/path/.*S([0-9]{2})E([0-9]{2}).*.mkv -r destination/path/Season\ $STEP/Name (Year) S$1E$2.mkv

const (
	RecursiveFlag      = "recursive"
	RecursiveFlagShort = "r"
	VerboseFlag        = "verbose"
	VerboseFlagShort   = "v"
	ConfirmFlag        = "confirm"
	ConfirmFlagShort   = "c"
	StepFlag           = "step"
	StepFlagShort      = "s"
)

const regexConstants = ".*+?[]()|{}"

var rootCmd = &cobra.Command{
	Use:  "supalink <source path template> <destination path template>",
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		srcPath := args[0]

		if !strings.HasSuffix(srcPath, "$") {
			srcPath += "$"
		}

		srcExp := regexp.MustCompile(srcPath)
		destPath := args[1]

		fmt.Println("Source Path:", srcPath)
		fmt.Println("Destination Path:", destPath)

		rootDirectory := findRootDirectory(srcPath)
		fmt.Println("Root Directory:", rootDirectory)

		filepath.Walk(rootDirectory, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if matches := srcExp.FindStringSubmatch(path); matches != nil {
				destPathFilled := destPath
				replacer := regexp.MustCompile(`\$([0-9]+)`)

				for matchIndex, match := range matches[1:] {
					fmt.Println("Replacing match:", matchIndex+1, "with", match)
					destPathFilled = replacer.ReplaceAllStringFunc(destPathFilled, func(s string) string {
						index, err := strconv.Atoi(s[1:])
						if err != nil {
							return s
						}
						return matches[index]
					})
					fmt.Println("Destination Path after replacement:", destPathFilled)
				}

				fmt.Println("Matched:", path)
				return nil
			}

			return nil
		})

		return nil
	},
}

func findRootDirectory(path string) string {
	for i, c := range path {
		if strings.ContainsRune(regexConstants, c) {
			return filepath.Dir(path[:i])
		}
	}
	return filepath.Dir(path)
}

func main() {
	flags := rootCmd.Flags()
	flags.BoolP(RecursiveFlag, RecursiveFlagShort, false, "Search source path recursively")
	flags.BoolP(VerboseFlag, VerboseFlagShort, false, "Enable verbose output (good for debugging)")
	flags.BoolP(ConfirmFlag, ConfirmFlagShort, false, "Asks for user confirmation before creating symlinks")
	flags.StringArrayP(StepFlag, StepFlagShort, make([]string, 0), "Step number to break destination path into subdirectories")
	rootCmd.Execute()
}
