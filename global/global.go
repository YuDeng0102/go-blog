package global

import (
	"server/config"

	"go.uber.org/zap"
)

var (
	Config *config.Config
	Log    *zap.Logger
)
