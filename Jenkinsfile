pipeline {  
    agent any  
     tools {
        go 'go-1.4'
    }
    environment {
        GO111MODULE = 'on'
    }
    stages {  
        stage ('Go Version') {  
            steps {  
                sh 'go version'  
            }  
        }  
    }
}  