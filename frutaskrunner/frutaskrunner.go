//
// Copyright (c) 2017 Dell Inc. or its subsidiaries.  All Rights Reserved.
// Dell EMC Confidential/Proprietary Information
//
//

package frutaskrunner

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dellemc-symphony/workflow-cli/auth"
	"github.com/dellemc-symphony/workflow-cli/models"
	"github.com/dellemc-symphony/workflow-cli/transport"
	"github.com/dellemc-symphony/workflow-cli/utils"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
)

// RunTask starts the logic to run tasks and act on the results
func RunTask(r models.Response, target string) error {

	if len(r.Links) == 0 {
		return fmt.Errorf("Server error: No next-step or retry-step received")
	}

	client, clientErr := transport.NewClient(target)
	if clientErr != nil {
		log.Fatalf(clientErr.Error())
	}

	for {
		var index int
		for i, link := range r.Links {
			if link.Rel == models.StepNext {
				index = i
				break
			}
		}

		hrefSplit := strings.Split(r.Links[index].Href, "/")
		log.Infof("Next Step: %s", hrefSplit[len(hrefSplit)-1])

		endpointBody := models.Endpoint{}
		var postBody io.Reader

		// Parse for what we are getting credentials to
		typeSplit := strings.Split(r.Links[index].Type, ".")

		if typeSplit[len(typeSplit)-1] == "endpoint+json" {

			if len(typeSplit) < 4 {
				log.Printf("r is %+v", r)
				log.Warnf("Expecting Type to be of format 'application/vnd.dellemc.SOMETHING.endpoint+json', not: %s", r.Links[index].Type)
			}
			authTarget := typeSplit[len(typeSplit)-2]

			// Call auth
			endpointURL, userName, password, err := auth.TargetAuth(authTarget)
			if err != nil {
				log.Warnf("error getting creds: %s", err)
				return err
			}

			endpointBody.EndpointURL = endpointURL
			endpointBody.Username = userName
			endpointBody.Password = password

			//Format post body
			postBody, err = utils.EncodeBody(endpointBody)
			if err != nil {
				log.Warnf("Could not encode: %+v", endpointBody)
				return fmt.Errorf("error encoding credentials body: %s", err)
			}

		} else if r.Links[index].Type == "application/vnd.dellemc.nodes.list.add+json" {
			nodeSelected, err := PresentNodesToUser(models.ActionAddNode, r.Nodes)
			if err != nil {
				return err
			}

			postBody, err = utils.EncodeBody(nodeSelected)
			if err != nil {
				log.Warnf("Could not encode: %+v", nodeSelected)
				return fmt.Errorf("error encoding credentials body: %s", err)
			}

		} else if r.Links[index].Type == "application/vnd.dellemc.nodes.list.remove+json" {
			nodeSelected, err := PresentNodesToUser(models.ActionRemoveNode, r.Nodes)
			if err != nil {
				return err
			}

			postBody, err = utils.EncodeBody(nodeSelected)
			if err != nil {
				log.Warnf("Could not encode: %+v", nodeSelected)
				return fmt.Errorf("error encoding credentials body: %s", err)
			}

		} else if r.Links[index].Type == "application/vnd.dellemc.cpsd.removenode" {
			fmt.Printf("\nIt is now safe to remove the failed node from the rack.\n")
			fmt.Println("When the new node has been racked and cabled, please power it on and type CONTINUE ...")
			ok := false
			for !ok {
				scanner := bufio.NewScanner(os.Stdin)

				scanner.Scan()
				if err := scanner.Err(); err != nil {
					return fmt.Errorf("Error reading user input: %s", err)
				}

				input := strings.ToUpper(scanner.Text())

				if input == "CONTINUE" {
					ok = true
				} else {
					fmt.Println("Type CONTINUE when ready ...")
				}
			}
		}

		reqNext, err := http.NewRequest(r.Links[index].Method, r.Links[index].Href, postBody)
		if err != nil {
			log.Warnf("%s", err)
		}

		reqNext.Header.Set("content-type", r.Links[index].Type)

		resp, err := client.Do(reqNext)
		if err != nil {
			return fmt.Errorf("error sending HTTP request: %s", err)
		}
		if resp.StatusCode < 200 || 300 < resp.StatusCode {
			log.Warnf("Request was: %+v\n", reqNext)
			log.Warnf("Response is: %+v\n", resp)

			rBody, readErr := ioutil.ReadAll(resp.Body)
			if readErr != nil {
				log.Warnf("Readall error %s", readErr)
			}
			log.Warnf("Response body is: %s\n", rBody)
			return fmt.Errorf("Non-Success status(%d): %s", resp.StatusCode, resp.Status)
		}

		r = models.Response{}
		if err = utils.DecodeBody(resp, &r); err != nil {
			log.Warnf("Resp is %s\n", resp)
			return fmt.Errorf("error decoding response: %s", err)

		}

		for i, link := range r.Links {
			if link.Rel == models.StepNext {
				index = i
				break
			}
		}

		// if r.Links is empty, assume end of steps
		if len(r.Links) == 0 {
			log.Infof("Step Complete: %s", r.CurrentStep)
			log.Printf("No next-step recieved...Workflow complete")
			return nil
		}

		delay := r.Links[index].Delay

		switch r.Links[index].Rel {

		case models.StepNext:

			if delay != 0 {
				log.Infof("This task is in progress. Please wait...")
				time.Sleep(time.Duration(delay) * time.Second)
			} else {
				log.Infof("Step Complete: %s", r.CurrentStep)
			}
			fmt.Print("\n")

		case models.StepRetry:
			log.Warnf("Step Failed: %s", r.CurrentStep)

			if delay != 0 {
				log.Infof("Waiting for %d seconds before retrying", delay)
				time.Sleep(time.Duration(delay) * time.Second)
			} else {
				log.Warnf("Attempting retry: %s\n", r.CurrentStep)
			}

		default:
			log.Warnf("Status Unknown: %+v\n", r)
		}

	}
}

// InitiateWorkflow starts the workflow
func InitiateWorkflow(target string) (models.Response, error) {

	client, err := transport.NewClient(target)
	if err != nil {
		log.Fatalf(err.Error())
	}

	quantaReplace := models.WorkflowRequest{
		Workflow: "quanta-replacement-d51b-esxi",
	}

	body, err := utils.EncodeBody(quantaReplace)
	if err != nil {
		return models.Response{}, fmt.Errorf("Error encoding request: %s", err)

	}

	log.Printf("Server target is %s", target)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(target+"/fru/api/workflow"), body)
	if err != nil {
		return models.Response{}, fmt.Errorf("Error creating new request: %s", err)

	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("%s", err)
	}

	if resp.StatusCode < 200 || 300 < resp.StatusCode {
		log.Warnf("Non-Success status(%d): %s", resp.StatusCode, resp.Status)
		log.Warnf("Request was: %+v\n", req)
		log.Warnf("Response is: %+v\n", resp)

		rBody, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			log.Warnf("Readall error %s", readErr)
		}
		log.Warnf("Response body is: %s\n", rBody)

		return models.Response{}, fmt.Errorf("Non-Success status(%d): %s", resp.StatusCode, resp.Status)

	}

	r := models.Response{}
	if err = utils.DecodeBody(resp, &r); err != nil {

		return models.Response{}, fmt.Errorf("Decoding Response to Start Workflow: %s", err)
	}

	return r, nil
}

//PresentNodesToUser presents the user with a list of nodes and asks them to choose one
func PresentNodesToUser(action string, nodes models.Nodes) (models.Node, error) {
	var err error
	var selector int
	var selectedNode models.Node
	ok := false

	for !ok {
		// Print result
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Select", "Hostname", "Service Tag", "Mgmt IP", "Power Status", "Connection State"})

		for index, node := range nodes {
			table.Append([]string{
				fmt.Sprintf("%d", index+1),
				node.Hostname,
				node.ServiceTag,
				node.ManagementIP,
				node.PowerStatus,
				node.ConnectionState,
			})
		}
		table.Render()

		scanner := bufio.NewScanner(os.Stdin)
		selector = 0
		for selector < 1 || selector > len(nodes) {

			okNodeSelect := false
			for !okNodeSelect {
				fmt.Printf("Select a node for action '%s': ", action)

				scanner.Scan()
				if err = scanner.Err(); err != nil {
					log.Errorf("Error reading user input: %s", err)
					continue
				}
				input := scanner.Text()

				selector, err = strconv.Atoi(input)
				if err != nil {
					log.Errorf("Error parsing user input: \"%s\"\n", input)
					continue
				}

				okNodeSelect = true
			}
			if selector < 1 || selector > len(nodes) {
				log.Warnf("Invalid node selection: %d", selector)
			}
		}
		// subtract 1 because nodes is an array and we counted from 1 in the table
		selectedNode = nodes[selector-1]

		selectedTable := tablewriter.NewWriter(os.Stdout)
		selectedTable.SetHeader([]string{"Hostname", "Service Tag", "Mgmt IP", "Power Status", "Connection State"})
		selectedTable.Append([]string{
			selectedNode.Hostname,
			selectedNode.ServiceTag,
			selectedNode.ManagementIP,
			selectedNode.PowerStatus,
			selectedNode.ConnectionState,
		})
		fmt.Printf("\n")
		selectedTable.Render()

		fmt.Printf("Is this the correct node? [Y/N] or Q to quit: ")
		scanner.Scan()
		fmt.Printf("\n")
		if err := scanner.Err(); err != nil {
			return models.Node{}, fmt.Errorf("Error reading user input: %s", err)
		}

		input := scanner.Text()

		if strings.ToLower(input) == "q" {
			return models.Node{}, fmt.Errorf("User selected quit")
		}
		if strings.ToLower(input) == "y" {
			ok = true
		}

	}

	return selectedNode, nil
}
