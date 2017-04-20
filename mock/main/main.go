//
// Copyright (c) 2017 Dell Inc. or its subsidiaries.  All Rights Reserved.
// Dell EMC Confidential/Proprietary Information
//
//

package main

import (
	"flag"

	"github.com/dellemc-symphony/workflow-cli/mock"
	log "github.com/sirupsen/logrus"
)

var https = flag.Bool("https", false, "Set 'true' to enable HTTPS for mock REST endpoint")

func init() {
	flag.Parse()
}

func main() {

	scheme := "http://"
	if *https {
		scheme = "https://"
	}

	log.Infof("Starting mock REST endpoint at " + scheme + "localhost:8080")
	log.Infof("HTTPS: %v", *https)

	mock.CreateMock(*https)
	defer mock.StopMock()

	select {}
}
