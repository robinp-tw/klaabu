package main

import (
	"flag"
	"github.com/erikkn/klaabu/klaabu"
	"log"
	"net"
	"os"
)

func ipCommand() {
	get := flag.NewFlagSet("ip", flag.ExitOnError)
	schemaFileName := get.String("schema", "schema.kml", "(Optional) Schema file path (defaults to ./schema.kml)")
	err := get.Parse(os.Args[2:])

	rawIp := get.Arg(0)

	if err != nil || len(os.Args) < 3 {
		log.Printf("Usage: klaabu ip [OPTIONS] IP_ADDRESS \n\n Subcommands: \n")
		get.PrintDefaults()
		os.Exit(1)
	}

	schema, err := klaabu.LoadSchemaFromKmlFile(*schemaFileName)
	if err != nil {
		log.Fatalln(err)
	}

	ip := net.ParseIP(rawIp)
	if ip == nil {
		log.Fatalf("Couldn't parse ip.")
	}

	matches := schema.SearchIp(ip)
	log.Println(matches)
}
