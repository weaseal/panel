package panel

import (
	"net/http"
)

// Panel is a Nanoleaf panel
type Panel struct {
	Status    PanelStatus
	apiClient PanelAPIClient
}

// PanelStatus contains known information about the panel's operation
type PanelStatus struct {
	PowerState  PowerState
	Brightness  Brightness
	Temperature Temperature
	Hue         Hue
	Saturation  Saturation
}

// PanelAPIClient sends requests to a panel
type PanelAPIClient struct {
	apiAddr    string // "http://some.ip.and:port"
	token      AuthToken
	httpClient http.Client
}

// Brightness 0 - 100
type Brightness uint8

// PowerState is the on/off value for the panel light
type PowerState bool

// Temperature colour in Kelvin 1200 - 6500
type Temperature uint16

// Hue colour 0 - 360
type Hue uint16

// Saturation for colour 0 - 100
type Saturation uint16

// ColourMode describes the current method for setting colour, either "ct" or
// "hs", for temperature or hue/saturation
type ColourMode string

// APIPath is the path after the auth-token, eg "state". See related consts
type APIPath string

// AuthToken allows you to interact with the panel
type AuthToken string

// PanelRequest contains all the information we need for a complete API request
type panelRequest struct {
	method                            string // PUT, POST, GET, etc
	client                            http.Client
	tokenizedAddress                  string
	tokenizedAddressWithCensoredToken string
	body                              interface{}
}

// StateSettings settings for the panel display
type StateSettings struct {
	Brightness  *brightnessSetting  `json:"brightness,omitempty"`
	Temperature *temperatureSetting `json:"ct,omitempty"`
	Hue         *hueSetting         `json:"hue,omitempty"`
	Saturation  *saturationSetting  `json:"saturation,omitempty"`
	ColourMode  *ColourMode         `json:"colorMode,omitempty"`
	On          *onSetting          `json:"on,omitempty"`
}

type newToken struct {
	AuthToken AuthToken `json:"auth_token"`
}

// onSetting is a struct to help marshal JSON
type onSetting struct {
	Value PowerState `json:"value"`
}

// brightnessSetting is a struct to help marshal JSON
type brightnessSetting struct {
	Value Brightness `json:"value"`
}

// temperatureSetting is a struct to help marshal JSON
type temperatureSetting struct {
	Value Temperature `json:"value"`
}

// hueSetting is a struct to help marshal JSON
type hueSetting struct {
	Value Hue `json:"value"`
}

// SaturationSetting is a struct to help marshal JSON
type saturationSetting struct {
	Value Saturation `json:"value"`
}
