//
// Copyright (c) 2017 Dell Inc. or its subsidiaries.  All Rights Reserved.
// Dell EMC Confidential/Proprietary Information
//
//

package main

import (
	"flag"
	"strconv"

	"github.com/dellemc-symphony/workflow-cli/mock"
	log "github.com/sirupsen/logrus"
)

var https = flag.Bool("https", false, "Set 'true' to enable HTTPS for mock REST endpoint")
var port  = flag.Int("port", 8080, "Port to start up the mock REST endpoint")

func init() {
	flag.Parse()
}

func main() {

	scheme := "http://"
	if *https {
		scheme = "https://"
	}

	bindPort := strconv.Itoa(*port)

	log.Infof("Starting mock REST endpoint at " + scheme + "localhost:" + bindPort)
	log.Infof("HTTPS: %v", *https)

	mock.CreateMock(*https, *port)
	defer mock.StopMock()

	select {}
}
