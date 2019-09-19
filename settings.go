package panel

import log "github.com/sirupsen/logrus"

// SetBrightness is a helper function that allows updating a StateSettings
// object with shorthand eg "settings.SetBrightness(20)"
func (settings *StateSettings) SetBrightness(brightness Brightness) {
	brightnessSetting := brightnessSetting{Value: brightness}
	settings.Brightness = &brightnessSetting
}

// SetTemperature is a helper function that allows updating a StateSettings
// object with shorthand eg "settings.SetTemperature(1200)"
func (settings *StateSettings) SetTemperature(temperature Temperature) {
	temperatureSetting := temperatureSetting{Value: temperature}
	settings.Temperature = &temperatureSetting
}

// SetHue is a helper function that allows updating a StateSettings
// object with shorthand eg "settings.SetHue(20)"
func (settings *StateSettings) SetHue(hue Hue) {
	hueSetting := hueSetting{Value: hue}
	settings.Hue = &hueSetting
}

// SetSaturation is a helper function that allows updating a StateSettings
// object with shorthand eg "settings.SetSaturation(20)"
func (settings *StateSettings) SetSaturation(saturation Saturation) {
	satSettings := saturationSetting{Value: saturation}
	settings.Saturation = &satSettings
}

// Validate adjusts settings to required min/max
func (settings *StateSettings) Validate() {

	if settings != nil {
		if settings.Temperature != nil {
			if settings.Temperature.Value < MinimumTemperature {
				log.Warnf("requested temperature %d is below minimum %d, adjusting", settings.Temperature.Value, MinimumTemperature)
				settings.Temperature.Value = MinimumTemperature
			}

			if settings.Temperature.Value > MaximumTemperature {
				log.Warnf("requested temperature %d is above maximum %d, adjusting", settings.Temperature.Value, MaximumTemperature)
				settings.Temperature.Value = MaximumTemperature
			}
		}

		if settings.Brightness != nil {
			if settings.Brightness.Value < MinimumBrightness {
				log.Warnf("requested brightness %d is below minimum %d, adjusting", settings.Brightness.Value, MinimumBrightness)
				settings.Brightness.Value = MinimumBrightness
			}

			if settings.Brightness.Value > MaximumBrightness {
				log.Warnf("requested brightness %d is above maximum %d, adjusting", settings.Brightness.Value, MaximumBrightness)
				settings.Brightness.Value = MaximumBrightness
			}
		}

		if settings.Saturation != nil {
			if settings.Saturation.Value < MinimumSaturation {
				log.Warnf("requested saturation %d is below minimum %d, adjusting", settings.Saturation.Value, MinimumSaturation)
				settings.Saturation.Value = MinimumSaturation
			}

			if settings.Saturation.Value > MaximumSaturation {
				log.Warnf("requested saturation %d is above maximum %d, adjusting", settings.Saturation.Value, MaximumSaturation)
				settings.Saturation.Value = MaximumSaturation
			}
		}

		if settings.Hue != nil {
			if settings.Hue.Value < MinimumHue {
				log.Warnf("requested hue %d is below minimum %d, adjusting", settings.Hue.Value, MinimumHue)
				settings.Hue.Value = MinimumHue
			}

			if settings.Hue.Value > MaximumHue {
				log.Warnf("requested hue %d is above maximum %d, adjusting", settings.Hue.Value, MaximumHue)
				settings.Hue.Value = MaximumHue
			}
		}
	}
}
