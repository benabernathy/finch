package flike

import (
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"log"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/devices/v3/ssd1306"
	"periph.io/x/devices/v3/ssd1306/image1bit"
	"periph.io/x/host/v3"
)

func main() {

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
		Dot: fixed.P(0, 0),
	}
	drawer.DrawString("ACED21 / ASA211")

	if err := dev.Draw(dev.Bounds(), img, image.Point{}); err != nil {
		log.Fatal(err)
	}
}