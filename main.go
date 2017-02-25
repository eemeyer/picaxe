package picaxe

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"strings"

	flags "github.com/jessevdk/go-flags"
)

type Options struct {
	ListenAddress string `short:"l" long:"listen" description:"Listen address." value-name:"ADDRESS"`
}

func main() {
	var options Options
	parser := flags.NewParser(&options, flags.HelpFlag|flags.PassDoubleDash)
	if _, err := parser.Parse(); err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
			parser.WriteHelp(os.Stdout)
			return
		}
		return
	}

	server := NewServer(ServerOptions{
		ResourceResolver: HTTPResourceResolver,
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
