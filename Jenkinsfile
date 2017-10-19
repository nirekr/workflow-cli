//
// Copyright (c) 2017 Dell Inc. or its subsidiaries.  All Rights Reserved.
// Dell EMC Confidential/Proprietary Information
//
//

pipeline {
    agent {
        docker {
            image 'rackhd/golang:1.8.3'
            label 'maven-builder'
	    customWorkspace "workspace/${env.JOB_NAME}"
        }
    }
    environment {
        GIT_CREDS = credentials('github-03')
        GITHUB_TOKEN = credentials('github-02')
        RELEASE_BRANCH = 'master'
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
                checkout([$class: 'GitSCM', branches: [[name: env.BRANCH_NAME]], 
			  doGenerateSubmoduleConfigurations: false, 
			  extensions: [[$class: 'CloneOption', depth: 0, noTags: false, reference: '', shallow: false], [$class: 'AuthorInChangelog']], 
			  gitTool: 'linux-git', submoduleCfg: [], 
			  userRemoteConfigs: [[credentialsId: 'github-oauth-token', url: 'https://github.com/dellemc-symphony/workflow-cli']]])
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
	    
	stage('Code Coverage') {
            steps {
                sh '''
                    cd /go/src/github.com/dellemc-symphony/workflow-cli/
                    make coverage
		    mkdir -p ${WORKSPACE}/Cobcov
		    find . -name '*coverage*.xml' -exec cp {} ${WORKSPACE}/Cobcov  \\;
		'''
		 step([$class: 'CoberturaPublisher', autoUpdateHealth: false, autoUpdateStability: false, coberturaReportFile: '**/Cobcov/*.xml',  failNoReports: false, failUnhealthy: false, failUnstable: false, maxNumberOfBuilds: 0, sourceEncoding: 'ASCII', zoomCoverageChart: false])
	       }
            }

	stage('Licenses') {
            steps {
                sh '''
                   cd /go/src/github.com/dellemc-symphony/workflow-cli/
                   mkdir -p ${WORKSPACE}/target/generated-sources/license
                   make licenses
                   cd ${WORKSPACE}
                '''
                archiveArtifacts '**/target/**'
	       }
           }
	    
        stage('NexB Scan') {
             steps {
                    checkout([$class: 'GitSCM',
                              branches: [[name: '*/master']],
                              doGenerateSubmoduleConfigurations: false,
                              extensions: [[$class: 'RelativeTargetDirectory', relativeTargetDir: 'nexB']],
                              submoduleCfg: [],
                              userRemoteConfigs: [[url: 'https://github.com/nexB/scancode-toolkit.git']]])
		     checkout changelog: false, poll: false, scm: [$class: 'GitSCM',
			      branches: [[name: '*/master']],
			      doGenerateSubmoduleConfigurations: false,
			      extensions: [[$class: 'RelativeTargetDirectory', relativeTargetDir: 'workflow-cli']],
			      gitTool: 'linux-git', submoduleCfg: [],
			      userRemoteConfigs: [[credentialsId: 'github-03', url: 'https://github.com/dellemc-symphony/workflow-cli.git']]]

		    sh "mkdir -p  ${WORKSPACE}/nexb-output/"
       		    sh "nexB/scancode --help"
		    sh "nexB/scancode --format html workflow-cli ${WORKSPACE}/nexb-output/workflow-cli.html"
		    sh "nexB/scancode --format html-app workflow-cli ${WORKSPACE}/nexb-output/workflow-cli-grap.html"
		    archiveArtifacts '**/nexb-output/**'

            }
        }
        stage('Release') {
            when {
                branch '${RELEASE_BRANCH}'
            }
            steps {
                sh '''
                    # Decide if bumping Major, Minor, or Patch
                    LAST_COMMIT=$(git log -1 --pretty=%B)

                    BUMP=""

                    # If number of times "MAJOR" appears is greater or equal to 1
                    if [ `echo ${LAST_COMMIT}  | grep -c "MAJOR"` -ge 1 ]; then
                        BUMP=M

                    elif [ `echo ${LAST_COMMIT}  | grep -c "MINOR"` -ge 1 ]; then
                        BUMP=m

                    # Default to patch bump
                    else
                        BUMP=p

                    fi

                    # Get new version number
                    NEW_VERSION=$(increment_version.sh -$BUMP $(git describe --abbrev=0 --tag))

                    go get -u github.com/aktau/github-release
                    cd /go/src/github.com/dellemc-symphony/workflow-cli/
                    make build

                    cd bin/windows && zip ../../release-$NEW_VERSION-windows.zip ./* && cd ../../
                    tar -czvf release-$NEW_VERSION-mac.tgz bin/darwin
                    tar -czvf release-$NEW_VERSION-linux.tgz bin/linux

                    github-release release \
                        --user dellemc-symphony \
                        --repo workflow-cli \
                        --tag $NEW_VERSION \
                        --name "Workflow CLI Release" \
                        --description "Workflow CLI Release" \
                        --target "${RELEASE_BRANCH}"

                    github-release upload \
                        --user dellemc-symphony \
                        --repo workflow-cli \
                        --tag $NEW_VERSION \
                        --name "WorkflowCLI-Windows.zip" \
                        --file release-$NEW_VERSION-windows.zip

                    github-release upload \
                        --user dellemc-symphony \
                        --repo workflow-cli \
                        --tag $NEW_VERSION \
                        --name "WorkflowCLI-Mac.tgz" \
                        --file release-$NEW_VERSION-mac.tgz

                    github-release upload \
                        --user dellemc-symphony \
                        --repo workflow-cli \
                        --tag $NEW_VERSION \
                        --name "WorkflowCLI-Linux.tgz" \
                        --file release-$NEW_VERSION-linux.tgz
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
