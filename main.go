package main

import (
	"github.com/vessel-app/vessel-cli/cmd"
	"github.com/vessel-app/vessel-cli/internal/logger"
)

func main() {
	logger := logger.GetLogger()
	defer logger.Close()
	cmd.Execute()
}
