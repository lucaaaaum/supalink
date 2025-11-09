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
	Use:  "supalink <source path regex> <destination path template>",
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		srcPath := args[0]
		destPath := args[1]

		addStopSuffixToPattern(&srcPath)

		matchingPathsAndDestinations := getMatchingPathsAndDestinations(srcPath, destPath)

		if len(matchingPathsAndDestinations) == 0 {
			fmt.Println("No matching paths found.")
			return nil
		}

		for path, destination := range matchingPathsAndDestinations {
			fmt.Printf("%s -> %s\n", path, destination)
		}

		return nil
	},
}

func addStopSuffixToPattern(pattern *string) {
	if !strings.HasSuffix(*pattern, "$") {
		*pattern += "$"
	}
}

func getMatchingPathsAndDestinations(srcPath, destPath string) map[string]string {
	matchingPathsAndDestinations := make(map[string]string)
	rootDirectory := findRootDirectory(srcPath)

	srcExp := regexp.MustCompile(srcPath)

	filepath.Walk(rootDirectory, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if matches := srcExp.FindStringSubmatch(path); matches != nil {
			parameterMatches := matches[1:]
			destPathWithFilledParameters := getDestPathWithFilledParameters(destPath, parameterMatches)
			matchingPathsAndDestinations[path] = destPathWithFilledParameters
			return nil
		}

		return nil
	})
	return matchingPathsAndDestinations
}

func findRootDirectory(path string) string {
	for i, c := range path {
		if strings.ContainsRune(regexConstants, c) {
			return filepath.Dir(path[:i])
		}
	}
	return filepath.Dir(path)
}

func getDestPathWithFilledParameters(destPath string, parameterMatches []string) string {
	parameterExp := regexp.MustCompile(`\$([0-9]+)`)
	return parameterExp.ReplaceAllStringFunc(destPath, func(s string) string {
		index, err := strconv.Atoi(s[1:])
		if err != nil {
			return s
		}
		return parameterMatches[index-1]
	})
}

func main() {
	flags := rootCmd.Flags()
	flags.BoolP(RecursiveFlag, RecursiveFlagShort, false, "Search source path recursively")
	flags.BoolP(VerboseFlag, VerboseFlagShort, false, "Enable verbose output (good for debugging)")
	flags.BoolP(ConfirmFlag, ConfirmFlagShort, false, "Asks for user confirmation before creating symlinks")
	flags.StringArrayP(StepFlag, StepFlagShort, make([]string, 0), "Step number to break destination path into subdirectories")
	rootCmd.Execute()
}
