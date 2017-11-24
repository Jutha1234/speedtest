package print

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/Jutha1234/speedtest/sthttp"
)

// Server prints the results in "human" format
func Server(server sthttp.Server) {
	fmt.Printf("%-4s | %s (%s, %s)\n", server.ID, server.Sponsor, server.Name, server.Country)
}

// ServerReport prints restults in a machine useable format
func ServerReport(client *sthttp.Client, server sthttp.Server) {
	fmt.Printf("%s%s%s%s%s(%s,%s)%s", time.Now(), client.ReportChar, server.ID, client.ReportChar, server.Sponsor, server.Name, server.Country, client.ReportChar)
}

// EnvironmentReport is a debugging report helpful for debugging
func EnvironmentReport(client *sthttp.Client) {
	log.Printf("Env Report")
	log.Printf("-------------------------------\n")
	log.Printf("[User Environment]\n")
	log.Printf("Arch: %v\n", runtime.GOARCH)
	log.Printf("OS: %v\n", runtime.GOOS)
	log.Printf("IP: %v\n", client.Config.IP)
	log.Printf("Lat: %v\n", client.Config.Lat)
	log.Printf("Lon: %v\n", client.Config.Lon)
	log.Printf("ISP: %v\n", client.Config.Isp)
	log.Printf("Config: %s\n", client.SpeedtestConfig.ConfigURL)
	log.Printf("Servers: %s\n", client.SpeedtestConfig.ServersURL)
	log.Printf("User Agent: %s\n", client.SpeedtestConfig.UserAgent)
	log.Printf("-------------------------------\n")
	log.Printf("[args]\n")
	log.Printf("%#v\n", os.Args)
	log.Printf("--------------------------------\n")
}
