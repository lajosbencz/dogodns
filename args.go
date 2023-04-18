package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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
		return err
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		if err := viper.Unmarshal(&cfg); err != nil {
			fail(err)
		}
		fmt.Printf("config file changed: %s\n", e.Name)
		doClient = godo.NewFromToken(cfg.Token)
	})
	viper.WatchConfig()

	return nil
}
