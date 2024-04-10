package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"
)

var baseURL = "https://maps.googleapis.com/maps/api/distancematrix/json?"

const timeout = 10 * time.Second // timeout for api request

// Google distance matrix api response structure
type DistanceMatrixResponse struct {
	DestinationAddr []string `json:"destination_addresses"`
	OriginAddr      []string `json:"origin_addresses"`
	Rows            []DistanceMatrixRow
	Status          string
	ErrMsg          string `json:"error_message,omitempty"`
}

type DistanceMatrixRow struct {
	Elements []DistanceMatrixElement
}

type DistanceMatrixElement struct {
	Status          string
	Distance        DistanceTextValueObject `json:"distance,omitempty"`
	Duration        DistanceTextValueObject `json:"duration,omitempty"`
	TrafficDuration DistanceTextValueObject `json:"duration_in_traffic,omitempty"`
	Fare            DistanceFare            `json:"fare,omitempty"`
}

type DistanceTextValueObject struct {
	Text  string
	Value float32
}

type DistanceFare struct {
	Currency string
	Text     string
	Value    float32
}

// Text template used to report the traffic information response
// gotten from Google API.
const respTempl = `
Traffic Information:
--------------------

Origin:         {{index .OriginAddr 0}}

Destination:    {{index .DestinationAddr 0}}

{{range .Rows}}
    *Distance:  {{(index .Elements 0).Distance.Text}}

    *Duration:  {{(index .Elements 0).Duration.Text}}

    *Status:    {{(index .Elements 0).Status}}
{{end}}
`

var trafficInfo *template.Template

func init() {
	trafficInfo = template.Must(
		template.New("trafficInfo").Parse(respTempl))
}

// Fetch traffic information using the origin and destination.
func FetchTrafficInfo(apiKey string) (DistanceMatrixResponse, error) {
	escapedOrigin := url.QueryEscape(origin)
	escapedDestination := url.QueryEscape(destination)
	URL := fmt.Sprintf(
		"%slanguage=en&key=%s&origins=%s&destinations=%s",
		baseURL, apiKey, escapedOrigin, escapedDestination)

	// Retry for 10 seconds in case of any error while fetching traffic
	// informatin from Google API.
	deadline := time.Now().Add(timeout)
	for tries := 0; time.Now().Before(deadline); tries++ {
		response, err := http.Get(URL)
		if err != nil {
			continue // retry
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusOK {
			var trafficResponse DistanceMatrixResponse
			err := json.NewDecoder(response.Body).Decode(&trafficResponse)
			if err != nil {
				return DistanceMatrixResponse{}, fmt.Errorf(
					"Error decoding response. %v", err)
			}
			return trafficResponse, nil
		}
	}
	return DistanceMatrixResponse{}, fmt.Errorf(
		"Error getting traffic information after %s.",
		timeout)
}

// Print a formatted version of the traffic information received
// from Google map API.
func PrintTrafficInfo(trafficResponse DistanceMatrixResponse) error {
	writer := new(strings.Builder)

	if err := trafficInfo.Execute(writer, trafficResponse); err != nil {
		return fmt.Errorf("Error formatting traffic information response.%v", err)
	}

	_, err := fmt.Println(writer.String())
	return err
}
