pipeline {
    agent any

    stages {
        stage('Verify Environment') {
            steps {
                bat 'go version'
                bat 'docker --version'
                bat 'docker-compose --version'
            }
        }

        stage('Checkout') {
            steps {
                // Clone your GitHub repository and ensure it's pointing to the correct branch ('main')
                git branch: 'main', url: 'https://github.com/NikosDouras/GoProductRest.git'
            }
        }

        stage('Build') {
            steps {
                // Cache Go dependencies
                bat 'go mod download'
                // Build the Go application
                bat 'go build -v ./...'
            }
        }

        stage('Test') {
            steps {
                // Run tests
                bat 'go test -v ./...'
            }
        }

        stage('Deploy') {
            steps {
                script {
                    // Use Docker Compose to deploy the application
                    bat 'docker-compose down || true'  // Stop any running containers
                    bat 'docker-compose up --build -d' // Build and start new containers
                }
            }
        }
    }

    post {
        success {
            echo 'Build and tests passed! Deployment successful.'
        }
        failure {
            echo 'Build or tests failed. No deployment performed.'
        }
    }
}
