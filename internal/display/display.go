package display

import (
	"fmt"
	"github.com/benabernathy/finch/internal/config"
	"github.com/benabernathy/finch/internal/paclient"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"log"
	"math"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/devices/v3/ssd1306"
	"periph.io/x/devices/v3/ssd1306/image1bit"
	"periph.io/x/host/v3"
	"time"
)

/*
 Desired Screen:
 HEX / FLT
 ALT FT MSL (or FL) / SPEED KN
 LAT / LON / DIST NM
*/

func Display(aircraft []paclient.Aircraft, config config.Config) {

	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Use i2creg I²C bus registry to find the first available I²C bus.
	b, err := i2creg.Open("")
	if err != nil {
		log.Fatal(err)
	}
	defer b.Close()

	for _, ac := range aircraft {

		dev, err := ssd1306.NewI2C(b, &ssd1306.DefaultOpts)
		if err != nil {
			log.Fatalf("failed to initialize ssd1306: %v", err)
		}
		// Draw on it.
		img := image1bit.NewVerticalLSB(dev.Bounds())

		f := basicfont.Face7x13
		drawer := font.Drawer{
			Dst:  img,
			Src:  &image.Uniform{image1bit.On},
			Face: f,
			//Dot:  fixed.P(0, img.Bounds().Dy()-1-f.Descent),
		}
		size := 13.0
		dpi := 72.0
		spacing := 1.0

		y := 10 + int(math.Ceil(size*dpi/72))
		dy := int(math.Ceil(size * spacing * dpi / 72))

		drawer.Dot = fixed.Point26_6{
			X: fixed.I(2),
			Y: fixed.I(y),
		}

		displayString := ac.Hex + " / " + ac.Flight
		drawer.DrawString(displayString)
		y = y + dy

		// Format altitude
		var displayAlt string

		if config.Display.UseFlightLevel {
			if ac.Altitude >= config.Display.TransitionAltitude {
				fl := ac.Altitude / 1000
				displayAlt = fmt.Sprintf("FL%d0", fl)
			} else if config.Display.AltitudeUnits == "M" {
				fl := int(float64(ac.Altitude) / 3.28)
				displayAlt = fmt.Sprintf("%.0f M", fl)
			} else {
				displayAlt = fmt.Sprintf("%d FT", ac.Altitude)
			}
		}

		var displaySpd string

		if config.Display.SpeedUnits == "KPH" {
			spd := ac.Speed * 1.85
			displaySpd = fmt.Sprintf("%.0f KPH", spd)
		} else if config.Display.SpeedUnits == "MPH" {
			spd := ac.Speed * 1.15078
			displaySpd = fmt.Sprintf("%.0f MPH", spd)
		} else {
			displaySpd = fmt.Sprintf("%.0f KN", ac.Speed)
		}

		//displaySpd := fmt.Sprintf("%.0f", math.Round(float64(ac.Speed)))
		displayLat := fmt.Sprintf("%.2f", ac.Latitude)
		displayLon := fmt.Sprintf("%.2f", ac.Longitude)

		var displayDistance string

		if config.Display.DistanceUnits == "KM" {
			dist := ac.Distance * 1.60934
			displayDistance = fmt.Sprintf("%.f KM", dist)
		} else if config.Display.DistanceUnits == "MI" {
			displayDistance = fmt.Sprintf("%.f MI", ac.Distance)
		} else {
			dist := ac.Distance * 0.868976
			displayDistance = fmt.Sprintf("%.f NM", dist)
		}

		var displayText = []string{
			displayAlt + " / " + displaySpd,
			displayLat + " / " + displayLon,
			displayDistance,
		}

		for _, s := range displayText {
			drawer.Dot = fixed.P(2, y)
			drawer.DrawString(s)
			y += dy
		}

		log.Println(displayString)

		if err := dev.Draw(dev.Bounds(), img, image.Point{}); err != nil {
			log.Fatal(err)
		}

		time.Sleep(5 * time.Second)
	}
}
