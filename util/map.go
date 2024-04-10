package util

import (
    "encoding/json"
    "net/http"
    "net/url"
    "strings"
    "text/template"
    "time"
)

var baseURL = "https://maps.googleapis.com/maps/api/distancematrix/json?"

const timeout = 10 * time.Second    // timeout for api request

// Google distance matrix api response structure
type DistanceMatrixResponse struct {
    DestinationAddr     []string            `json:"destination_addresses"`
    OriginAddr          []string            `json:"origin_addresses"`
    Rows                []DistanceMatrixRow
    Status              string
    ErrMsg              string              `json:"error_message,omitempty"`
}

type DistanceMatrixRow struct {
    Elements            []DistanceMatrixElement
}

type DistanceMatrixElement struct {
    Status              string
    Distance            DistanceTextValueObject `json:"omitempty"`
    Duration            DistanceTextValueObject `json:"omitempty"`
    TrafficDuration     DistanceTextValueObject `json:"duration_in_traffic,omitempty"`
    Fare                DistanceFare            `json:"omitempty"`
}

type DistanceTextValueObject struct {
    Text    string
    Value   float32
}

type DistanceFare struct {
    Currency    string
    Text        string
    Value       float32
}

// Text template used to report the traffic information response
// gotten from Google API.
const respTempl = `
Traffic Information:
--------------------

Origins:        {{.OriginAddr[0]}}

Destination:    {{.DestinationAddr[0]}}

{{range .Rows}}
    Distance:   {{.Elements[0].Distance.Text}}
    Duration:   {{.Elements[0].Duration.Text}}
{{end}}
`

var trafficInfo *template.Template

func init() {
    trafficInfo = template.Must(
        template.New("trafficInfo").Parse(respTempl))
}

// Fetch traffic information using the origin and destination.
func FetchTrafficInfo(apiKey string) (DistanceMatrixResponse, error) {
    URL := url.QueryEscape(fmt.Sprintf(
        "%slanguage=en&key=%s&origins=%s&destinations=%s",
        baseURL, apiKey, origin, destination))

    // Retry for 10 seconds in case of any error while fetching traffic
    // informatin from Google API.
    deadline := time.Now().Add(timeout)
    for tries := 0; time.Now().Before(deadline); tries++ {
        response, err := http.Get(URL)
        if err != nil {
            continue    // retry 
        }
        defer response.Body.Close()

        if response.StatusCode == http.StatusOK {
            var trafficResponse DistanceMatrixResponse
            err := json.NewDecoder(response.Body).Decode(trafficResponse)
            if err != nil {
                return DistanceMatrixResponse{}, fmt.Errorf("Error decoding response.")
            }
            return trafficResponse, nil
        }
    }
    return DistanceMatrixResponse{}, fmt.Errorf(
            "Error getting traffic information after %s",
            timeout)
}

func printTrafficInfo(trafficResponse DistanceMatrixResponse) error {
    writer := new(strings.Builder)

    if err := trafficInfo.Execute(writer, trafficResponse); err != nil {
        return fmt.Errorf("Error formatting traffic information response.%v", err)
    }

    fmt.Println(writer.String())
}
