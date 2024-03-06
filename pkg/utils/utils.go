package utils

import (
	"flag"
	"fmt"
	"github.com/go-errors/errors"
	"os"
)

func RequireParamAndExit(param string) {
	fmt.Fprintf(os.Stderr, "error: --%s parameter is requried\n", param)
	fmt.Println("Usage:")
	flag.PrintDefaults()
	os.Exit(1)
}

func PrintErrorWithStackAndExit(err error) {
	fmt.Printf("error: %s\n", err.Error())
	fmt.Printf("%s\n", errors.Wrap(err, 0).Stack())
	os.Exit(1)
}
