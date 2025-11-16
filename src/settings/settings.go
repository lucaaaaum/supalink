package settings

import (
	"fmt"
	"strconv"

	"github.com/spf13/pflag"
)

type Settings struct {
	Verbose bool
	Confirm bool
	DryRun  bool
	Steps   []int
	Format  string
}

func GetSettings(flags *pflag.FlagSet) (Settings, error) {
	settings := Settings{
		Verbose: flags.Changed(VerboseFlag) && flags.Lookup(VerboseFlag).Value.String() == "true",
		Confirm: flags.Changed(ConfirmFlag) && flags.Lookup(ConfirmFlag).Value.String() == "true",
		DryRun:  flags.Changed(DryRunFlag) && flags.Lookup(DryRunFlag).Value.String() == "true",
		Format:  flags.Lookup(FormatFlag).Value.String(),
		Steps:   make([]int, 0),
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
