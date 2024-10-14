pipeline {
    agent any

    stages {
        stage('Verify Environment') {
            steps {
                sh 'go version'
                sh 'docker --version'
                sh 'docker-compose --version'
            }
        }

        stage('Checkout') {
            steps {
                // Clone your GitHub repository
                git url: 'https://github.com/NikosDouras/GoProductRest.git'
            }
        }

        stage('Build') {
            steps {
                // Cache Go dependencies
                sh 'go mod download'  
                // Build the Go application
                sh 'go build -v ./...'
            }
        }

        stage('Test') {
            steps {
                // Run tests
                sh 'go test -v ./...'
            }
        }
    }

    post {
        success {
            echo 'Build and tests passed! Deployment successful.'
            script {
                // Use Docker Compose to deploy the application
                sh 'docker-compose down --remove-orphans || true'  // Clean up orphan containers
                sh 'docker-compose up --build -d' // Build and start new containers
            }
        }
        failure {
            echo 'Build or tests failed. No deployment performed.'
        }
    }
}
