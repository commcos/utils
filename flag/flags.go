package flag

import (
	"github.com/spf13/pflag"
	"github.com/xsbull/utils/logger"
)

// PrintFlags logs the flags in the flagset
func PrintFlags(flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		logger.Log(logger.DebugLevel, "FLAG: --%s=%q", flag.Name, flag.Value)
	})
}
