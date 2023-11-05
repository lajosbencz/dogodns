package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cmdIp = &cobra.Command{
	Use:   "ip",
	Short: "Prints public IP",
	Run:   runIp,
}

func runIp(cmd *cobra.Command, args []string) {
	publicIP, err := getIpResolver().Resolve()
	if err != nil {
		fmt.Printf("failed to resolve public IP: %v\n", err)
	} else {
		fmt.Println(publicIP)
	}
}
