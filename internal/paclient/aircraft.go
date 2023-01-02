package paclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Aircraft struct {
	Hex       string  `json:"hex"`
	Flight    string  `json:"flight"`
	Altitude  int     `json:"alt_baro"`
	Speed     float32 `json:"gs"`
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"lon"`
	Distance  float32
}

type AircraftResponse struct {
	Aircraft []Aircraft `json:"aircraft"`
}

type AircraftClientError struct {
	Message string
	Url     string
	Err     error
	Data    string
}

func (e *AircraftClientError) Error() string {

	return fmt.Sprintf("%s: %s, %v, %s", e.Message, e.Url, e.Err, e.Data)
}

func GetAircraft(url string, aircraftResponse *AircraftResponse) error {

	response, getError := http.Get(url)

	if getError != nil {

		return &AircraftClientError{
			Message: "Error while communicating with PiAware aircraft service",
			Url:     url,
			Err:     getError,
		}
	}

	buf := new(strings.Builder)
	if _, err := io.Copy(buf, response.Body); err != nil {

		return &AircraftClientError{
			Message: "Error while getting response data from PiAware aircraft service",
			Url:     url,
			Err:     err,
		}
	}

	if err := response.Body.Close(); err != nil {

		return &AircraftClientError{
			Message: "Error while closing response body",
			Url:     url,
			Err:     err,
		}

	}

	if err := json.Unmarshal([]byte(buf.String()), aircraftResponse); err != nil {

		return &AircraftClientError{
			Message: "Error while unmarshalling response body",
			Url:     url,
			Err:     err,
			Data:    buf.String(),
		}
	}

	return nil
}
