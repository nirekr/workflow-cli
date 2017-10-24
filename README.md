[![License](http://img.shields.io/badge/License-EPL%201.0-red.svg)](http://opensource.org/licenses/EPL-1.0)
[![Codecov](https://img.shields.io/codecov/c/github/dellemc-symphony/workflow-cli.svg)](https://codecov.io/gh/dellemc-symphony/workflow-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/dellemc-symphony/workflow-cli)](https://goreportcard.com/report/github.com/dellemc-symphony/workflow-cli)
[![Build Status](https://travis-ci.org/dellemc-symphony/workflow-cli.svg?branch=master)](https://travis-ci.org/dellemc-symphony/workflow-cli)

# workflow-cli
## Description
The workflow CLI is written in Golang. It supports native execution on a variety of hosts (Windows, Linux, and OS X). The CLI provides access to the FRU PAQX service, which facilitates FRU replacement and debugging.

## Documentation

You can find additional documentation for Project Symphony at [dellemc-symphony.readthedocs.io][documentation].

## Before you begin

Verify that the following tool is installed:

* Go Programming Language release 1.8 (go1.8) or higher  

## Building

Dell EMC strongly encourages you to build this project on a Linux or Mac environment. The instructions that follow will work on either Linux or Mac.

To install dependencies:
```
make deps
```

To build a binary for all operating systems:
```
make build
```

To build a Linux binary:
```
make build-linux
```

To build a Mac binary:
```
make build-mac
```

To build a Windows binary:
```
make build-windows
```

To run all tests:
```
make test
```

To run Integration Test:
```
make integration-test
```

To run Unit Tests:
```
make unit-test
```

The default behavior (running `make` with no commands or arguments) will run `deps`, `build`, and `test`.

## Deploying

If you want to deploy the released version of the workflow CLI, you can simply download it from the following location: https://github.com/dellemc-symphony/workflow-cli/releases

To execute the workflow CLI:

1. Change the directory to the location where you downloaded the workflow-cli executable.
2. To start the workflow CLI, type the following command:
```
./workflow-cli <command>
```
  Where *&lt;command&gt;* is one of the following:

| Command | Description |
| --- | --- |
| `fru` | Executes a field replacement workflow. The **fru** command supports several child commands, including **data**, **resume**, and **start**. The **fru data** command displays the data gathered for the FRU workflow. The **fru resume** command restarts execution of a failed FRU workflow from the last successful step in the process. The **fru start** command begins a FRU workflow.|
| `help` | Displays help about any command. |
| `status` | Retrieves the current status of the system. |
| `target` | Sets the target location (IP address and port) for the FRU PAQX service. |
| `version` | Prints version information for the CLI tool. |

To execute the workflow CLI for FRU:

1. Set the target IP address and port for the FRU PAQX service by running this command: `./workflow-cli target http://<ip address:port>`

2. Start the FRU workflow by running this command: `./workflow-cli fru start`

  The workflow CLI performs a series of steps to complete the workflow. These steps are determined dynamically by the FRU PAQX service. The CLI prompts for information, as needed, and displays informational messages to guide you through the FRU process.

  When you initiate the FRU workflow, the CLI may prompt you for the following pieces of information:
  - RackHD IP endpoint, user name, and password
  - HostBMC IP endpoint, user name, and password
  - vCenter IP endpoint, user name, and password
  - ScaleIOGateway IP endpoint, user name, and password

  The workflow CLI can read these values from the `endpoint.yaml` file, if you provide them. For any value not provided in the `endpoint.yaml` file, the workflow CLI displays a prompt to allow you to specify the value at runtime.  

After collecting the data needed to execute the FRU workflow, the CLI will prompt you to select the degraded node to remove. Once this node has been removed, the CLI will prompt you to select the new node that will replace the degraded node.   


## Contributing

Project Symphony is a collection of services and libraries housed at [GitHub][github].

Contribute code and make submissions at the relevant GitHub repository level. See [our documentation][contributing] for details on how to contribute.

## Community

Reach out to us on the Slack [#symphony][slack] channel. Request an invite at [{code}Community][codecommunity].

You can also join [Google Groups][googlegroups] and start a discussion.

[slack]: https://codecommunity.slack.com/messages/symphony
[googlegroups]: https://groups.google.com/forum/#!forum/dellemc-symphony
[codecommunity]: http://community.codedellemc.com/
[contributing]: http://dellemc-symphony.readthedocs.io/en/latest/contributingtosymphony.html
[github]: https://github.com/dellemc-symphony
[documentation]: https://dellemc-symphony.readthedocs.io/en/latest/
