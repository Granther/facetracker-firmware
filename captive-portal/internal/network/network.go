package network

// Discover local networks and make permanent connections to them across boot

import (
	"log"
    "fmt"

	"github.com/schollz/wifiscan"
)

func DiscoverNets() {
        wifis, err := wifiscan.Scan()
        if err != nil {
                log.Fatal(err)
        }
        for _, w := range wifis {
                fmt.Println(w.SSID, w.RSSI)
        }
}
