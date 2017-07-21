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

	"github.com/dellemc-symphony/workflow-cli/resources"
	log "github.com/sirupsen/logrus"
)

// TargetAuth gets auth
func TargetAuth(target string) (string, string, string, error) {

	var endpoint, username, password string
	scanner := bufio.NewScanner(os.Stdin)

	// Check if endpoints are set in file first
	fileEndpoints := resources.ParseEndpointsFile(resources.UseEndpointFile)
	endpoint = fileEndpoints[target].EndpointURL
	username = fileEndpoints[target].Username
	password = fileEndpoints[target].Password

	for endpoint == "" {

		// Get Address
		fmt.Printf("Enter %s endpoint: ", target)

		scanner.Scan()
		if err := scanner.Err(); err != nil {
			log.Warnf("Error reading addr: %s", err)
			return "", "", "", err
		}

		endpoint = scanner.Text()
		if !(resources.ValidateEndpoint(endpoint)) {
			endpoint = ""
		}
	}

	if username == "" {
		// Get Username
		fmt.Printf("Enter %s Username: ", target)

		scanner.Scan()
		if err := scanner.Err(); err != nil {
			log.Warnf("Error reading username: %s", err)
			return "", "", "", err
		}

		username = scanner.Text()

	}

	if password == "" {
		confirmPassword := "default"

		for confirmPassword != password {
			// Get Password
			fmt.Printf("Enter %s Password: ", target)

			if terminal.IsTerminal(int(syscall.Stdin)) {
				passwordBytes, err := terminal.ReadPassword(int(syscall.Stdin))
				if err != nil {
					log.Warnf("\nError reading password: %s", err)
					return "", "", "", err
				}
				password = string(passwordBytes)

			} else {
				scanner.Scan()
				if err := scanner.Err(); err != nil {
					log.Warnf("\nError reading password: %s", err)
					return "", "", "", err
				}

				password = scanner.Text()
			}
			fmt.Printf("\n")

			fmt.Printf("Confirm %s Password: ", target)
			if terminal.IsTerminal(int(syscall.Stdin)) {
				passwordBytes, err := terminal.ReadPassword(int(syscall.Stdin))
				if err != nil {
					log.Warnf("\nError reading password: %s", err)
					return "", "", "", err
				}
				confirmPassword = string(passwordBytes)

			} else {
				scanner.Scan()
				if err := scanner.Err(); err != nil {
					log.Warnf("\nError reading password: %s", err)
					return "", "", "", err
				}

				confirmPassword = scanner.Text()
			}
			fmt.Printf("\n")

			if password != confirmPassword {
				log.Warnf("Passwords for %s don't match.\n", target)
			}
		}
	}

	return endpoint, username, password, nil

}
