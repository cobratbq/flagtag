flagtag
=======
*... a Go command line flag support package*

*flagtag* is a little package that focuses on providing very quick &amp; easy support for connecting command line arguments to variables. The package provides a **flag** tag and a mechanism for automatically discovering tagged variables for which to set up flags.

By using this mechanism you have to do only very little programming. Variables with either user-provided or default values are readily available in the \(tagged\) variables after calling the flag package to parse command line arguments.

This flag configurator is based on Go's [flag package](http://golang.org/pkg/flag/) and facilitates struct field tags. The function scans the provided variables for tags and automatically defines flags for tagged variables based on the content of the flag tag.

When to use this package?
------------------------
I have created this package specifically for those occasions where you do not want to spend a lot of time defining and fine tuning flags and command line arguments. For example, when creating a rather basic command line tool. The package creates a direct mapping between flags and struct fields. Tags are simply declared using a '*flag*'-tag. A single function call will make all arrangements for your tagged variables in the background.

The tag formats
---------------

The flag:
~~~
flag:"<flag-name>,<default-value>,<usage-description>"
~~~

The flag options:
~~~
flagopt:"<option>"
~~~

Available options:

* **skipFlagValue** - Skip testing for *flag.Value* implementation and immediately continue with primitive types.

A basic example
---------------
A basic example follows. Below the example there will be a small description of what the tags accomplish.

~~~
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
~~~

This basic program defines a type that stores some configuration values. Apart from the default values you may wish to provide support for users to provide alternative configurations via command line arguments. Using this package you accomplish the same by defining an appropriate tag.

3 flags are being defined:

* *-greet* with default value '*Hello*',
* *-name* with the default value '*User*',
* *-times* with the default amount of *1*.

Go's *flag* package also provides a flag *-help* which prints help information. The last part of the tag defines a usage description that will be shown when help information is printed.

Features
--------

* Based on the default behavior of *flag* package.
* Supports *flag*'s primitive types,
* and supports types derived from these primitive types.
* Support for type [*time.Duration*](http://golang.org/pkg/time/#Duration), as this is also supported by *flag*.
* Support for pointers and interfaces to variables. (It does **not** appreciate *nil* though.)
* Any types that implement the [*flag.Value*](http://golang.org/pkg/flag/#Value) interface.
* Recursively configuring nested structs (unless they themselves are tagged).
* Either returning an error or panicking, whatever suits your needs.
* Do a one-pass **configure &amp; parse** and be done with it, or configure multiple structs and/or define your own additional flags yourself. You can define your own flags interchangeably with using the flagtag package.
* Default value for flag.Value implementors. If a default value is provided (i.e. non-zero length) then Set() will be called first with the default value. Then again if the flag was specified as a command line argument.
* Flag options using the tag '*flagopt*'.
  * '*skipFlagValue*' can be provided to indicate that we should skip testing for the implementation of flag.Value. This can be used if the struct field accidentally implements flag.Value but you do not want to use it.

Under consideration
-------------------

* More advanced syntax for '*flagopt*' tag for other options.
  * Define syntax for use w/ and w/o value, such that the flagopt tag can be extended at a later time.
  * Separate multiple options with commas.
  * Support for more advanced values using double quotes or something ...

Compatibility notes
-------------------

*flagtag* is fully compatible with Go's flag package. It simply uses the facilities offered by the [flag package](http://golang.org/pkg/flag/). It is also possible to use *flagtag* interchangeably with the flag package itself. As with the flag package, you have to be sure that flags have not been parsed yet, while still configuring the flags.

Support for default values is not consistent with the flag package. The flag package does not provide a parameter for a default value. *flagtag* simply calls Set(..) first with the default value *if* a default value is provided, then again if the flag is specified as a program argument.
