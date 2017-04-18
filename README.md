[![License](http://img.shields.io/badge/License-EPL%201.0-red.svg)](http://opensource.org/licenses/EPL-1.0)

# workflow-cli
## Description
The workflow CLI is written in Golang. It supports native execution on a variety of hosts (Windows, Linux, and OS X). The CLI provides access to the FRU PAQX service, which facilitates FRU replacement and debugging.

## Documentation

The documentation is hosted at http://dellemc-symphony.readthedocs.io/.

## Before you begin
## Building
## Packaging
## Deploying

If you want to deploy the released version of the workflow CLI, you can simply download it from the following location: https://github.com/dellemc-symphony/workflow-cli/releases

To execute the workflow CLI:

1. Change the directory to the location where the workflow-cli executable has been downloaded.
2. To start the workflow CLI, type the following command: 
```
./workflow-cli <command>
```
  Where *&lt;command&gt;* is one of the following:
  
| Command | Description |
| --- | --- |
| `fru` | Executes a field replacement workflow. The **fru** command supports several child commands, including **data**, **resume**, and **start**. The **fru data** command displays the data gathered for the FRU workflow. The **fru resume** command restarts execution of a failed FRU workflow from the last successful step in the process. The **fru start** command begins a FRU workflow.|
| `status` | Retrieves the current status of the system. |
| `target` | Sets the target location (IP address and port) for the FRU PAQX service. |
| `version` | Prints version information for the CLI tool. |
 
To execute the workflow CLI for FRU: 

1. Set the target IP address and port for the FRU PAQX service by running this command: `./workflow-cli target http://<ip address:port>`

2. Start the FRU workflow by running this command: `./workflow-cli fru start`

  The workflow CLI performs a series of steps to complete the workflow. These steps are determined dynamically by the FRU PAQX service. The CLI prompts for information, as needed, and displays informational messages to guide you through the FRU process.

  When you initiate the FRU workflow, the CLI may prompt you for the following pieces of information:
  - RackHD IP endpoint, user name, and password
  - CoprHD IP endpoint, user name, and password
  - vCenter IP endpoint, user name, and password
  - ScaleIO IP endpoint, user name, and password 

  The workflow CLI can read these values from the `endpoint.yaml` file, if you provide them. For any value not provided in the `endpoint.yaml` file, the workflow CLI displays a prompt to allow you to specify the value at runtime.  
  
  At the conclusion of processing, the workflow CLI returns a UUID that represents the workflow ID. You will need this value to complete the next step to display the data (`./workflow-cli fru data <workflow-id>`).  

3. Display the data collected for the FRU workflow by running this command: `./workflow-cli fru data <workflow-id>`

 The workflow CLI returns a JSON data structure that shows the collected data. 
 
## Contributing

The Symphony project is a collection of services and libraries housed at https://github.com/dellemc-symphony.
Contribute code and make submissions at the relevant GitHub repository level. See our documentation for details on how to contribute.

## Community

Reach out to us on Slack #symphony channel. Request an invite at http://community.codedellemc.com.
You can also join [Google Groups] (https://groups.google.com/d/forum/dellemc-symphony) and start a discussion. 
