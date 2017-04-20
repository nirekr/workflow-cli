//
// Copyright (c) 2017 Dell Inc. or its subsidiaries.  All Rights Reserved.
// Dell EMC Confidential/Proprietary Information
//
//

package resources

import (
	"fmt"
	"os"
)

// WriteEndpointsFile writes out the endpoints file
func WriteEndpointsFile(configuration, fileLocation string) error {
	_, err := os.Stat(fileLocation)
	if err == nil {
		return fmt.Errorf("File already exists: %s", fileLocation)
	}

	f, err := os.Create(fileLocation)
	if err != nil {
		return fmt.Errorf("Error opening file: %s", err.Error())
	}

	switch configuration {

	case "AllFields":
		f.WriteString(allFields)

	case "MissingEndpoint":
		f.WriteString(missingEndpoints)

	case "MissingCredentials":
		f.WriteString(missingCredentials)

	default:
		err = fmt.Errorf("unknown configuration: %s", configuration)
	}

	return err

}

const allFields = `
# Endpoints for FRU-PAQX workflow
---
MinimumVersion: v0.0.1-32

RackHD:
  endpoint:
    - "https://10.10.10.10:9090"
  username:
    - "admin"
  password:
    - "admin"

CoprHD:
  endpoint:
    - "https://10.10.10.10:9090"
  username:
    - "admin"
  password:
    - "admin"

vCenter:
  endpoint:
    - "https://10.10.10.10:9090"
  username:
    - "admin"
  password:
    - "admin"

ScaleIO:
  endpoint:
    - "https://10.10.10.10:9090"
  username:
    - "admin"
  password:
    - "admin"`

const missingEndpoints = `
# Endpoints for FRU-PAQX workflow
---
MinimumVersion: v0.0.1-32
RackHD:
  endpoint:
#    - "https://10.10.10.10:9090"
  username:
    - "admin"
  password:
    - "admin"
CoprHD:
  endpoint:
#    - "https://10.10.10.10:9090"
  username:
    - "admin"
  password:
    - "admin"
vCenter:
  endpoint:
#    - "https://10.10.10.10:9090"
  username:
    - "admin"
  password:
    - "admin"
ScaleIO:
  endpoint:
#    - "https://10.10.10.10:9090"
  username:
    - "admin"
  password:
    - "admin"`

const missingCredentials = `
# Endpoints for FRU-PAQX workflow
---
MinimumVersion: v0.0.1-32
RackHD:
  endpoint:
    - "https://10.10.10.10:9090"
  username:
#    - "admin"
  password:
#    - "admin"
CoprHD:
  endpoint:
    - "https://10.10.10.10:9090"
  username:
#    - "admin"
  password:
#    - "admin"
vCenter:
  endpoint:
    - "https://10.10.10.10:9090"
  username:
#    - "admin"
  password:
#    - "admin"
ScaleIO:
  endpoint:
    - "https://10.10.10.10:9090"
  username:
#    - "admin"
  password:
#    - "admin"`
