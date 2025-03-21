package main

import (
	"os"
	"fmt"
	"time"

	"github.com/joho/godotenv"

	"glorp-track/internal/errors"
	"glorp-track/internal/setup"
	"glorp-track/internal/track"
	"glorp-track/internal/network"
)

const IP_CHK_TIMEOUT = 4 // 4 Seconds between tries
const IP_CHK_ATTS = 5

func main() {
	envMap, err := godotenv.Read("./.env")
	errors.ProcessError(errors.FS_ERR, "unable to load .env", err)

	//ip, err := network.GetIP(envMap["WLAN_IFACE"])
	//errors.ProcessError(errors.NET_ERR, "unable to get ip", err)

	// So setup runs after Timeout * atts (20) of not finding net
	ip := attemptIp(IP_CHK_ATTS, IP_CHK_TIMEOUT, envMap["WLAN_IFACE"])

	if !ip { // No ip, move to setup mode
		fmt.Println("Run setup")
		setup.Init(envMap["HS_IFACE"], envMap["HS_GATEWAY_IP"])
	} else {
		fmt.Println("Run track")
		track.Init()
	}
	os.Exit(0)
}

func attemptIp(maxAtt, sleepSecs int, wlanIface string) bool {
	for i := 0; i < maxAtt; i++ {
		ip, err := network.GetIP(wlanIface)
		if err != nil { break }
		if ip != nil {
			return true
		}
		time.Sleep(time.Duration(sleepSecs) * time.Second)
	}
	return false
}
