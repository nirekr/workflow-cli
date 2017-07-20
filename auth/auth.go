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
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/dellemc-symphony/workflow-cli/models"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// TargetAuth gets auth
func TargetAuth(target string) (string, string, string, error) {

	var endpoint, username, password string
	scanner := bufio.NewScanner(os.Stdin)

	// Check if endpoints are set in file first
	fileEndpoints := ParseEndpointsFile()
	endpoint = fileEndpoints[target].EndpointURL
	username = fileEndpoints[target].Username
	password = fileEndpoints[target].Password

	if endpoint == "" {

		// Get Address
		fmt.Printf("Enter %s endpoint: ", target)

		scanner.Scan()
		if err := scanner.Err(); err != nil {
			log.Warnf("Error reading addr: %s", err)
			return "", "", "", err
		}

		endpoint = scanner.Text()

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

//ParseEndpointsFile parses the file for the endpoints
func ParseEndpointsFile() map[string]models.Endpoint {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	services := []string{"rackhd", "hostbmc", "vcenter", "scaleiogateway"}
	endpoints := make(map[string]models.Endpoint, len(services))

	viper.SetConfigName("endpoint")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(dir)

	err = viper.ReadInConfig()
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			log.Warnf(`Config file "endpoint.yaml" not found.`)
		} else {
			log.Warnf("Invalid endpoint.yaml: %s", err)
		}

		log.Warnf("Will prompt user for endpoints.")
		return endpoints
	}

	for _, service := range services {

		entry := models.Endpoint{}

		endpoint := viper.GetStringSlice(service + ".endpoint")
		if len(endpoint) == 1 {
			entry.EndpointURL = endpoint[0]
		} else {
			entry.EndpointURL = ""
		}

		username := viper.GetStringSlice(service + ".username")
		if len(username) == 1 {
			entry.Username = username[0]
		} else {
			entry.Username = ""
		}

		password := viper.GetStringSlice(service + ".password")
		if len(password) == 1 {
			entry.Password = password[0]
		} else {
			entry.Password = ""
		}

		endpoints[service] = entry

	}

	return endpoints
}
