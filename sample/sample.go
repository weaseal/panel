package main

import (
	"time"

	"github.com/weaseal/panel"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)

	// Set up a new panel
	p := panel.NewPanel()

	// Set the API address of the panel
	if err := p.SetAPIAddr("http://192.168.1.131:16021"); err != nil {
		log.Fatal(err)
	}

	// You can set the token if you know it via p.SetToken(), or hold the power
	// button on the panel to set it in pairing mode, and use p.GetToken().
	// Be sure to save the returned token (unlike in the example here) so that you
	// can set it later with p.SetToken(), otherwise it is saved in-memory only.
	if _, err := p.GetNewToken(); err != nil {
		log.Fatalf("could not get new token: %s", err.Error())
	}

	// store the panel state so it can be restored later
	originalSettings, err := p.GetStateSettings()
	if err != nil {
		log.Fatalf("could not get current status: %s", err.Error())
	}

	// turn on the panel
	if err := p.On(); err != nil {
		log.Fatalf("failed to turn on light: %s", err.Error())
	}

	time.Sleep(1 * time.Second)

	// turn off the panel
	p.Off()

	// methods that operate on panel.Panel modify a single setting immediately
	p.SetBrightness(20)
	p.SetTemperature(6000)
	time.Sleep(1 * time.Second)

	p.SetBrightness(40)
	time.Sleep(1 * time.Second)

	p.SetBrightness(60)
	p.SetTemperature(1200)
	time.Sleep(1 * time.Second)

	p.SetBrightness(80)
	time.Sleep(1 * time.Second)

	p.SetBrightness(100)
	time.Sleep(1 * time.Second)

	p.SetBrightness(0)
	time.Sleep(1 * time.Second)

	// you can optionally use a settings object to change multiple settings in a
	// single request with p.Apply(settings)
	settings := &panel.StateSettings{}
	settings.SetBrightness(20)
	settings.SetSaturation(100)
	// don't action any returned error from Apply(), due to a bug with Canvas. See
	// the comment on the function definition for more information
	p.Apply(settings)

	// do a rainbow
	for i := 0; i <= 20; i++ {
		settings.SetHue(panel.Hue(i * 18))
		p.Apply(settings)
		time.Sleep(10 * time.Millisecond)
	}

	if err := p.Apply(originalSettings); err != nil {
		log.Warnf("could not restore original settings: %s", err.Error())
	}

}
