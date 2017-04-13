//
// Copyright (c) 2017 Dell Inc. or its subsidiaries.  All Rights Reserved.
// Dell EMC Confidential/Proprietary Information
//
//

package auth

import (
	"bufio"
	"fmt"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	log "github.com/sirupsen/logrus"
)

// TargetAuth gets vcenter auth
func TargetAuth(target string) (string, string, string, error) {

	// Get Address
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Enter %s endpoint: ", target)
	scanner.Scan()
	endpoint := scanner.Text()

	if err := scanner.Err(); err != nil {
		log.Warnf("Error reading addr: %s", err)
		return "", "", "", err
	}

	// Get Username
	fmt.Printf("Enter %s Username: ", target)
	scanner.Scan()
	userName := scanner.Text()

	if err := scanner.Err(); err != nil {
		log.Warnf("Error reading username: %s", err)
		return "", "", "", err
	}

	// Get Password
	fmt.Printf("Enter %s Password: ", target)
	var password string

	if terminal.IsTerminal(int(syscall.Stdin)) {
		passwordBytes, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Warnf("\nError reading password: %s", err)
			return "", "", "", err
		}
		password = string(passwordBytes)

	} else {
		password = scanner.Text()

		if err := scanner.Err(); err != nil {
			log.Warnf("\nError reading password: %s", err)
			return "", "", "", err
		}
	}

	return endpoint, userName, password, nil

}
