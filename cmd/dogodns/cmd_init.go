package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/lajosbencz/dogodns/pkg/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

var cmdInit = &cobra.Command{
	Use:   "init",
	Short: "Creates sample config file specified by --config",
	Run:   runInit,
}

func init() {
	cmdInit.Flags().Bool("prompt", false, "Prompt for DigitalOcean R+W token")
	cmdInit.Flags().Bool("force", false, "Force overwrite file")
	cmdInit.Flags().Bool("fail", false, "Exit with error if file already exists")
}

func runInit(cmd *cobra.Command, args []string) {
	flagPrompt, _ := cmd.Flags().GetBool("prompt")
	flagForce, _ := cmd.Flags().GetBool("force")
	flagFail, _ := cmd.Flags().GetBool("fail")
	_, err := os.Stat(cfgPath)
	if err != nil && !os.IsNotExist(err) {
		fmt.Printf("failed to stat config file %s: %s\n", cfgPath, err)
		os.Exit(3)
	}
	createFile := func() error {
		cfg := config.DefaultConfig("dogodns.example.tld", "<secret>")
		if flagPrompt {
			fmt.Print("Token: ")
			buf, _ := term.ReadPassword(int(syscall.Stdin))
			fmt.Println()
			cfg.Token = string(buf)
		}
		cfgBytes, _ := yaml.Marshal(cfg)
		return os.WriteFile(cfgPath, cfgBytes, 0664)
	}
	if os.IsNotExist(err) {
		fmt.Printf("creating config file %s\n", cfgPath)
		if err := createFile(); err != nil {
			fmt.Printf("failed to create config file: %s\n", err)
			os.Exit(2)
		}
	} else {
		fmt.Printf("config file exists: %s\n", cfgPath)
		if flagFail {
			os.Exit(1)
		}
		if flagForce {
			fmt.Printf("force creating config file %s\n", cfgPath)
			if err := createFile(); err != nil {
				fmt.Printf("failed to create config file: %s\n", err)
				os.Exit(2)
			}
		}
	}
	os.Exit(0)
}
