package runtime

import (
	"github.com/commcos/utils/logger"
)

// RecoveryFunction panic recovery
func RecoveryFunction() {
	if recoveryMessage := recover(); recoveryMessage != nil {
		logger.Log(logger.ErrorLevel, "panic info:%s", recoveryMessage)
	}
}
