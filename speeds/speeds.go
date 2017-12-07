/*
* @Author: Juthathip Cheawcharn
* @Date:   2017-11-13
 */
package speeds

import (
	"fmt"
	"log"
	"strings"

	"github.com/Jutha1234/speedtest/misc"
	"github.com/Jutha1234/speedtest/print"
	"github.com/Jutha1234/speedtest/shttp"
)

var (
	// DefaultDLSizes defines the default download sizes
	DefaultDLSizes = []int{350, 500, 750, 1000, 1500, 2000, 2500, 3000, 3500, 4000}
	// DefaultULSizes defines the default upload sizes
	DefaultULSizes = []int{int(0.25 * 1024 * 1024), int(0.5 * 1024 * 1024), int(1.0 * 1024 * 1024), int(1.5 * 1024 * 1024), int(2.0 * 1024 * 1024)}
)

// SpeedTest defines a Speedtester client tester
type SpeedTest struct {
	Client   *sthttp.Client
	DLSizes  []int
	ULSizes  []int
	Quiet    bool
	Report   bool
	Debug    bool
	AlgoType string
}

func NewSpeed(client *sthttp.Client, dlsizes []int, ulsizes []int, quiet bool, report bool) *SpeedTest {
	return &SpeedTest{
		Client:  client,
		DLSizes: dlsizes,
		ULSizes: ulsizes,
		Quiet:   quiet,
		Report:  report,
	}
}

// Download will perform the "normal" speedtest download test
func (SpeedTest *SpeedTest) Download(server sthttp.Server) float64 {
	var urls []string
	var maxSpeed float64

	for size := range SpeedTest.DLSizes {
		url := server.URL
		splits := strings.Split(url, "/")
		baseURL := strings.Join(splits[1:len(splits)-1], "/")
		randomImage := fmt.Sprintf("random%dx%d.jpg", SpeedTest.DLSizes[size], SpeedTest.DLSizes[size])
		downloadURL := "http:/" + baseURL + "/" + randomImage
		urls = append(urls, downloadURL)
	}

	if !SpeedTest.Quiet && !SpeedTest.Report {
		log.Printf("Testing download speed")
	}

	for u := range urls {
		if SpeedTest.Debug {
			log.Printf("Download Test Run: %s\n", urls[u])
		}
		dlSpeed, err := SpeedTest.Client.DownloadSpeed(urls[u])
		if err != nil {
			log.Printf("Can't get download speed")
			dlSpeed = 0.0
			//log.Fatal(err)
		}
		//fmt.Println("dlSpeed :", dlSpeed)
		if !SpeedTest.Quiet && !SpeedTest.Debug && !SpeedTest.Report {
			fmt.Printf(".")
		}
		if SpeedTest.Debug {
			log.Printf("Dl Speed: %v\n", dlSpeed)
		}

		if dlSpeed > maxSpeed {
			maxSpeed = dlSpeed
		}

	}

	if !SpeedTest.Quiet && !SpeedTest.Report {
		fmt.Printf("\n")
	}

	return maxSpeed

}

// Upload runs a "normal" speedtest upload test
func (SpeedTest *SpeedTest) Upload(server sthttp.Server) float64 {
	// https://github.com/sivel/speedtest-cli/blob/master/speedtest-cli
	var ulsize []int
	var maxSpeed float64
	//var avgSpeed float64

	for size := range SpeedTest.ULSizes {
		ulsize = append(ulsize, SpeedTest.ULSizes[size])
	}

	if !SpeedTest.Quiet && !SpeedTest.Report {
		log.Printf("Testing upload speed")
	}

	//fmt.Println("ulsize", ulsize, " len ul size :", len(ulsize))

	for i := 0; i < len(ulsize); i++ {
		if SpeedTest.Debug {
			log.Printf("Upload Test Run: %v\n", i)
		}
		r := misc.Urandom(ulsize[i])
		ulSpeed, err := SpeedTest.Client.UploadSpeed(server.URL, "text/xml", r)
		if err != nil {
			fmt.Println("Error Post Upload ", err)
			ulSpeed = 0.0

		}
		if !SpeedTest.Quiet && !SpeedTest.Debug && !SpeedTest.Report {
			fmt.Printf(".")
		}
		if SpeedTest.Debug {
			log.Printf("Ul Amount: %v bytes\n", len(r))
			log.Printf("Ul Speed: %vMbps\n", ulSpeed)
		}

		if ulSpeed > maxSpeed {
			maxSpeed = ulSpeed
		}

	}

	if !SpeedTest.Quiet && !SpeedTest.Report {
		fmt.Printf("\n")
	}

	return maxSpeed
}

// FindServer will find a specific server in the servers list
func (SpeedTest *SpeedTest) FindServer(id string, serversList []sthttp.Server) sthttp.Server {
	var foundServer sthttp.Server
	for s := range serversList {
		if serversList[s].ID == id {
			foundServer = serversList[s]
		}
	}
	if foundServer.ID == "" {
		log.Fatalf("Cannot locate server Id '%s' in our list of speedtest servers!\n", id)
	}
	return foundServer
}

// ListServers prints a list of all "global" servers
func (SpeedTest *SpeedTest) ListServers(configURL string, serversURL string, blacklist []string) (err error) {
	if SpeedTest.Debug {
		fmt.Printf("Loading config from speedtest.net\n")
	}
	c, err := SpeedTest.Client.GetConfig()
	if err != nil {
		return err
	}
	SpeedTest.Client.Config = &c

	if SpeedTest.Debug {
		fmt.Printf("\n")
	}

	if SpeedTest.Debug {
		fmt.Printf("Getting servers list...")
	}
	allServers, err := SpeedTest.Client.GetServers()
	if err != nil {
		log.Fatal(err)
	}
	if SpeedTest.Debug {
		fmt.Printf("(%d) found\n", len(allServers))
	}
	for s := range allServers {
		server := allServers[s]
		print.Server(server)
	}
	return nil
}
