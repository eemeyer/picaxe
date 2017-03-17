package main

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"

	"github.com/t11e/picaxe/iiif"
	"github.com/t11e/picaxe/resources"
	"github.com/t11e/picaxe/server"
)

type Options struct {
	ListenAddress string `short:"l" long:"listen" description:"Listen address." value-name:"[HOST][:PORT]"`
	MaxAge        string `short:"m" long:"max-age" default:"31536000s" description:"max-age for cache-control response header." value-name:"[integer][unit h,m, or s]"`
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
	var maxAge time.Duration
	{
		var err error
		if maxAge, err = time.ParseDuration(options.MaxAge); err != nil {
			fmt.Fprintf(os.Stderr, "max-age %s\n", err.Error())
			os.Exit(1)
		}
	}

	server := server.NewServer(server.ServerOptions{
		ResourceResolver: resources.HTTPResolver,
		Processor:        iiif.DefaultProcessor,
		MaxAge:           maxAge,
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
