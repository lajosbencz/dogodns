package main

import (
	"fmt"
	"os"

	"github.com/lajosbencz/dogodns/pkg/config"
	"github.com/lajosbencz/dogodns/pkg/dns"
	"github.com/lajosbencz/dogodns/pkg/ip"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgPath string

var cfg config.Config

var ipResolver ip.Resolver

var dnsUpdater dns.Updater

var cmdRoot = &cobra.Command{
	Use:   "dogodns",
	Short: "DigitalOcean Dynamic IP Client",
	Long: `This app will use the public IP of it's network
to update the DNS registry of DigitalOcean`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmdRoot.PersistentFlags().StringVarP(&cfgPath, "config", "c", config.DefaultPath, "Path to config file without extension")
	cmdRoot.PersistentFlags().StringVarP(&cfg.Token, "token", "t", "", "DigitalOcean R+W token")
	cmdRoot.PersistentFlags().StringArrayVarP(&cfg.Domains, "domain", "d", []string{}, "List of domain names")
	cmdRoot.PersistentFlags().StringVarP(&cfg.PIP, "pip", "p", config.DefaultPIP, "HTTP URL to fetch public IP from")
	cmdRoot.PersistentFlags().IntVarP(&cfg.Interval, "interval", "i", config.DefaultInterval, "Interval between checks")
	cmdRoot.PersistentFlags().IntVarP(&cfg.TTL, "ttl", "l", config.DefaultTTL, "Domain record TTL")
	viper.BindPFlags(cmdRoot.Flags())
	cobra.OnInitialize(initConfig)
	cmdRoot.AddCommand(cmdInit)
	cmdRoot.AddCommand(cmdIp)
	cmdRoot.AddCommand(cmdStatus)
	cmdRoot.AddCommand(cmdService)
}

func initConfig() {
	var err error
	_, err = os.Stat(cfgPath)
	if err == nil {
		viper.SetConfigFile(cfgPath)
		if err = viper.ReadInConfig(); err != nil {
			fmt.Printf("error reading config file: %s\n", err)
		}
	}
	if err = viper.Unmarshal(&cfg); err != nil {
		fmt.Printf("error unmarshaling config file: %s\n", err)
	}
}

func getIpResolver() ip.Resolver {
	var err error
	if ipResolver == nil {
		if ipResolver, err = ip.DefaultResolver(cfg); err != nil {
			fmt.Printf("failed to create IP Resolver: %s\n", err)
		}
	}
	return ipResolver
}

func getDnsUpdater() dns.Updater {
	var err error
	if dnsUpdater == nil {
		if dnsUpdater, err = dns.DefaultUpdater(cfg); err != nil {
			fmt.Printf("failed to create DNS Updater: %s\n", err)
		}
	}
	return dnsUpdater
}
