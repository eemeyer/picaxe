package main

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"

	"github.com/t11e/picaxe/iiif"
	"github.com/t11e/picaxe/server"
)

type Options struct {
	ListenAddress string `short:"l" long:"listen" description:"Listen address." value-name:"[HOST][:PORT]"`
}

func main() {
	var options Options
	parser := flags.NewParser(&options, flags.HelpFlag|flags.PassDoubleDash)
	if _, err := parser.Parse(); err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
			parser.WriteHelp(os.Stdout)
			os.Exit(2)
		}
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	server := server.NewServer(server.ServerOptions{
		ResourceResolver: server.HTTPResourceResolver,
		ProcessorFactory: iiif.DefaultProcessorFactory,
	})
	if err := server.Run(ensureAddressWithPort(options.ListenAddress, 7073)); err != nil {
		log.Fatal(err)
	}
}

func ensureAddressWithPort(address string, defaultPort int) string {
	if address == "" {
		return fmt.Sprintf(":%d", defaultPort)
	} else if !strings.Contains(address, ":") {
		return fmt.Sprintf("%s:%d", address, defaultPort)
	}
	return address
}
