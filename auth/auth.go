//
// Copyright © 2017 Dell Inc. or its subsidiaries. All Rights Reserved.
// VCE Confidential/Proprietary Information
//
//

package auth

import (
	"bufio"
	"fmt"
	"os"

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
	fmt.Print("\033[8m") // Hide input
	scanner.Scan()
	password := scanner.Text()
	fmt.Print("\033[28m") // Show input

	if err := scanner.Err(); err != nil {
		log.Warnf("Error reading password: %s", err)
		return "", "", "", err
	}

	return endpoint, userName, password, nil
}
