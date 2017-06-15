//
// Copyright (c) 2017 Dell Inc. or its subsidiaries.  All Rights Reserved.
// Dell EMC Confidential/Proprietary Information
//
//

pipeline {
    agent {
        docker {
            image 'rackhd/golang:1.8.0'
            label 'maven-builder'
	    customWorkspace "workspace/${env.JOB_NAME}"
        }
    }
    environment {
        GIT_CREDS = credentials('github-03')
        GITHUB_TOKEN = credentials('github-02')
        RELEASE_BRANCH = 'develop'
    }
    options { 
        skipDefaultCheckout()
        buildDiscarder(logRotator(artifactDaysToKeepStr: '30', artifactNumToKeepStr: '30', daysToKeepStr: '30', numToKeepStr: '30'))
        timestamps()
        disableConcurrentBuilds()
    }
    stages {
        stage('Checkout') {
            steps {
                doCheckout()
	    }
	}
        stage('Dependencies') {
            steps {
                sh '''
                    export GIT_SSL_NO_VERIFY=1
                    mkdir -p /go/src/github.com/dellemc-symphony/workflow-cli
                    cp -r . /go/src/github.com/dellemc-symphony/workflow-cli/
                    cd /go/src/github.com/dellemc-symphony/workflow-cli/

                    make creds
                    make deps
                '''
            }
        }
	stage('Unit Tests') {
            steps {
                sh '''
                    cd /go/src/github.com/dellemc-symphony/workflow-cli/
                    make unit-test
                '''
            }
        }
        stage('Integration Tests') {
            steps {
                sh '''
                    cd /go/src/github.com/dellemc-symphony/workflow-cli/
                    make integration-test
                '''
            }
        }
        stage('NexB Scan') {
            steps {
                doNexbScanNonMaven()
            }
        }
        stage('Release') {
            when {
                branch '${RELEASE_BRANCH}'
            }
            steps {
                sh '''
		    export BUILD_ID=$(git describe --always --dirty)
		    
                    go get -u github.com/aktau/github-release
                    cd /go/src/github.com/dellemc-symphony/workflow-cli/
                    make build

                    tar -czvf release-v0.0.1-${BUILD_ID}-windows.tgz bin/windows
                    tar -czvf release-v0.0.1-${BUILD_ID}-mac.tgz bin/darwin
                    tar -czvf release-v0.0.1-${BUILD_ID}-linux.tgz bin/linux

                    github-release release \
                        --user dellemc-symphony \
                        --repo workflow-cli \
                        --tag v0.0.1-${BUILD_ID} \
                        --name "Workflow CLI Release" \
                        --description "Workflow CLI Release" \
                        --target "${RELEASE_BRANCH}"

                    github-release upload \
                        --user dellemc-symphony \
                        --repo workflow-cli \
                        --tag v0.0.1-${BUILD_ID} \
                        --name "WorkflowCLI-Windows.tgz" \
                        --file release-v0.0.1-${BUILD_ID}-windows.tgz

                    github-release upload \
                        --user dellemc-symphony \
                        --repo workflow-cli \
                        --tag v0.0.1-${BUILD_ID} \
                        --name "WorkflowCLI-Mac.tgz" \
                        --file release-v0.0.1-${BUILD_ID}-mac.tgz

                    github-release upload \
                        --user dellemc-symphony \
                        --repo workflow-cli \
                        --tag v0.0.1-${BUILD_ID} \
                        --name "WorkflowCLI-Linux.tgz" \
                        --file release-v0.0.1-${BUILD_ID}-linux.tgz
                '''
            }
        }
    }
    post {
        always {
            cleanWorkspace()   
        }
        success {
            successEmail()
        }
        failure {
            failureEmail()
        }
    }
}
