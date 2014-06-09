package main

import (
	"github.com/cobratbq/flagtag"
	"log"
	"time"
)

type config struct {
	Sleep   time.Duration `flag:"s,1s,The amount of time to sleep."`
	Verbose bool          `flag:"v,false,Verbose output."`
}

func main() {
	// Prepare configuration
	var config config
	flagtag.MustConfigureAndParse(&config)
	// Start sleeping.
	if config.Verbose {
		log.Println("Sleeping for " + config.Sleep.String())
	}
	time.Sleep(config.Sleep)
	if config.Verbose {
		log.Println("Done.")
	}
}
