package runtime

import (
	"github.com/xsbull/utils/logger"
)

// RecoveryFunction panic recovery
func RecoveryFunction() {
	if recoveryMessage := recover(); recoveryMessage != nil {
		logger.Log(logger.ErrorLevel, "panic info:%s", recoveryMessage)
	}
}
