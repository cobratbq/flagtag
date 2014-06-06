package main

import (
	"fmt"
	"github.com/cobratbq/flagtag"
)

type Configuration struct {
	Greeting string `flag:"greet,Hello,The greeting."`
	Name     string `flag:"name,User,The user's name."`
	Times    int    `flag:"times,1,Number of repeats."`
}

func main() {
	var config Configuration
	flagtag.MustConfigureAndParse(&config)

	for i := 0; i < config.Times; i++ {
		fmt.Printf("%s %s!\n", config.Greeting, config.Name)
	}
}
