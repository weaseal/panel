// Package panel exports functions related to a Nanoleaf Canvas controller
package panel

import (
	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"
)

// NewPanel instantiates a new Panel
func NewPanel() *Panel {
	Panel := &Panel{}
	return Panel
}

// On turn on the light
func (panel *Panel) On() error {
	settings := &StateSettings{}
	onSetting := onSetting{Value: On}
	settings.On = &onSetting

	log.Debugf("Settings for On(): %+v", settings)
	if err := panel.apiClient.state(settings).put(); err != nil {
		return err
	}

	panel.Status.PowerState = On
	log.Infof("Panel state change: on")

	return nil
}

// Off turn on the light
func (panel *Panel) Off() error {
	settings := &StateSettings{}
	onSetting := onSetting{Value: Off}
	settings.On = &onSetting

	log.Debugf("Settings for Off(): %+v", settings)
	if err := panel.apiClient.state(settings).put(); err != nil {
		return err
	}

	panel.Status.PowerState = Off
	log.Infof("Panel state change: off")

	return nil
}

// SetBrightness adjusts panel brightness
func (panel *Panel) SetBrightness(brightness Brightness) error {

	settings := &StateSettings{}
	brightnessSetting := brightnessSetting{Value: brightness}
	settings.Brightness = &brightnessSetting

	log.Debugf("Settings for SetBrightness(): %+v", settings)
	if err := panel.apiClient.state(settings).put(); err != nil {
		return err
	}

	panel.Status.Brightness = brightness
	log.Infof("Brightness changed: %v", brightness)

	return nil
}

// SetTemperature adjusts panel colour by Kelvin temperature 1200 - 6500
func (panel *Panel) SetTemperature(temperature Temperature) error {

	settings := &StateSettings{}
	temperatureSettings := temperatureSetting{Value: temperature}
	settings.Temperature = &temperatureSettings

	log.Debugf("Settings for SetTemperature(): %+v", settings)
	if err := panel.apiClient.state(settings).put(); err != nil {
		return err
	}

	panel.Status.Temperature = temperature
	log.Infof("Temperature changed: %v", temperature)

	return nil
}

// SetHue adjusts panel colour by hue 0 - 360
func (panel *Panel) SetHue(hue Hue) error {

	settings := &StateSettings{}
	hueSettings := hueSetting{Value: hue}
	settings.Hue = &hueSettings

	log.Debugf("Settings for SetHue(): %+v", settings)
	if err := panel.apiClient.state(settings).put(); err != nil {
		return err
	}

	panel.Status.Hue = hue
	log.Infof("Hue changed: %v", hue)

	return nil
}

// SetSaturation adjusts panel saturation 0 - 100
func (panel *Panel) SetSaturation(saturation Saturation) error {

	settings := &StateSettings{}
	satSettings := saturationSetting{Value: saturation}
	settings.Saturation = &satSettings

	log.Debugf("Settings for SetHue(): %+v", settings)
	if err := panel.apiClient.state(settings).put(); err != nil {
		return err
	}

	panel.Status.Saturation = saturation
	log.Infof("Hue changed: %v", saturation)

	return nil
}

// Apply adjusts multiple settings in a single request.
// NOTE: As of firmware v1.5.0, the Canvas returns a 404 when multiple settings
// are sent -- although it does actually succeed at applying the change.  Given
// this bug, consider discarding the returned error here.
func (panel *Panel) Apply(settings *StateSettings) error {

	// if brightness is 0, only send that, otherwise the panel flickers
	if settings.Brightness != nil && settings.Brightness.Value == 0 {
		settings.Temperature = nil
		settings.Saturation = nil
		settings.Hue = nil
	}

	if (settings.Hue != nil || settings.Saturation != nil) && settings.Temperature != nil {
		log.Warnf("Temperature and hue/saturation are mutually exclusive. Dropping temperature. If you want to set temperature, don't set hue/sat.")
		settings.Temperature = nil
	}

	log.Debugf("Settings for Apply(): %+v", *settings)
	if err := panel.apiClient.state(settings).put(); err != nil {
		return err
	}

	log.Infof("Settings changed:")
	// print output about changes and update panel.Status
	if settings.Brightness != nil {
		log.Infof("\tBrightness changed: %v", settings.Brightness.Value)
		panel.Status.Brightness = settings.Brightness.Value
	}
	if settings.Temperature != nil {
		log.Infof("\tTemperature changed: %v", settings.Temperature.Value)
		panel.Status.Temperature = settings.Temperature.Value
	}
	if settings.Hue != nil {
		log.Infof("\tHue changed: %v", settings.Hue.Value)
		panel.Status.Hue = settings.Hue.Value
	}
	if settings.Saturation != nil {
		log.Infof("\tSaturation changed: %v", settings.Saturation.Value)
		panel.Status.Saturation = settings.Saturation.Value
	}

	return nil
}

// SetToken sets the token to be used when communicating with the API
func (panel *Panel) SetToken(token AuthToken) {
	panel.apiClient.token = token
}

// GetNewToken attempts to get a new auth token from the panel.  The panel must
// be in pairing mode for this to work. See notes under the API section 5.1,
// Authorization - https://forum.nanoleaf.me/docs/openapi
// the token is returned, and stored on the current panel object. The user
// should take care to save the returned token for future use.
func (panel *Panel) GetNewToken() (AuthToken, error) {

	if panel.apiClient.apiAddr == "" {
		return "", fmt.Errorf("could not get token, API address is not set")
	}

	token, err := panel.apiClient.newToken().post()
	if err != nil {
		return "", err
	}

	log.Debugf("got new token: %s", token)
	panel.apiClient.token = token

	return token, nil
}

// SetAPIAddr sets the address to be used when communicating with the API. Must
// be a valid URI.
func (panel *Panel) SetAPIAddr(apiAddr string) error {
	// check the API address
	if _, err := url.ParseRequestURI(apiAddr); err != nil {
		return fmt.Errorf("could not parse API addresss: %s", err.Error())
	}
	panel.apiClient.apiAddr = apiAddr
	return nil
}

// GetStateSettings returns the current state of the panel, omitting fields
// which are not really in use, such as hue when the colour mode is temperature,
// as they would fail to apply together. This is somewhat buggy, as the returned
// values are generally incorrect if the panel was configured by something other
// than direct API calls, such as via mobile app.
func (panel *Panel) GetStateSettings() (*StateSettings, error) {

	settings := &StateSettings{}

	state, err := panel.apiClient.state(settings).get()
	if err != nil {
		return nil, err
	}

	// if it's off, save no other information, as this confuses the panel
	settings.On = state.On
	if state.On.Value == Off {
		*settings.On = *state.On
		return settings, nil
	}

	settings.Brightness = state.Brightness

	if *state.ColourMode == ColourModeTemperature {
		settings.Temperature = state.Temperature
	}

	if *state.ColourMode == ColourModeHueSaturation {
		settings.Hue = state.Hue
		settings.Saturation = state.Saturation
	}

	log.Debug("got state:")
	if settings.ColourMode != nil {
		log.Debugf("colourmode: %v", state.ColourMode)
	}
	if settings.Hue != nil {
		log.Debugf("hue: %v", settings.Hue.Value)
	}
	if settings.On != nil {
		log.Debugf("on: %v", settings.On.Value)
	}
	if settings.Brightness != nil {
		log.Debugf("brightness: %v", settings.Brightness.Value)
	}
	if settings.Saturation != nil {
		log.Debugf("saturation: %v", settings.Saturation.Value)
	}
	if settings.Temperature != nil {
		log.Debugf("temperature: %v", settings.Temperature.Value)
	}

	return settings, err
}
