pipeline {
    agent any

    environment {
        POSTGRES_USER = 'postgres'
        POSTGRES_PASSWORD = 'oOYyyha5lFkEyiWsy855'
        POSTGRES_DB = 'products_db'
        DB_HOST = 'db'
        DB_PORT = '5432'
        GO_ENV = 'production'
    }

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
                    // Use Docker Compose to deploy the application with environment variables
                    bat '''
                    docker-compose down
                    docker-compose up --build -d
                    '''
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
