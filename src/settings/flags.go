package settings

import (
	"supalink/src/output"

	"github.com/spf13/pflag"
)

const (
	VerboseFlag      = "verbose"
	VerboseFlagShort = "v"
	ConfirmFlag      = "confirm"
	ConfirmFlagShort = "c"
	StepFlag         = "step"
	StepFlagShort    = "s"
	DryRunFlag       = "dry-run"
	DryRunFlagShort  = "d"
	FormatFlag       = "format"
	FormatFlagShort  = "f"
)

func SetupFlags(flags *pflag.FlagSet) {
	flags.BoolP(VerboseFlag, VerboseFlagShort, false, "Enable verbose output (good for debugging)")
	flags.BoolP(ConfirmFlag, ConfirmFlagShort, false, "Asks for user confirmation before creating symlinks")
	flags.BoolP(DryRunFlag, DryRunFlagShort, false, "Perform a trial run with no changes made")
	flags.StringArrayP(StepFlag, StepFlagShort, make([]string, 0), "Step number to break destination path into subdirectories")
	flags.StringP(FormatFlag, FormatFlagShort, output.TreeFormat, "Format of the destination path")
}
