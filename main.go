package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/digitalocean/godo"
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
	ticker          *time.Ticker
	paused          bool
)

func commit(ctx context.Context, req *godo.DomainRecordEditRequest) (err error) {
	var res *godo.Response
	if existingRecord != nil {
		_, res, err = doClient.Domains.EditRecord(ctx, topDomain, existingRecord.ID, req)
	} else {
		_, res, err = doClient.Domains.CreateRecord(ctx, topDomain, req)
	}
	if err != nil && res.StatusCode >= 500 {
		err = ErrDOInternal
	}
	return
}

func main() {

	if err := initArgs(); err != nil {
		fail(err)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)

	doClient = godo.NewFromToken(cfg.Token)

	ticker = time.NewTicker(time.Second * time.Duration(cfg.Interval))

	paused = false

	f := func() error {
		ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*10)
		defer ctxCancel()

		if err := initFacts(ctx); err != nil {
			return err
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
					return err
				}
			}

			fmt.Printf("%s %s\n", cfg.Domain, publicAddress)

			previousAddress = publicAddress
		} else {
			fmt.Printf("no change in %s\n", publicAddress)
		}
		return nil
	}

	err := f()
	if err != nil {
		errDO(err)
		fmt.Println(err, ", ignoring")
	}

out:
	for {
		select {
		case <-shutdown:
			break out
		case <-ticker.C:
			if paused {
				continue
			}
			err := f()
			if err != nil {
				errDO(err)
				fmt.Println(err, ", suspending for 10 minutes")
				paused = true
				go func() {
					time.Sleep(time.Minute * 10)
					paused = false
				}()
			}
		}
	}
}
