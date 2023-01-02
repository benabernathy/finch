package main

import (
	"flag"
	"fmt"
	"github.com/benabernathy/finch/internal/config"
	"github.com/benabernathy/finch/internal/display"
	"github.com/benabernathy/finch/internal/flike"
	"github.com/benabernathy/finch/internal/paclient"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/umahmood/haversine"
	"log"
	"os"
	"strings"
)

type Args struct {
	ConfigPath string
}

var receiver paclient.Receiver
var receiverLoc haversine.Coord

func main() {

	var cfg config.Config

	args := processArgs(&cfg)

	if err := config.GetConfig(args.ConfigPath, &cfg); err != nil {
		log.Fatal(err)
	}

	if err := cleanenv.ReadConfig(args.ConfigPath, &cfg); err != nil {
		log.Fatal(err)
	}

	if err := paclient.GetReceiverInfo(cfg.PiAwareConfig.ReceiverUrl, &receiver); err != nil {
		log.Fatal(err)
	}

	receiverLoc = haversine.Coord{Lat: float64(receiver.Latitude), Lon: float64(receiver.Longitude)}

	log.Println(receiver)

	for {
		var aircraftResponse paclient.AircraftResponse
		if err := paclient.GetAircraft(cfg.PiAwareConfig.AircraftUrl, &aircraftResponse); err != nil {
			log.Println(err)

		}
		aircraft := processAircraftData(aircraftResponse)
		display.Display(aircraft, cfg)
	}
}

func processArgs(cfg interface{}) Args {
	var a Args

	f := flag.NewFlagSet("Example server", 1)
	f.StringVar(&a.ConfigPath, "c", "config.yml", "Path to configuration file")

	fu := f.Usage
	f.Usage = func() {
		fu()
		envHelp, _ := cleanenv.GetDescription(cfg, nil)
		fmt.Fprintln(f.Output())
		fmt.Fprintln(f.Output(), envHelp)
	}

	f.Parse(os.Args[1:])
	return a
}

func removeNoLoc(aircraft paclient.Aircraft) bool {
	if aircraft.Latitude == 0.0 || aircraft.Longitude == 0.0 {
		return false
	}
	return true
}

func computeDistance(aircraft paclient.Aircraft) paclient.Aircraft {
	acLoc := haversine.Coord{Lat: float64(aircraft.Latitude), Lon: float64(aircraft.Longitude)}
	mi, _ := haversine.Distance(receiverLoc, acLoc)

	aircraft.Distance = float32(mi)

	return aircraft
}

func sanitize(aircraft paclient.Aircraft) paclient.Aircraft {
	aircraft.Flight = strings.TrimSpace(aircraft.Flight)
	aircraft.Hex = strings.ToUpper(aircraft.Hex)

	return aircraft
}

func processAircraftData(response paclient.AircraftResponse) []paclient.Aircraft {

	aircraft := flike.Filter(response.Aircraft, removeNoLoc)
	aircraft = flike.Map(aircraft, computeDistance)
	aircraft = flike.Map(aircraft, sanitize)

	return aircraft
}
