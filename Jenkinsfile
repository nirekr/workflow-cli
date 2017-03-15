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
                   #export GIT_SSL_NO_VERIFY=1
                   #mkdir -p /go/src/github.com/dellemc-symphony/workflow-cli
                   #cp -r . /go/src/github.com/dellemc-symphony/workflow-cli/
                   #cd /go/src/github.com/dellemc-symphony/workflow-cli/
                   
                   printenv
                   echo $GIT_BRANCH

                   #make creds
                   #make deps
                '''
            }
        }

        stage('Release') {
            when {
                branch "master"
            }
            steps {
                sh '''
                    echo "RELEASE~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~"
                '''
            }
        }
    }
}
