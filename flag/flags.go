package flag

import (
	"github.com/commcos/utils/logger"
	"github.com/spf13/pflag"
)

// PrintFlags logs the flags in the flagset
func PrintFlags(flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		logger.Log(logger.DebugLevel, "FLAG: --%s=%q", flag.Name, flag.Value)
	})
}
