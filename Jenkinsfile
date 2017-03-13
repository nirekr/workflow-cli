pipeline {
    agent {
        docker {
            image 'rackhd/golang:1.8.0'
            label 'maven-builder'
        }
    }
    environment {
        GIT_CREDS = credentials('github-03')
        GITHUB_TOKEN = credentials('github-01')
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
            when {
                branch 'master'
            }
            steps {
                sh '''
                    go get -u github.com/aktau/github-release
                    tar -czvf release.tar.gz /go/src/github.com/dellemc-symphony/workflow-cli/bin
                    github-release release \
                        --user dellemc-symphony \
                        --repo workflow-cli \
                        --tag v0.0.1-TEST \
                        --name "TEST RELEASE UNO" \
                        --description "Release before you release!"
                    github-release upload \
                        --user dellemc-symphony \
                        --repo workflow-cli \
                        --tag v0.0.1-TEST \
                        --name "release_test" \
                        --file release.tar.gz
                '''
            }
        }
    }
}
