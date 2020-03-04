package voalarm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Client represents the main parameters to generate an alarm in VictorOps plattform
type Client struct {
	Endpoint string
	APIkey   string
}

// Alarm represents the properties needed to set an alarm
type Alarm struct {
	MessageType       MessageType `json:"message_type"`
	EntityState       MessageType `json:"entity_state"`
	EntityID          string      `json:"entity_id"`
	EntityDisplayName string      `json:"entity_display_name"`
	StateMessage      string      `json:"state_message"`
	StateStartTime    string      `json:"state_start_time"`
}

// Response represents the response of VictorOps plattform
type Response struct {
	Message  string `json:"message"`
	Result   string `json:"result"`
	EntityID string `json:"entity_id"`
}

// MessageType is a representation of a string type
type (
	MessageType string
)

// const are the different values to set the MessageType parameter
const (
	Info            MessageType = "INFO"
	Recovery        MessageType = "RECOVERY"
	Acknowledgement MessageType = "ACKNOWLEDGEMENT"
	Warning         MessageType = "WARNING"
	Critical        MessageType = "CRITICAL"
)

// NewClient is a method to instantiate an API object setting the main parameters
func NewClient(apikey string) *Client {
	if len(apikey) == 0 {
		apikey = "2f616629-de63-4162-bb6f-11966bbb538d/test"
	}

	client := &Client{
		APIkey:   apikey,
		Endpoint: "https://alert.victorops.com/integrations/generic/20131114/alert",
	}
	return client
}

// SendAlarm is the main method to generate an alarm.
// Needs a MessageType parameter and the error that we need to log.
// Returns the response of VictorOps plattform and nil if success
func (a *Client) SendAlarm(entity string, ms MessageType, err error) (*Response, error) {
	if len(a.APIkey) == 0 || len(a.Endpoint) == 0 {
		a = NewClient("")
	}

	alarm := Alarm{
		MessageType:       ms,
		EntityState:       ms,
		EntityID:          entity,
		EntityDisplayName: entity,
		StateMessage:      err.Error(),
		StateStartTime:    time.Now().Format("2006-01-02 15:04:05"),
	}

	var r Response

	bytevalues, err := json.Marshal(alarm)
	if err != nil {
		return &Response{
			Result:  "failure",
			Message: fmt.Sprintf("Error marshaling Alarm struct %v", err),
		}, err
	}

	url := fmt.Sprint(a.Endpoint, "/", a.APIkey)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bytevalues))
	if err != nil {
		return &Response{
			Result:  "failure",
			Message: fmt.Sprintf("Unable to send alarm to VictorOps due to: %v", err),
		}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &Response{
			Result:  "failure",
			Message: resp.Status,
		}, fmt.Errorf("Error. Verify API key and enpoint URL: %v", resp.Status)
	}

	data, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(data, &r)

	return &r, nil
}
