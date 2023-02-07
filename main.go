package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/digitalocean/godo"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	ENV_PREFIX = "DOGODNS_"

	defPublicIP = "https://api.ipify.org/?format=raw"
	defInterval = 60
	defTTL      = 300
)

type config struct {
	Domain   string
	Token    string
	IP6      bool
	TTL      int
	PIP      string
	Interval int
	Dry      bool
}

// globals
var (
	cfg             config
	doClient        *godo.Client
	domains         []godo.Domain
	existingRecord  *godo.DomainRecord
	topDomain       string
	subDomain       string
	publicAddress   string
	previousAddress string
)

func initArgs() error {
	viper.RegisterAlias("d", "domain")
	viper.RegisterAlias("t", "token")
	viper.RegisterAlias("i", "interval")
	flag.String("pip", defPublicIP, "Public IP fetch URL")
	flag.String("token", "", "DigitalOcean API R+W token")
	flag.String("domain", "", "Domain name")
	flag.Int("interval", defInterval, "Check interval")
	flag.Int("ttl", defTTL, "Record TTL")
	flag.Bool("dry", false, "Dry run, don't commit any changes")
	viper.SetConfigName("dogodns")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/dogodns/")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(err)
		}
	}
	viper.SetEnvPrefix("dogodns")
	viper.AutomaticEnv()
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	if !viper.IsSet("token") {
		return errors.New("missing token")
	}
	if !viper.IsSet("domain") {
		return errors.New("missing domain")
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		fail(err)
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		if err := viper.Unmarshal(&cfg); err != nil {
			fail(err)
		}
		fmt.Printf("config file changed: %s\n", e.Name)
	})
	viper.WatchConfig()

	return nil
}

func initFacts(ctx context.Context) error {

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		// detect public IP address
		defer wg.Done()
		var err error
		publicAddress, err = getPublicIP(cfg.PIP, cfg.IP6)
		if err != nil {
			fail(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// client for DigitalOcean
		doClient = godo.NewFromToken(cfg.Token)

		// init top and sub domain names
		var err error
		domains, _, err = doClient.Domains.List(ctx, &godo.ListOptions{
			Page:    1,
			PerPage: 100,
		})
		if err != nil {
			fail(err)
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
		records, _, err := doClient.Domains.Records(ctx, topDomain, nil)
		if err != nil {
			fail(err)
		}
		for _, record := range records {
			if record.Name == subDomain {
				existingRecord = &record
				previousAddress = record.Data
			}
		}
	}()

	wg.Wait()

	if topDomain == "" {
		return fmt.Errorf("TLD for %s does not exist in DigitalOcean domain list", cfg.Domain)
	}

	return nil
}

func commit(ctx context.Context, req *godo.DomainRecordEditRequest) (err error) {
	if existingRecord != nil {
		_, _, err = doClient.Domains.EditRecord(ctx, topDomain, existingRecord.ID, req)
	} else {
		_, _, err = doClient.Domains.CreateRecord(ctx, topDomain, req)
	}
	return
}

func main() {

	if err := initArgs(); err != nil {
		fail(err)
	}

	for {
		func() {
			ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*10)
			defer ctxCancel()

			if err := initFacts(ctx); err != nil {
				fail(err)
			}

			if previousAddress != publicAddress {
				req := &godo.DomainRecordEditRequest{
					Type: "A",
					Name: subDomain,
					Data: publicAddress,
					TTL:  cfg.TTL,
				}

				if !cfg.Dry {
					if err := commit(ctx, req); err != nil {
						fail(err)
					}
				}

				fmt.Printf("%s %s\n", cfg.Domain, publicAddress)

				previousAddress = publicAddress
			} else {
				fmt.Printf("no change in %s\n", publicAddress)
			}
		}()
		time.Sleep(time.Second * time.Duration(cfg.Interval))
	}
}
