pipeline {
    agent {
        docker {
            image 'rackhd/golang:1.8.0'
            label 'maven-builder'
        }
    }
    environment {
        GIT_CREDS = credentials('github-03')
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
        stage('Release') {
            steps {
                sh 'echo "Release"'
            }
        }
    }
}
