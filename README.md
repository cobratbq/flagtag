flagtag
=======

A flag configurator for Go's [flag package](http://golang.org/pkg/flag/) based on struct field tags.

This package scans a struct's field for 'flag' tags and configures flags to use the struct's fields as target variables for command line arguments.

Compatibility note
------------------

*flagtag* is fully compatible with Go's flag package. It simply uses the facilities offered by the [flag package](http://golang.org/pkg/flag/). It is also possible to use *flagtag* interchangably with the flag package itself. As with the flag package, you have to be sure that flags have not been parsed yet, while still configuring the flags.
