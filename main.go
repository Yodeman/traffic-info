package main

import (
    "fmt"
    "os"

    "github.com/yodeman/traffic-info/util"
)

func main() {
    apiKey := os.Getenv("GOOGLE_MAP_API_KEY")
    if !apiKey {
        log.Fatalln("Error retreiving Google map api key.")
    }

    // Verify command line argument inputs
    util.CheckArgs()

    // Fetch traffic information from Google.
    trafficResp, err := util.FetchTrafficInfo(apiKey)
    if err != nil {
        log.Fatalf("Fetch Traffic Info: %v", err)
    }

    err = printTrafficInfo(trafficResp)
    if err != nil {
        log.Fatalf("Print Traffic Info: %v", err)
    }
}
