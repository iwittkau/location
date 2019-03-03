package main

import (
	"flag"
	"log"

	"github.com/iwittkau/location/bbolt"

	"github.com/iwittkau/location/http"
)

var (
	hash string
)

func main() {
	opts := http.DefaultOpts
	flag.StringVar(&opts.Address, "address", opts.Address, "set listening address")
	flag.StringVar(&opts.StaticPath, "static", opts.StaticPath, "set path to static files")
	flag.StringVar(&opts.TemplatesPath, "templates", opts.TemplatesPath, "set path to template files")
	flag.StringVar(&hash, "secret", "", "set application secret as bcrypt hash")
	flag.BoolVar(&opts.Debug, "debug", opts.Debug, "enable debug mode")
	flag.Parse()
	s := bbolt.New()
	if err := s.Open("data"); err != nil {
		panic(err)
	}
	defer s.Close()

	opts.SecretHash = hash
	a, err := http.New(opts, s)
	if err != nil {
		panic(err)
	}
	log.Fatal(a.Open())
}
