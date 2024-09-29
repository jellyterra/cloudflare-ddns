// Copyright 2024 Jelly Terra
// Use of this source code form is governed under the MIT license.

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/jellyterra/cloudflare-ddns/config"
	"github.com/jellyterra/cloudflare-ddns/ddns"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"sync/atomic"
	"syscall"
	"time"
)

func main() {
	err := _main()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func _main() error {
	var (
		optConfig = flag.String("c", "config.yaml", "Path to config TOML file")

		configFile config.File
	)

	flag.Parse()

	b, err := os.ReadFile(*optConfig)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, &configFile)
	if err != nil {
		return err
	}

	env, err := ddns.LoadConfig(&configFile)
	if err != nil {
		return err
	}

	var (
		haveToUpdate atomic.Bool
		notification = make(chan struct{}, 1)
	)

	go func() {
		for {
			// To keep the result latest, it must be set before the update.
			haveToUpdate.Store(false)

			err := env.UpdateAllZones(context.Background())
			if err != nil {
				log.Println("Error updating DNS records:", err)
			}

			log.Println("All zones updated.")

			for !haveToUpdate.Load() {
				time.Sleep(1 * time.Second)
			}
		}
	}()

	go func() {
		for {
			select {
			case <-notification:
				log.Println("Network change detected.")
				haveToUpdate.Store(true)
			}
		}
	}()

	return notify(notification)
}

func notify(c chan<- struct{}) error {

	fd, err := syscall.Socket(syscall.AF_NETLINK, syscall.SOCK_DGRAM, syscall.NETLINK_ROUTE)
	if err != nil {
		return err
	}

	err = syscall.Bind(fd, &syscall.SockaddrNetlink{
		Family: syscall.AF_NETLINK,
		Pid:    0,
		Groups: (1 << (syscall.RTNLGRP_LINK - 1)) | (1 << (syscall.RTNLGRP_IPV4_IFADDR - 1)) | (1 << (syscall.RTNLGRP_IPV6_IFADDR - 1)),
	})
	if err != nil {
		return err
	}

	for {
		packet := make([]byte, 2048)

		n, err := syscall.Read(fd, packet)
		if err != nil {
			return err
		}

		messages, err := syscall.ParseNetlinkMessage(packet[:n])
		if err != nil {
			return err
		}

		for _, message := range messages {
			if message.Header.Type == syscall.RTM_NEWADDR || message.Header.Type == syscall.RTM_DELADDR || message.Header.Type == syscall.RTM_GETADDR {
				c <- struct{}{}
				break
			}
		}
	}
}
