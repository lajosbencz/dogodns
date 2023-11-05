package main

import (
	"fmt"

	"github.com/lajosbencz/dogodns/pkg/dns"
	"github.com/spf13/cobra"
)

var cmdStatus = &cobra.Command{
	Use:   "status",
	Short: "Shows a brief status page",
	Run:   runStatus,
}

func runStatus(cmd *cobra.Command, args []string) {
	for _, domain := range cfg.GetDomains() {
		fmt.Printf("%s: ", domain)
		r := &dns.Record{Name: domain}
		if !getDnsUpdater().Has(r) {
			fmt.Printf("N/A\n")
		} else {
			err := getDnsUpdater().Get(r)
			if err != nil {
				fmt.Printf(" error\nfailed to fetch dns record for %s: %s\n", r.Name, err)
				return
			}
			fmt.Printf("%s\n", r.Data)
		}
	}
}
