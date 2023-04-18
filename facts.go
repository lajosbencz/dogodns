package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/digitalocean/godo"
)

func initFacts(ctx context.Context) error {
	var err error

	// detect public IP address
	if publicAddress, err = getPublicIP(cfg.PIP, cfg.IP6); err != nil {
		publicAddress = ""
		return err
	}

	// init top and sub domain names
	var res *godo.Response
	domains, res, err = doClient.Domains.List(ctx, &godo.ListOptions{
		Page:    1,
		PerPage: 100,
	})
	if err != nil {
		if res != nil && res.StatusCode >= 500 {
			return ErrDOInternal
		}
		return err
	}
	dotCount := strings.Count(cfg.Domain, ".")
	for _, domain := range domains {
		if strings.Contains(cfg.Domain, domain.Name) {
			topDomain = domain.Name
			if dotCount > 1 {
				subDomain = cfg.Domain[:len(cfg.Domain)-(1+len(topDomain))]
			} else {
				subDomain = "@"
			}
			break
		}
	}

	// check if record exists
	existingRecord = nil
	previousAddress = ""
	records, _, err := doClient.Domains.Records(ctx, topDomain, nil)
	if err != nil {
		return err
	}
	for _, record := range records {
		if record.Name == subDomain {
			existingRecord = &record
			previousAddress = record.Data
		}
	}

	if topDomain == "" {
		return fmt.Errorf("TLD for %s does not exist in DigitalOcean domain list", cfg.Domain)
	}

	return nil
}
