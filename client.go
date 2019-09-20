package panel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// state takes settings to be sent. Do not perform any logic here which could
// error, because State cannot return an error, as it is a chained method.
func (p apiClient) state(settings *StateSettings) panelRequest {

	settings.Validate()

	request := panelRequest{
		tokenizedAddress:                  fmt.Sprintf("%s/api/v1/%s/%s", p.apiAddr, p.token, APIPathState),
		tokenizedAddressWithCensoredToken: fmt.Sprintf("%s/api/v1/%s/%s", p.apiAddr, "$token", APIPathState),
		client:                            p.httpClient,
		body:                              settings,
	}

	return request
}

// newToken deals with credentials
func (p apiClient) newToken() panelRequest {

	request := panelRequest{
		tokenizedAddress:                  fmt.Sprintf("%s/api/v1/%s", p.apiAddr, APIPathNew),
		tokenizedAddressWithCensoredToken: fmt.Sprintf("%s/api/v1/%s", p.apiAddr, APIPathNew),
		client:                            p.httpClient,
	}

	return request
}

// put takes a panelRequest and sends it
func (request panelRequest) put() error {
	request.method = http.MethodPut
	if _, err := request.do(); err != nil {
		return err
	}

	return nil
}

// post takes a panelRequest and sends it
func (request panelRequest) post() (AuthToken, error) {
	request.method = http.MethodPost
	body, err := request.do()
	if err != nil {
		return "", fmt.Errorf("failed to send new token request: %s", err.Error())
	}

	token := newToken{}
	err = json.Unmarshal(body, &token)
	if err != nil {
		return "", fmt.Errorf("failed to process new token JSON: %s", err.Error())
	}

	return token.AuthToken, nil
}

// get gets panel current state
func (request panelRequest) get() (*StateSettings, error) {
	request.method = http.MethodGet
	body, err := request.do()
	if err != nil {
		return nil, err
	}

	settings := &StateSettings{}
	err = json.Unmarshal(body, settings)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

// do sends the API request
func (request panelRequest) do() ([]byte, error) {

	bodyJSON, err := json.MarshalIndent(request.body, "", "  ")
	if err != nil {
		return nil, err
	}

	log.Debugf("Sending to: %s", request.tokenizedAddressWithCensoredToken)
	log.Debugf("Sending JSON: %s", bodyJSON)

	httpRequest, err := http.NewRequest(request.method, request.tokenizedAddress, bytes.NewBuffer(bodyJSON))
	if err != nil {
		return nil, err
	}

	resp, err := request.client.Do(httpRequest)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Failed to send request: %s", resp.Status)
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("could not read response body")
		return nil, err
	}

	if string(responseBody) != "" {
		return responseBody, nil
	}

	log.Debugf("do() got nil body response")
	return nil, nil
}
