package main

import (
	"os"
	"fmt"

	"github.com/joho/godotenv"

	"glorp-track/internal/errors"
	"glorp-track/internal/setup"
	"glorp-track/internal/track"
	"glorp-track/internal/network"
)

func main() {
	envMap, err := godotenv.Read("./.env")
	errors.ProcessError(errors.FS_ERR_COD, "unable to load .env", err)
	fmt.Println("Env map: ", envMap)

	ip, err := network.GetIP()
	errors.ProcessError(errors.FS_ERR_COD, "unable to load .env", err)

	if ip == nil { // No ip, move to setup mode
		err = setup.Init()
		errors.ProcessError(errors.UNK_ERR_COD, "setup mode init failed", err)
		os.Exit(0)
	} else {
		err = track.Init()
		errors.ProcessError(errors.UNK_ERR_COD, "tracking mode init failed", err)
		os.Exit(0)
	}
}
