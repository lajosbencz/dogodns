package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/lajosbencz/dogodns/pkg/dns"
	"github.com/spf13/cobra"
)

var cmdService = &cobra.Command{
	Use:   "service",
	Short: "Run the service",
	Run:   runService,
}

func init() {
	cmdService.Flags().Bool("dry", false, "Dry run")
	cmdService.Flags().BoolP("verbose", "v", false, "Verbose output")
}

func runService(cmd *cobra.Command, args []string) {
	dryRun, _ := cmd.Flags().GetBool("dry")
	verbose, _ := cmd.Flags().GetBool("verbose")
	commit := func(record *dns.Record) error {
		if dryRun {
			fmt.Printf("dry run, skipping %#v\n", record)
			return nil
		}
		return getDnsUpdater().Set(record)
	}
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	ticker := time.NewTicker(time.Second * time.Duration(cfg.Interval))
	tick := func() {
		ipAddress, err := getIpResolver().Resolve()
		if err != nil {
			log.Fatal(err)
		}
		if verbose {
			log.Printf("public IP address is %s\n", ipAddress)
		}
		for _, domainName := range cfg.GetDomains() {
			record := dns.NewRecord(domainName, ipAddress)
			if !getDnsUpdater().Has(record) {
				if verbose {
					log.Printf("%s does not exist, creating\n", domainName)
				}
				log.Printf("%s (%s) created\n", record.Name, record.Data)
				if err := commit(record); err != nil {
					log.Fatalln(err)
				}
			} else {
				if verbose {
					log.Printf("%s exists, updating\n", domainName)
				}
				if err := getDnsUpdater().Pull(record); err != nil {
					log.Fatalln(err)
				}
				if record.Data != ipAddress {
					record.Data = ipAddress
					log.Printf("%s (%s) updated\n", record.Name, record.Data)
					if err := commit(record); err != nil {
						log.Fatalln(err)
					}
				} else {
					if verbose {
						log.Printf("%s is up to date\n", domainName)
					}
				}
			}
		}
		if verbose {
			log.Printf("will check again in about %ds\n", cfg.Interval)
		}
	}
	tick()
	if dryRun {
		fmt.Println("dry run, exiting after first tick")
		os.Exit(0)
	}
out:
	for {
		select {
		case <-shutdown:
			break out
		case <-ticker.C:
			tick()
		}
	}
}
