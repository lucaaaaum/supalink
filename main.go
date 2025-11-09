package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
	DryRunFlag         = "dry-run"
	DryRunFlagShort    = "d"
)

const regexConstants = ".*+?[]()|{}"

type settings struct {
	Recursive bool
	Verbose   bool
	Confirm   bool
	Steps     []int
}

var rootCmd = &cobra.Command{
	Use:  "supalink <source path regex> <destination path template>",
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		srcPath := args[0]
		destPath := args[1]

		settings, err := getSettings(cmd.Flags())
		if err != nil {
			return err
		}

		addStopSuffixToPattern(&srcPath)

		matchingPathsAndDestinations := getMatchingPathsAndDestinations(srcPath, destPath, settings)

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

func getSettings(flags *pflag.FlagSet) (settings, error) {
	settings := settings{
		Recursive: flags.Changed(RecursiveFlag) && flags.Lookup(RecursiveFlag).Value.String() == "true",
		Verbose:   flags.Changed(VerboseFlag) && flags.Lookup(VerboseFlag).Value.String() == "true",
		Confirm:   flags.Changed(ConfirmFlag) && flags.Lookup(ConfirmFlag).Value.String() == "true",
		Steps:     make([]int, 0),
	}
	stepsAsStringArray, err := flags.GetStringArray(StepFlag)
	if err != nil {
		return settings, err
	}

	for _, stepAsString := range stepsAsStringArray {
		step, err := strconv.Atoi(stepAsString)
		if err != nil {
			return settings, fmt.Errorf("invalid step value: %s", stepAsString)
		}
		settings.Steps = append(settings.Steps, step)
	}

	return settings, err
}

func addStopSuffixToPattern(pattern *string) {
	if !strings.HasSuffix(*pattern, "$") {
		*pattern += "$"
	}
}

func getMatchingPathsAndDestinations(srcPath, destPath string, settings settings) map[string]string {
	matchingPathsAndDestinations := make(map[string]string)

	rootDirectory := findRootDirectory(srcPath)
	printIfVerbose(settings, "Searching in root directory: %s\n", rootDirectory)

	srcExp := regexp.MustCompile(srcPath)

	filepath.Walk(rootDirectory, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if matches := srcExp.FindStringSubmatch(path); matches != nil {
			printIfVerbose(settings, "Path matched: %s\n", path)
			parameterMatches := matches[1:]
			destPathWithFilledParameters := getDestPathWithFilledParameters(destPath, parameterMatches, settings)
			matchingPathsAndDestinations[path] = destPathWithFilledParameters
			return nil
		}

		return nil
	})

	return matchingPathsAndDestinations
}

func printIfVerbose(settings settings, message string, args ...any) {
	if settings.Verbose {
		fmt.Printf(message, args...)
	}
}

func findRootDirectory(path string) string {
	for i, c := range path {
		if strings.ContainsRune(regexConstants, c) {
			return filepath.Dir(path[:i])
		}
	}
	return filepath.Dir(path)
}

func getDestPathWithFilledParameters(destPath string, parameterMatches []string, settings settings) string {
	printIfVerbose(settings, "Filling parameters for destination path: %s\n", destPath)
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
	flags.BoolP(DryRunFlag, DryRunFlagShort, false, "Perform a trial run with no changes made")
	flags.StringArrayP(StepFlag, StepFlagShort, make([]string, 0), "Step number to break destination path into subdirectories")
	rootCmd.Execute()
}
