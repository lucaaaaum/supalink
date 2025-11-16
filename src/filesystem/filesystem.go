package filesystem

import (
	"io/fs"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"supalink/src/settings"
	"supalink/src/step"
	"supalink/src/utils"
)

func GetMatchingPathsAndDestinations(srcPath, destPath string, settings settings.Settings) map[string]string {
	matchingPathsAndDestinations := make(map[string]string)

	rootDirectory := findRootDirectory(srcPath)

	srcExp := regexp.MustCompile(srcPath)

	stepManager := &step.StepManager{}

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
		if strings.ContainsRune(utils.RegexConstants, c) {
			return filepath.Dir(path[:i])
		}
	}
	return filepath.Dir(path)
}

func getDestPathWithFilledParameters(destPath string, parameterMatches []string, settings settings.Settings, stepManager *step.StepManager) string {
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
