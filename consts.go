package panel

const (
	// On is true
	On PowerState = true
	// Off is false
	Off PowerState = false

	// APIPathState is the path for the power on/off state
	APIPathState APIPath = "state"
	// APIPathNew is the path for token credentials
	APIPathNew APIPath = "new"

	// MinimumTemperature lowest colour temperature in Kelvin
	MinimumTemperature Temperature = 1200
	// MaximumTemperature highest colour temperature in Kelvin
	MaximumTemperature Temperature = 6500
	// MinimumBrightness lowest possible Brightness
	MinimumBrightness Brightness = 0
	// MaximumBrightness lowest possible Brightness
	MaximumBrightness Brightness = 100
	// MinimumHue lowest possible Hue
	MinimumHue Hue = 0
	// MaximumHue lowest possible Hue
	MaximumHue Hue = 360
	// MinimumSaturation lowest possible saturation
	MinimumSaturation Saturation = 0
	// MaximumSaturation highest possible saturation
	MaximumSaturation Saturation = 100
	// ColourModeTemperature is the API's colourMode name for temperature
	ColourModeTemperature ColourMode = "ct"
	// ColourModeHueSaturation is the API's colourMode name for hue/saturation
	ColourModeHueSaturation ColourMode = "hs"
)
