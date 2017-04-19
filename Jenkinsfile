pipeline {
    agent {
        docker {
            image 'rackhd/golang:1.8.0'
            label 'maven-builder'
        }
    }
    environment {
        GIT_CREDS = credentials('github-03')
        GITHUB_TOKEN = credentials('github-02')
    }
    options { 
        buildDiscarder(logRotator(artifactDaysToKeepStr: '30', artifactNumToKeepStr: '5', daysToKeepStr: '30', numToKeepStr: '5'))
        timestamps()
    }
    stages {
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
//        stage('Unit Tests') {
//            steps {
//                sh '''
//                    cd /go/src/github.com/dellemc-symphony/workflow-cli/
//                    make unit-test
//                '''
//            }
//        }
//        stage('Integration Tests') {
//            steps {
//                sh '''
//                    cd /go/src/github.com/dellemc-symphony/workflow-cli/
//                    make integration-test
//                '''
//            }
//        }
         stage('NexB Scan') {
            steps {
		    sh "cd /go/"
                    checkout([$class: 'GitSCM', 
                              branches: [[name: '*/master']], 
                              doGenerateSubmoduleConfigurations: false, 
                              extensions: [[$class: 'RelativeTargetDirectory', relativeTargetDir: 'nexB']], 
                              submoduleCfg: [], 
                              userRemoteConfigs: [[url: 'https://github.com/nexB/scancode-toolkit.git']]])
		    sh "mkdir -p /go/nexB/nexb-output/"       
//		    sh "sh /go/nexB/scancode --help"
//                  sh "sh /go/nexB/scancode --format html ${WORKSPACE} /go/nexB/nexb-output/minimal.html"
//		    sh "sh /go/nexB/scancode --format html-app ${WORKSPACE} /go/nexB/nexb-output/scancode_result.html"
		    sh "go/nexB/scancode --help"
                    sh "go/nexB/scancode --format html /go/src/github.com/dellemc-symphony/workflow-cli /go/nexB/nexb-output/workflow-cli.html"
		    sh "go/nexB/scancode --format html-app /go/src/github.com/dellemc-symphony/workflow-cli/ /go/nexB/nexb-output/workflow-cli-grap.html"	       
//	            sh "mv /go/nexB/nexb-output/ ${WORKSPACE}/"
	       	    archiveArtifacts '**/**/nexb-output/**'
            }
        }
        stage('Third Party Audit') {
            steps {
                sh '''
                    mvn org.apache.maven.plugins:maven-dependency-plugin:2.10:analyze-report license:add-third-party org.apache.maven.plugins:maven-dependency-plugin:2.10:tree \
                    -DoutputType=dot \
                    -DoutputFile=${WORKSPACE}/report/dependency-tree.dot
                    '''   
                archiveArtifacts '**/dependency-analysis.html, **/THIRD-PARTY.txt, **/dependency-check-report.html, **/dependency-tree.dot'
            }
        }
        stage('Release') {
            when {
                environment name: "JOB_NAME", value: "workflow-cli-master"
            }
            steps {
                sh '''
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
                        --description "Workflow CLI Release"

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
}
