package track

import (
	"glorp-track/internal/setup"
	"glorp-track/internal/errors"
	"glorp-track/internal/utils"
)

func Init() {
	// Make sure setup is not running
	setup.Stop()
	err := utils.ExecCommand("systemctl start publish-cam")
	errors.ProcessError(errors.PUB_ERR, "unable to start publish-cam service for tracking", err)
}
