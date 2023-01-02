package paclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Receiver struct {
	Version   string  `json:"version"`
	Refresh   int     `json:"refresh"`
	History   int     `json:"history"`
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"lon"`
}

type ReceiverClientError struct {
	Message string
	Url     string
	Err     error
}

func (e *ReceiverClientError) Error() string {

	return fmt.Sprintf("%s: %s, %v", e.Message, e.Url, e.Err)
}

func GetReceiverInfo(url string, receiver *Receiver) error {

	response, getError := http.Get(url)

	if getError != nil {

		return &ReceiverClientError{
			Message: "Error while communicating with PiAware receiver service",
			Url:     url,
			Err:     getError,
		}

	}

	buf := new(strings.Builder)
	if _, err := io.Copy(buf, response.Body); err != nil {

		return &ReceiverClientError{
			Message: "Error while getting response data from PiAware receiver service",
			Url:     url,
			Err:     err,
		}
	}

	if err := response.Body.Close(); err != nil {

		return &ReceiverClientError{
			Message: "Error while closing response body",
			Url:     url,
			Err:     err,
		}

	}

	if err := json.Unmarshal([]byte(buf.String()), receiver); err != nil {
		return &ReceiverClientError{
			Message: "Error while unmarshalling response body",
			Url:     url,
			Err:     err,
		}
	}

	return nil
}
