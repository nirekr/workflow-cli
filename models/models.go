package models

// Constants for steps
const (
	StepNext  = "step-next"
	StepRetry = "step-retry"
)

// Endpoint is a struct ...
type Endpoint struct {
	EndpointURL string `json:"endpointUrl"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

// Response is a struct ...
type Response struct {
	ID          string `json:"id,omitempty"`
	Workflow    string `json:"workflow, omitempty"`
	CurrentStep string `json:"currentStep,omitempty"`
	Nodes       Nodes  `json:"nodes"`
	Links       Links  `json:"links"`
}

// Link is a struct ...
type Link struct {
	Rel    string `json:"rel,omitempty"`
	Href   string `json:"href,omitempty"`
	Type   string `json:"type,omitempty"`
	Method string `json:"method,omitempty"`
}

// Node is a struct ...
type Node struct {
	ID           string `json:"id,omitempty"`
	SerialNumber string `json:"serial-number,omitempty"`
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

// Workflow is a ...
type Workflow struct {
	URI string `json:"uri,omitempty"`
}

// WorkflowRequest is a ...
type WorkflowRequest struct {
	Workflow string `json:"workflow"`
}
