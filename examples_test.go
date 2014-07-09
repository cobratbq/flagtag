package flagtag

import (
	"fmt"
	"log"
	"time"
)

func ExampleMustConfigureAndParse_basic() {
	// Prepare configuration
	var config struct {
		Greeting string `flag:"greet,Hello,The greeting."`
		Name     string `flag:"name,User,The user's name."`
		Times    int    `flag:"times,1,Number of repeats."`
	}
	MustConfigureAndParse(&config)

	// Start greeting
	for i := 0; i < config.Times; i++ {
		fmt.Printf("%s %s!\n", config.Greeting, config.Name)
	}
}

func ExampleMustConfigureAndParse_sleep() {
	// Prepare configuration
	var config struct {
		Sleep   time.Duration `flag:"s,1s,The amount of time to sleep."`
		Verbose bool          `flag:"v,false,Verbose output."`
	}
	MustConfigureAndParse(&config)

	// Start sleeping.
	if config.Verbose {
		log.Println("Sleeping for " + config.Sleep.String())
	}
	time.Sleep(config.Sleep)
	if config.Verbose {
		log.Println("Done.")
	}
}
