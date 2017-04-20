# Workflow-cli Mock Server
## Usage
Start this mock server to test the workflow-cli against.
Run it with `go run main.go`

By default, HTTPS is disabled. It can be explicitly set with the `--https` flag, like so
`go run main.go --https=true`
or
`go run main.go --https=false`

**WARNING**
The HTTPS certs included in this mock server were generated for testing ONLY by a non-secure script. Do NOT use in production or for securing real traffic.

## Copyright
Copyright (c) 2017 Dell Inc. or its subsidiaries.  All Rights Reserved.
Dell EMC Confidential/Proprietary Information
