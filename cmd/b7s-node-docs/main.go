package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/spf13/pflag"

	"github.com/blocklessnetwork/b7s/internal/config"
)

//go:embed assets/*
var assets embed.FS

func main() {

	var (
		flagAddress string
		flagOutput  string
		flagEmbed   bool
	)
	pflag.StringVarP(&flagAddress, "address", "a", "127.0.0.1:8080", "address to serve on")
	pflag.StringVarP(&flagOutput, "output", "o", "", "output file to write the documentation to")
	pflag.BoolVarP(&flagEmbed, "embed", "e", true, "use embedded files for assets")
	pflag.Parse()

	configs := config.GetConfigDocumentation()
	component := page(configs)

	if flagOutput != "" {

		f, err := os.Create(flagOutput)
		if err != nil {
			log.Fatalf("could not open file: %s", err)
		}

		err = component.Render(context.Background(), f)
		if err != nil {
			log.Fatalf("could not render component: %s", err)
		}

		f.Close()
		return
	}

	mux := http.NewServeMux()

	var fh http.Handler
	if flagEmbed {
		fh = http.FileServer(http.FS(assets))
	} else {
		fh = http.StripPrefix("/assets/", http.FileServer(http.Dir("assets")))
	}

	mux.Handle("/assets/", fh)
	mux.Handle("/", templ.Handler(component))

	fmt.Printf("Documentation served on http://%s/", flagAddress)

	err := http.ListenAndServe(flagAddress, mux)
	if err != nil {
		log.Fatalf("failed to start server: %s", err)
	}
}
