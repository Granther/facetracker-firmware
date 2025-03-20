package app

import (
	"captive-portal/internal/web"
)

func Start() error {
	r := web.SetupRouter()
	r.Run(":9090")
	return nil
}
