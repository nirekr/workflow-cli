//
// Copyright (c) 2017 Dell Inc. or its subsidiaries.  All Rights Reserved.
// Dell EMC Confidential/Proprietary Information
//
//

package mock

import (
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/braintree/manners"
	"github.com/dellemc-symphony/workflow-cli/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// StopMock stops the Server
func StopMock() {
	manners.Close()
}

// CreateMock starts a mock REST Endpoint for the FRU workflow
func CreateMock(https bool) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	id := "123abc-456def-789ghi"
	retry := true
	longRunningCount := 0

	// About corresponds to the "status" command
	router.GET("/fru/api/about", func(c *gin.Context) {

		c.String(http.StatusOK, "up and running")

	})

	//Data returns info from Data Collection
	router.GET("/fru/api/data/:taskid", func(c *gin.Context) {
		taskid := c.Param("taskid")

		response := models.MockDataResponse{
			Response: "YooHoo!",
			TaskID:   taskid,
		}

		c.JSON(http.StatusOK, response)

	})

	// GET on /workflow returns list of partially completed workflows (for Resume command)
	router.GET("/fru/api/workflow", func(c *gin.Context) {
		var url string
		if https {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s", "https://", id)
		} else {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s", "http://", id)
		}

		workflow := models.Workflow{
			URI: url,
		}

		workflows := models.Workflows{workflow}

		c.JSON(http.StatusOK, workflows)

	})

	// This resumes a partially completes workflow. trackingID comes from GET /workflow
	router.GET("/fru/api/workflow/:trackingid", func(c *gin.Context) {
		var url string
		if https {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/capture-vcenter-endpoint", "https://", id)
		} else {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/capture-vcenter-endpoint", "http://", id)
		}

		stepNext := models.Link{
			Rel:    "step-next",
			Href:   url,
			Type:   "application/vnd.dellemc.vcenter.endpoint+json",
			Method: "POST",
		}

		links := models.Links{stepNext}

		response := models.Response{
			ID:          id,
			Workflow:    "quanta-replacement-d51b-esxi",
			CurrentStep: "capture-vcenter-endpoint",
			Links:       links,
		}

		c.JSON(http.StatusOK, response)

	})

	steps := make(map[string]string)
	steps[""] = "capture-rackhd-endpoint"
	steps["capture-rackhd-endpoint"] = "capture-hostbmc-endpoint"
	steps["capture-hostbmc-endpoint"] = "capture-vcenter-endpoint"
	steps["capture-vcenter-endpoint"] = "capture-scaleio-endpoint"
	steps["capture-scaleio-endpoint"] = "start-scaleio-data-collection"
	steps["start-scaleio-data-collection"] = "start-vcenter-data-collection"
	steps["start-vcenter-data-collection"] = "present-system-list-remove"
	steps["present-system-list-remove"] = "start-scaleio-remove-workflow"
	steps["start-scaleio-remove-workflow"] = "longrunning/scaleio-remove-workflow"
	steps["longrunning/scaleio-remove-workflow"] = "power-off-scaleio-vm"
	steps["power-off-scaleio-vm"] = "enter-maintanence-mode"
	steps["enter-maintanence-mode"] = "instruct-physical-removal"
	steps["instruct-physical-removal"] = "longrunning/wait-for-rackhd-discovery"
	steps["longrunning/wait-for-rackhd-discovery"] = "present-system-list-add"
	steps["present-system-list-add"] = "configure-disks-rackhd"
	steps["configure-disks-rackhd"] = "install-esxi"
	steps["install-esxi"] = "add-host-to-vcenter"
	steps["add-host-to-vcenter"] = "install-scaleio-vib"
	steps["install-scaleio-vib"] = "exit-vcenter-maintenance-mode"
	steps["exit-vcenter-maintenance-mode"] = "deploy-svm"
	steps["deploy-svm"] = "wait-for-svm-deploy"
	steps["wait-for-svm-deploy"] = "start-scaleio-add-workflow"
	steps["start-scaleio-add-workflow"] = "wait-for-scaleio-add-complete"
	steps["wait-for-scaleio-add-complete"] = "map-scaleio-volumes-to-host"
	steps["map-scaleio-volumes-to-host"] = "NONE"

	// Initiates the workflow
	router.POST("/fru/api/workflow", func(c *gin.Context) {
		var url string
		nextStep := steps[""]
		if https {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "https://", id, nextStep)
		} else {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "http://", id, nextStep)
		}

		stepNext := models.Link{
			Rel:    "step-next",
			Href:   url,
			Type:   "application/vnd.dellemc.rackhd.endpoint+json",
			Method: "POST",
		}

		links := models.Links{stepNext}

		response := models.Response{
			ID:          id,
			Workflow:    "quanta-replacement-d51b-esxi",
			CurrentStep: "Initiate Workflow",
			Links:       links,
		}

		c.JSON(http.StatusCreated, response)
	})

	router.POST("/fru/api/workflow/:trackingid/capture-rackhd-endpoint", func(c *gin.Context) {
		id = c.Param("trackingid")
		nextStep := steps["capture-rackhd-endpoint"]
		// Validate JSON Body
		var rackhdCreds models.Endpoint
		if c.BindJSON(&rackhdCreds) == nil {
			var url string
			if https {
				url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "https://", id, nextStep)
			} else {
				url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "http://", id, nextStep)
			}

			stepNext := models.Link{
				Rel:    "step-next",
				Href:   url,
				Type:   "application/vnd.dellemc.hostbmc.endpoint+json",
				Method: "POST",
			}

			links := models.Links{stepNext}

			response := models.Response{
				ID:          id,
				Workflow:    "quanta-replacement-d51b-esxi",
				CurrentStep: "capture-rackhd-endpoint",
				Links:       links,
			}

			c.JSON(http.StatusCreated, response)
		}
	})

	router.POST("/fru/api/workflow/:trackingid/capture-hostbmc-endpoint", func(c *gin.Context) {
		id = c.Param("trackingid")
		nextStep := steps["capture-hostbmc-endpoint"]
		// Validate JSON Body
		var hostBMCCreds models.Endpoint
		if c.BindJSON(&hostBMCCreds) == nil {
			var url string
			if https {
				url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "https://", id, nextStep)
			} else {
				url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "http://", id, nextStep)
			}

			stepNext := models.Link{
				Rel:    "step-next",
				Href:   url,
				Type:   "application/vnd.dellemc.vcenter.endpoint+json",
				Method: "POST",
			}

			links := models.Links{stepNext}

			response := models.Response{
				ID:          id,
				Workflow:    "quanta-replacement-d51b-esxi",
				CurrentStep: "capture-hostbmc-endpoint",
				Links:       links,
			}

			c.JSON(http.StatusCreated, response)
		}
	})

	router.POST("/fru/api/workflow/:trackingid/capture-vcenter-endpoint", func(c *gin.Context) {
		id = c.Param("trackingid")
		nextStep := steps["capture-vcenter-endpoint"]
		// Validate JSON Body
		var vcenterCreds models.Endpoint
		if c.BindJSON(&vcenterCreds) == nil {
			var url string
			if https {
				url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "https://", id, nextStep)
			} else {
				url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "http://", id, nextStep)
			}

			stepNext := models.Link{
				Rel:    "step-next",
				Href:   url,
				Type:   "application/vnd.dellemc.scaleiogateway.endpoint+json",
				Method: "POST",
			}

			links := models.Links{stepNext}

			response := models.Response{
				ID:          id,
				Workflow:    "quanta-replacement-d51b-esxi",
				CurrentStep: "capture-vcenter-endpoint",
				Links:       links,
			}

			c.JSON(http.StatusCreated, response)
		}
	})

	router.POST("/fru/api/workflow/:trackingid/capture-scaleio-endpoint", func(c *gin.Context) {
		id = c.Param("trackingid")
		nextStep := steps["capture-scaleio-endpoint"]
		// Validate JSON Body
		var scaleioCreds models.Endpoint
		if c.BindJSON(&scaleioCreds) == nil {
			var url string
			if https {
				url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "https://", id, nextStep)
			} else {
				url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "http://", id, nextStep)
			}

			stepNext := models.Link{
				Rel:    "step-next",
				Href:   url,
				Type:   "application/json",
				Method: "POST",
			}

			links := models.Links{stepNext}

			response := models.Response{
				ID:          id,
				Workflow:    "quanta-replacement-d51b-esxi",
				CurrentStep: "capture-scaleio-endpoint",
				Links:       links,
			}

			c.JSON(http.StatusCreated, response)
		}
	})

	router.POST("/fru/api/workflow/:trackingid/start-scaleio-data-collection", func(c *gin.Context) {
		var stepNext models.Link
		nextStep := steps["start-scaleio-data-collection"]

		// intentionally fail once to test retry
		if retry == true {
			retry = false
			id = c.Param("trackingid")
			var url string
			if https {
				url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "https://", id, nextStep)
			} else {
				url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "http://", id, nextStep)
			}

			stepNext = models.Link{
				Rel:    "step-retry",
				Href:   url,
				Type:   "application/json",
				Method: "POST",
			}

		} else {
			id = c.Param("trackingid")
			var url string
			if https {
				url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "https://", id, nextStep)
			} else {
				url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "http://", id, nextStep)
			}

			stepNext = models.Link{
				Rel:    "step-next",
				Href:   url,
				Type:   "application/json",
				Method: "POST",
			}
		}

		links := models.Links{stepNext}

		response := models.Response{
			ID:          id,
			Workflow:    "quanta-replacement-d51b-esxi",
			CurrentStep: "start-scaleio-data-collection",
			Links:       links,
		}

		c.JSON(http.StatusCreated, response)
	})

	router.POST("/fru/api/workflow/:trackingid/start-vcenter-data-collection", func(c *gin.Context) {
		id = c.Param("trackingid")
		nextStep := steps["start-vcenter-data-collection"]
		// Validate JSON Body

		var url string
		if https {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "https://", id, nextStep)
		} else {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "http://", id, nextStep)
		}

		stepNext := models.Link{
			Rel:    "step-next",
			Href:   url,
			Type:   "application/vnd.dellemc.nodes.list.remove+json",
			Method: "POST",
		}

		links := models.Links{stepNext}

		node1 := models.Node{
			Hostname:        "node01",
			ServiceTag:      "1234567",
			ManagementIP:    "10.10.10.1",
			PowerStatus:     "poweredOn",
			ConnectionState: "Connected",
			UUID:            "213894yu924h6ao",
		}

		node2 := models.Node{
			Hostname:        "node02",
			ServiceTag:      "98765432",
			ManagementIP:    "10.10.10.2",
			PowerStatus:     "poweredOn",
			ConnectionState: "Connected",
			UUID:            "654edcvbhnjki87",
		}
		node3 := models.Node{
			Hostname:        "node03",
			ServiceTag:      "91827465",
			ManagementIP:    "10.10.10.3",
			PowerStatus:     "poweredOn",
			ConnectionState: "Connected",
			UUID:            "fdr56765redfghj",
		}
		nodes := models.Nodes{node1, node2, node3}

		response := models.Response{
			ID:          id,
			Workflow:    "quanta-replacement-d51b-esxi",
			CurrentStep: "start-vcenter-data-collection",
			Links:       links,
			Nodes:       nodes,
		}

		c.JSON(http.StatusCreated, response)
	})

	router.POST("/fru/api/workflow/:trackingid/present-system-list-remove", func(c *gin.Context) {
		id = c.Param("trackingid")
		nextStep := steps["present-system-list-remove"]

		var nodeToRemove models.Node
		if c.BindJSON(&nodeToRemove) == nil {

			var url string
			if https {
				url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "https://", id, nextStep)
			} else {
				url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "http://", id, nextStep)
			}

			stepNext := models.Link{
				Rel:    "step-next",
				Href:   url,
				Type:   "application/json",
				Method: "POST",
			}

			links := models.Links{stepNext}

			response := models.Response{
				ID:          id,
				Workflow:    "quanta-replacement-d51b-esxi",
				CurrentStep: "present-system-list-remove",
				Links:       links,
			}

			c.JSON(http.StatusCreated, response)
		}
	})

	router.POST("/fru/api/workflow/:trackingid/start-scaleio-remove-workflow", func(c *gin.Context) {
		id = c.Param("trackingid")
		nextStep := steps["start-scaleio-remove-workflow"]
		var url string
		if https {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "https://", id, nextStep)
		} else {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "http://", id, nextStep)
		}

		stepNext := models.Link{
			Rel:    "step-next",
			Href:   url,
			Type:   "application/json",
			Method: "POST",
		}

		links := models.Links{stepNext}

		response := models.Response{
			ID:          id,
			Workflow:    "quanta-replacement-d51b-esxi",
			CurrentStep: "start-scaleio-remove-workflow",
			Links:       links,
		}

		c.JSON(http.StatusCreated, response)
	})

	router.POST("/fru/api/workflow/:trackingid/longrunning/scaleio-remove-workflow", func(c *gin.Context) {
		id = c.Param("trackingid")

		var nextStep string
		var delay int

		// For this step, do 3 "wait" cycles. Overwrite "nextStep" and "delay"
		if longRunningCount < 3 {
			longRunningCount++
			nextStep = "longrunning/scaleio-remove-workflow"
			delay = 5

		} else {
			nextStep = steps["longrunning/scaleio-remove-workflow"]
			delay = 0
			longRunningCount = 0
		}

		var url string
		if https {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "https://", id, nextStep)
		} else {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "http://", id, nextStep)
		}

		stepNext := models.Link{
			Rel:    "step-next",
			Href:   url,
			Type:   "application/json",
			Method: "POST",
			Delay:  delay,
		}

		links := models.Links{stepNext}

		response := models.Response{
			ID:          id,
			Workflow:    "quanta-replacement-d51b-esxi",
			CurrentStep: "wait-for-scaleio-workflow",
			Links:       links,
		}

		c.JSON(http.StatusCreated, response)
	})

	router.POST("/fru/api/workflow/:trackingid/power-off-scaleio-vm", func(c *gin.Context) {
		id = c.Param("trackingid")
		nextStep := steps["power-off-scaleio-vm"]
		var url string
		if https {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "https://", id, nextStep)
		} else {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "http://", id, nextStep)
		}

		stepNext := models.Link{
			Rel:    "step-next",
			Href:   url,
			Type:   "application/json",
			Method: "POST",
		}

		links := models.Links{stepNext}

		response := models.Response{
			ID:          id,
			Workflow:    "quanta-replacement-d51b-esxi",
			CurrentStep: "power-off-scaleio-vm",
			Links:       links,
		}

		c.JSON(http.StatusCreated, response)
	})

	router.POST("/fru/api/workflow/:trackingid/enter-maintanence-mode", func(c *gin.Context) {
		id = c.Param("trackingid")
		nextStep := steps["enter-maintanence-mode"]
		var url string
		if https {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "https://", id, nextStep)
		} else {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "http://", id, nextStep)
		}

		stepNext := models.Link{
			Rel:    "step-next",
			Href:   url,
			Type:   "application/vnd.dellemc.cpsd.removenode",
			Method: "POST",
		}

		links := models.Links{stepNext}

		response := models.Response{
			ID:          id,
			Workflow:    "quanta-replacement-d51b-esxi",
			CurrentStep: "enter-maintanence-mode",
			Links:       links,
		}

		c.JSON(http.StatusCreated, response)
	})

	router.POST("/fru/api/workflow/:trackingid/instruct-physical-removal", func(c *gin.Context) {
		id = c.Param("trackingid")
		nextStep := steps["instruct-physical-removal"]
		var url string
		if https {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "https://", id, nextStep)
		} else {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "http://", id, nextStep)
		}

		stepNext := models.Link{
			Rel:    "step-next",
			Href:   url,
			Type:   "application/json",
			Method: "POST",
		}

		links := models.Links{stepNext}

		response := models.Response{
			ID:          id,
			Workflow:    "quanta-replacement-d51b-esxi",
			CurrentStep: "instruct-physical-removal",
			Links:       links,
		}

		c.JSON(http.StatusCreated, response)
	})

	router.POST("/fru/api/workflow/:trackingid/longrunning/wait-for-rackhd-discovery", func(c *gin.Context) {
		id = c.Param("trackingid")

		var nextStep string
		var delay int
		var nextType string

		// For this step, do 3 "wait" cycles. Overwrite "nextStep" and "delay"
		if longRunningCount < 3 {
			longRunningCount++
			nextStep = "longrunning/wait-for-rackhd-discovery"
			delay = 5
			nextType = "application/json"

		} else {
			nextStep = steps["longrunning/wait-for-rackhd-discovery"]
			delay = 0
			longRunningCount = 0
			nextType = "application/vnd.dellemc.nodes.list.add+json"
		}

		var url string
		if https {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "https://", id, nextStep)
		} else {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "http://", id, nextStep)
		}

		stepNext := models.Link{
			Rel:    "step-next",
			Href:   url,
			Type:   nextType,
			Method: "POST",
			Delay:  delay,
		}

		links := models.Links{stepNext}

		node1 := models.Node{
			Hostname:        "node01",
			ServiceTag:      "1234567",
			ManagementIP:    "10.10.10.1",
			PowerStatus:     "poweredOn",
			ConnectionState: "Connected",
			UUID:            "213894yu924h6ao",
		}
		node2 := models.Node{
			Hostname:        "node02",
			ServiceTag:      "98765432",
			ManagementIP:    "10.10.10.2",
			PowerStatus:     "poweredOn",
			ConnectionState: "Connected",
			UUID:            "654edcvbhnjki87",
		}
		node3 := models.Node{
			Hostname:        "node03",
			ServiceTag:      "91827465",
			ManagementIP:    "10.10.10.3",
			PowerStatus:     "poweredOn",
			ConnectionState: "Connected",
			UUID:            "fdr56765redfghj",
		}

		nodes := models.Nodes{node1, node2, node3}

		response := models.Response{
			ID:          id,
			Workflow:    "quanta-replacement-d51b-esxi",
			CurrentStep: "wait-for-rackhd-discovery",
			Links:       links,
		}

		if delay == 0 {
			response.Nodes = nodes
		}

		c.JSON(http.StatusCreated, response)
	})

	router.POST("/fru/api/workflow/:trackingid/present-system-list-add", func(c *gin.Context) {
		id = c.Param("trackingid")
		nextStep := steps["present-system-list-add"]

		var nodeToRemove models.Node
		if c.BindJSON(&nodeToRemove) == nil {

			var url string
			if https {
				url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "https://", id, nextStep)
			} else {
				url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "http://", id, nextStep)
			}

			stepNext := models.Link{
				Rel:    "step-next",
				Href:   url,
				Type:   "application/json",
				Method: "POST",
			}

			links := models.Links{stepNext}

			response := models.Response{
				ID:          id,
				Workflow:    "quanta-replacement-d51b-esxi",
				CurrentStep: "present-system-list-add",
				Links:       links,
			}

			c.JSON(http.StatusCreated, response)
		}
	})

	router.POST("/fru/api/workflow/:trackingid/configure-disks-rackhd", func(c *gin.Context) {
		id = c.Param("trackingid")
		nextStep := steps["configure-disks-rackhd"]
		var url string
		if https {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "https://", id, nextStep)
		} else {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "http://", id, nextStep)
		}

		stepNext := models.Link{
			Rel:    "step-next",
			Href:   url,
			Type:   "application/json",
			Method: "POST",
		}

		links := models.Links{stepNext}

		response := models.Response{
			ID:          id,
			Workflow:    "quanta-replacement-d51b-esxi",
			CurrentStep: "configure-disks-rackhd",
			Links:       links,
		}

		c.JSON(http.StatusCreated, response)
	})

	router.POST("/fru/api/workflow/:trackingid/install-esxi", func(c *gin.Context) {
		id = c.Param("trackingid")
		nextStep := steps["install-esxi"]
		var url string
		if https {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "https://", id, nextStep)
		} else {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "http://", id, nextStep)
		}

		stepNext := models.Link{
			Rel:    "step-next",
			Href:   url,
			Type:   "application/json",
			Method: "POST",
		}

		links := models.Links{stepNext}

		response := models.Response{
			ID:          id,
			Workflow:    "quanta-replacement-d51b-esxi",
			CurrentStep: "install-esxi",
			Links:       links,
		}

		c.JSON(http.StatusCreated, response)
	})

	router.POST("/fru/api/workflow/:trackingid/add-host-to-vcenter", func(c *gin.Context) {
		id = c.Param("trackingid")
		nextStep := steps["add-host-to-vcenter"]
		var url string
		if https {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "https://", id, nextStep)
		} else {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "http://", id, nextStep)
		}

		stepNext := models.Link{
			Rel:    "step-next",
			Href:   url,
			Type:   "application/json",
			Method: "POST",
		}

		links := models.Links{stepNext}

		response := models.Response{
			ID:          id,
			Workflow:    "quanta-replacement-d51b-esxi",
			CurrentStep: "add-host-to-vcenter",
			Links:       links,
		}

		c.JSON(http.StatusCreated, response)
	})

	router.POST("/fru/api/workflow/:trackingid/install-scaleio-vib", func(c *gin.Context) {
		id = c.Param("trackingid")
		nextStep := steps["install-scaleio-vib"]
		var url string
		if https {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "https://", id, nextStep)
		} else {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "http://", id, nextStep)
		}

		stepNext := models.Link{
			Rel:    "step-next",
			Href:   url,
			Type:   "application/json",
			Method: "POST",
		}

		links := models.Links{stepNext}

		response := models.Response{
			ID:          id,
			Workflow:    "quanta-replacement-d51b-esxi",
			CurrentStep: "install-scaleio-vib",
			Links:       links,
		}

		c.JSON(http.StatusCreated, response)
	})

	router.POST("/fru/api/workflow/:trackingid/exit-vcenter-maintenance-mode", func(c *gin.Context) {
		id = c.Param("trackingid")
		nextStep := steps["exit-vcenter-maintenance-mode"]
		var url string
		if https {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "https://", id, nextStep)
		} else {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "http://", id, nextStep)
		}

		stepNext := models.Link{
			Rel:    "step-next",
			Href:   url,
			Type:   "application/json",
			Method: "POST",
		}

		links := models.Links{stepNext}

		response := models.Response{
			ID:          id,
			Workflow:    "quanta-replacement-d51b-esxi",
			CurrentStep: "exit-vcenter-maintenance-mode",
			Links:       links,
		}

		c.JSON(http.StatusCreated, response)
	})

	router.POST("/fru/api/workflow/:trackingid/deploy-svm", func(c *gin.Context) {
		id = c.Param("trackingid")
		nextStep := steps["deploy-svm"]
		var url string
		if https {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "https://", id, nextStep)
		} else {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "http://", id, nextStep)
		}

		stepNext := models.Link{
			Rel:    "step-next",
			Href:   url,
			Type:   "application/json",
			Method: "POST",
		}

		links := models.Links{stepNext}

		response := models.Response{
			ID:          id,
			Workflow:    "quanta-replacement-d51b-esxi",
			CurrentStep: "deploy-svm",
			Links:       links,
		}

		c.JSON(http.StatusCreated, response)
	})

	router.POST("/fru/api/workflow/:trackingid/wait-for-svm-deploy", func(c *gin.Context) {
		id = c.Param("trackingid")
		nextStep := steps["wait-for-svm-deploy"]
		var url string
		if https {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "https://", id, nextStep)
		} else {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "http://", id, nextStep)
		}

		stepNext := models.Link{
			Rel:    "step-next",
			Href:   url,
			Type:   "application/json",
			Method: "POST",
		}

		links := models.Links{stepNext}

		response := models.Response{
			ID:          id,
			Workflow:    "quanta-replacement-d51b-esxi",
			CurrentStep: "wait-for-svm-deploy",
			Links:       links,
		}

		c.JSON(http.StatusCreated, response)
	})

	router.POST("/fru/api/workflow/:trackingid/start-scaleio-add-workflow", func(c *gin.Context) {
		id = c.Param("trackingid")
		nextStep := steps["start-scaleio-add-workflow"]
		var url string
		if https {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "https://", id, nextStep)
		} else {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "http://", id, nextStep)
		}

		stepNext := models.Link{
			Rel:    "step-next",
			Href:   url,
			Type:   "application/json",
			Method: "POST",
		}

		links := models.Links{stepNext}

		response := models.Response{
			ID:          id,
			Workflow:    "quanta-replacement-d51b-esxi",
			CurrentStep: "start-scaleio-add-workflow",
			Links:       links,
		}

		c.JSON(http.StatusCreated, response)
	})

	router.POST("/fru/api/workflow/:trackingid/wait-for-scaleio-add-complete", func(c *gin.Context) {
		id = c.Param("trackingid")
		nextStep := steps["wait-for-scaleio-add-complete"]
		var url string
		if https {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "https://", id, nextStep)
		} else {
			url = fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/%s", "http://", id, nextStep)
		}

		stepNext := models.Link{
			Rel:    "step-next",
			Href:   url,
			Type:   "application/json",
			Method: "POST",
		}

		links := models.Links{stepNext}

		response := models.Response{
			ID:          id,
			Workflow:    "quanta-replacement-d51b-esxi",
			CurrentStep: "wait-for-scaleio-add-complete",
			Links:       links,
		}

		c.JSON(http.StatusCreated, response)
	})

	router.POST("/fru/api/workflow/:trackingid/map-scaleio-volumes-to-host", func(c *gin.Context) {
		id = c.Param("trackingid")

		response := models.Response{
			ID:          id,
			Workflow:    "quanta-replacement-d51b-esxi",
			CurrentStep: "map-scaleio-volumes-to-host",
		}

		c.JSON(http.StatusCreated, response)
	})

	if https {
		// Find path to certs file. Should always be in same dir as mock.go
		_, filename, _, _ := runtime.Caller(0)
		dir, err := filepath.Abs(filepath.Dir(filename))
		if err != nil {
			log.Fatal(err)
		}

		go manners.ListenAndServeTLS(":8080", dir+"/cert.pem", dir+"/key.pem", router)
	} else {
		go manners.ListenAndServe(":8080", router)
	}
}
