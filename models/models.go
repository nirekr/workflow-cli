//
// Copyright (c) 2017 Dell Inc. or its subsidiaries.  All Rights Reserved.
// Dell EMC Confidential/Proprietary Information
//
//

package models

// Constants for steps
const (
	StepNext  = "step-next"
	StepRetry = "step-retry"
)

// Constants for node actions
const (
	ActionAddNode    = "Add Node"
	ActionRemoveNode = "Remove Node"
)

// Endpoint is a struct ...
type Endpoint struct {
	EndpointURL string `json:"endpointUrl"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

// Response is a struct ...
type Response struct {
	ID          string     `json:"id,omitempty"`
	Workflow    string     `json:"workflow, omitempty"`
	CurrentStep string     `json:"currentStep,omitempty"`
	Nodes       Nodes      `json:"nodes"`
	Links       Links      `json:"links"`
	Networking  Networking `json:"networking"`
}

// Link is a struct ...
type Link struct {
	Rel    string `json:"rel,omitempty"`
	Href   string `json:"href,omitempty"`
	Type   string `json:"type,omitempty"`
	Method string `json:"method,omitempty"`
	Delay  int    `json:"nextStepDelay,omitempty"`
}

// Node is a struct ...
type Node struct {
	Hostname        string `json:"hostname,omitempty"`
	ServiceTag      string `json:"serviceTag,omitempty"`
	ManagementIP    string `json:"mgmtIP,omitempty"`
	PowerStatus     string `json:"powerStatus,omitempty"`
	ConnectionState string `json:"connectionState,omitempty"`
	UUID            string `json:"uuid,omitempty"`
}

// Networking is a struct ...
type Networking struct {
	Hostname    string      `json:"hostname,omitempty"`
	DvsNetworks DvsNetworks `json:"dvsNetworkList,omitempty"`
}

// VmkAdapter is a struct ...
type VmkAdapter struct {
	Device         string `json:"device,omitempty"`
	Network        string `json:"network,omitempty"`
	IpAddress      string `json:"ipAddress,omitempty"`
	SubnetMask     string `json:"subnetMask,omitempty"`
	Mtu            string `json:"mtu,omitempty"`
	EnabledService string `json:"enabledService,omitempty"`
}

// DvsNetwork is a struct ...
type DvsNetwork struct {
	DvsName      string      `json:"dvsName,omitempty"`
	PhysicalNics []string    `json:"physicalNics,omitempty"`
	VmkAdapters  VmkAdapters `json:"vmkAdapters,omitempty"`
}

// MockDataResponse is a struct ...
type MockDataResponse struct {
	Response string `json:"response,omitempty"`
	TaskID   string `json:"taskid,omitempty"`
}

// Links is an array of link objs
type Links []Link

// Nodes is an array of node objs
type Nodes []Node

// Workflows is a ...
type Workflows []Workflow

// DvsNetworks is a ...
type DvsNetworks []DvsNetwork

// VmkAdapters is a ...
type VmkAdapters []VmkAdapter

// Workflow is a ...
type Workflow struct {
	URI string `json:"uri,omitempty"`
}

// WorkflowRequest is a ...
type WorkflowRequest struct {
	Workflow string `json:"workflow"`
}
