package main

import (
	"flag"
	"fmt"
	"os"
)

// Defaults.
const (
	namespaceDef             = ""
	resyncIntervalSecondsDef = 30
	wildcardsDef             = false
	dryRunDef                = false
	developmentDef           = false
	debugDef                 = false
)

type stringArray []string

func (s *stringArray) String() string {
	return fmt.Sprintf("%s", *s)
}

func (s *stringArray) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// Flags are the flags of the program.
type Flags struct {
	ServiceName   string
	listenAddress string
	Endpoints     stringArray
}

// NewFlags returns the flags of the commandline.
func NewFlags() *Flags {
	flags := &Flags{}
	fl := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	fl.StringVar(&flags.ServiceName, "service-name", "tracing-example", "the name of the service")
	fl.StringVar(&flags.listenAddress, "listen-address", ":8080", "address where the server will be listening")
	fl.Var(&flags.Endpoints, "endpoint", "endpoint list (flag can be repeated)")

	fl.Parse(os.Args[1:])

	return flags
}
