package util

import (
    "fmt"
    "flag"
    "log"
    "strconv"
    "strings"
)

// command line input arguments destinations
var (
    origin      string  // starting coordinate.
    destination string  // ending coordinate.
)

func init() {
    flag.StringVar(
            &origin, "origin", "",
            "Travel starting coordinate i.e latitude,longitude")
    flag.StringVar(
            &destination, "destination", "",
            "Travel ending coordinate i.e. latitude,logitude")
    flag.Parse()
}

// verify the status of the command line arguments.
func CheckArgs() {
    if origin=="" || destination=="" {
        log.Fatalln("Command line argument `origin` and `destination` are required.")
    }

    if valid, err := verifyCoordinate(origin); !valid {
        log.Fatalln("origin coordinate: %v", err)
    }

    if valid, err := verifyCoordinate(destination); !valid {
        log.Fatalf("destination coordinate: %v\n", err)
    }
}

// validates the command line latitude/longitude input
func verifyCoordinate(coordinate string) (bool, error) {
    coord := strings.Split(coordinate, ",")
    if len(coord) != 2 {
        return false, fmt.Errorf(
                        "Error decoding coordinate: %s\n Expected latitude,longitude",
                        coordinate)
    }
    // the received latitude must be convertible to float and it
    // should be between -90 and 90.
    latitude, err := strconv.ParseFloat(coord[0], 32)
    if err != nil {
        return false, fmt.Errorf(
                        "Invalid latitude: %s.Value should be convertible to float.",
                        coord[0])
    }
    if latitude > 90 || latitude < -90 {
        return false, fmt.Errorf(
                        "Invalid latitude: %f.Value should be between -90 and 90.",
                        latitude)
    }

    // the received longitude must be convertible to float and it
    // should be between -180 and 180.
    longitude, err := strconv.ParseFloat(coord[1], 32)
    if err != nil {
        return false, fmt.Errorf(
                        "Invalid longitude: %s.Value should be convertible to float.",
                        coord[0])
    }
    if longitude > 180 || longitude < -180 {
        return false, fmt.Errorf(
                        "Invalid longitude: %f.Value should be between -180 and 180.",
                        latitude)
    }

    return true, nil
}
