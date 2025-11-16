package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"supalink/src/output"
	"supalink/src/settings"
	"supalink/src/utils"

	"github.com/spf13/cobra"
)

// objective: supalink src/path/.*S([0-9]{2})E([0-9]{2}).*.mkv -r destination/path/Season\ $STEP/Name (Year) S$1E$2.mkv

const regexConstants = ".*+?[]()|{}"

type stepManager struct {
	currentStep      int
	currentStepCount int
}

func (sm *stepManager) NextStep(settings settings.Settings) (int, int, error) {
	if len(settings.Steps) == 0 {
		return 0, 0, fmt.Errorf("no steps defined")
	}

	if sm.currentStep == 0 {
		sm.currentStep = 1
		sm.currentStepCount = 1
		return sm.currentStep, sm.currentStepCount, nil
	}

	if sm.currentStepCount >= settings.Steps[sm.currentStep-1] {
		if sm.currentStep >= len(settings.Steps) {
			return 0, 0, fmt.Errorf("exceeded the number of defined steps")
		}
		sm.currentStep++
		sm.currentStepCount = 1
	} else {
		sm.currentStepCount++
	}

	return sm.currentStep, sm.currentStepCount, nil
}

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

		matchingPathsAndDestinations := getMatchingPathsAndDestinations(srcPath, destPath, settings)

		if len(matchingPathsAndDestinations) == 0 {
			fmt.Println("No matching paths found.")
			return nil
		}

		createSymlinks(matchingPathsAndDestinations, settings)

		return nil
	},
}

func getMatchingPathsAndDestinations(srcPath, destPath string, settings settings.Settings) map[string]string {
	matchingPathsAndDestinations := make(map[string]string)

	rootDirectory := findRootDirectory(srcPath)

	srcExp := regexp.MustCompile(srcPath)

	stepManager := &stepManager{}

	filepath.Walk(rootDirectory, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if matches := srcExp.FindStringSubmatch(path); matches != nil {
			parameterMatches := matches[1:]
			destPathWithFilledParameters := getDestPathWithFilledParameters(destPath, parameterMatches, settings, stepManager)
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

func getDestPathWithFilledParameters(destPath string, parameterMatches []string, settings settings.Settings, stepManager *stepManager) string {
	parameterExp := regexp.MustCompile(`\$[0-9]+`)
	destPathWithFilledParameters := parameterExp.ReplaceAllStringFunc(destPath, func(s string) string {
		index, err := strconv.Atoi(s[1:])
		if err != nil {
			return s
		}
		return parameterMatches[index-1]
	})

	if len(settings.Steps) == 0 {
		return destPathWithFilledParameters
	}

	step, stepCount, _ := stepManager.NextStep(settings)

	stepCountParameterExp := regexp.MustCompile(`\$STEP_COUNT`)
	destPathWithFilledParameters = stepCountParameterExp.ReplaceAllStringFunc(destPathWithFilledParameters, func(s string) string {
		return strconv.Itoa(stepCount)
	})

	stepParameterExp := regexp.MustCompile(`\$STEP`)
	destPathWithFilledParameters = stepParameterExp.ReplaceAllStringFunc(destPathWithFilledParameters, func(s string) string {
		return strconv.Itoa(step)
	})

	return destPathWithFilledParameters
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
