package frutaskrunner

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dellemc-symphony/workflow-cli/auth"
	"github.com/dellemc-symphony/workflow-cli/models"
	"github.com/dellemc-symphony/workflow-cli/utils"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	log "github.com/sirupsen/logrus"
)

// RunTask starts the logic to run tasks and act on the results
func RunTask(r models.Response) error {

	if len(r.Links) == 0 {
		return fmt.Errorf("Server error: No next-step or retry-step received")
	}

	client := cleanhttp.DefaultClient()
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
			endpoint, userName, password, err := auth.TargetAuth(authTarget)
			if err != nil {
				log.Warnf("error getting creds: %s", err)
			}

			endpointBody.Endpoint = endpoint
			endpointBody.Username = userName
			endpointBody.Password = password

			//Format post body
			postBody, err = utils.EncodeBody(endpointBody)
			if err != nil {
				log.Warnf("Could not encode: %+v", endpointBody)
				return fmt.Errorf("error encoding credentials body: %s", err)
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

		if len(r.Nodes) != 0 {
			log.Printf("Received node data:")
			for _, node := range r.Nodes {
				fmt.Printf("  %+v\n", node)
			}
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

		switch r.Links[index].Rel {

		case models.StepNext:
			log.Infof("Step Complete: %s\n", r.CurrentStep)

		case models.StepRetry:
			log.Warnf("Step Failed: %s", r.CurrentStep)
			log.Warnf("Attempting retry: %s\n", r.CurrentStep)

		default:
			log.Warnf("Status Unknown: %+v\n", r)
		}

	}
}

// InitiateWorkflow starts the workflow
func InitiateWorkflow(target string) (models.Response, error) {

	client := cleanhttp.DefaultClient()

	quantaReplace := models.WorkflowRequest{
		Workflow: "quanta-replacement-d51b-esxi",
	}

	body, err := utils.EncodeBody(quantaReplace)
	if err != nil {
		return models.Response{}, fmt.Errorf("Error encoding request: %s", err)

	}

	log.Printf("Server target is %s\n", target)
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
