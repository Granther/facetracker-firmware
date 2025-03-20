package main

import (
	//"os"
	//"fmt"

	//"github.com/joho/godotenv"

	//"glorp-track/internal/errors"
	"glorp-track/internal/setup"
	//"glorp-track/internal/track"
	//"glorp-track/internal/network"
)

func main() {
/*
	envMap, err := godotenv.Read("./.env")
	errors.ProcessError(errors.FS_ERR, "unable to load .env", err)

	ip, err := network.GetIP(envMap["WLAN_IFACE"])
	errors.ProcessError(errors.NET_ERR, "unable to get ip", err)

	if ip == nil { // No ip, move to setup mode
		fmt.Println("Run setup")
		//setup.Init()
	} else {
		fmt.Println("Run track")
		//track.Init()
	}
	os.Exit(0)
*/
	setup.Init()
}
