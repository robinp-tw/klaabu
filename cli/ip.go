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

	if err != nil || len(os.Args) < 3 {
		log.Printf("Usage: klaabu ip [OPTIONS] CIDR \n\n Subcommands: \n")
		get.PrintDefaults()
		os.Exit(1)
	}

	schema, err := klaabu.LoadSchemaFromKmlFile(*schemaFileName)
	if err != nil {
		log.Fatalln(err)
	}

	rawIpNet := get.Arg(0)
	_, ipnet, err := net.ParseCIDR(rawIpNet)
	if err != nil {
		log.Fatalf("Couldn't parse CIDR. For raw IPv4, stick /32 after? Error: %v", err)
	}

	matches := schema.SearchNet(ipnet)
	log.Println(matches)
}
