pipeline {
    agent any

    stages {
        stage('Checkout') {
            steps {
                // Clone your GitHub repository
                git url: 'https://github.com/NikosDouras/GoProductRest.git'
            }
        }

        stage('Build') {
            steps {
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

        stage('Deploy') {
            when {
                // Only proceed if the build and tests pass
                expression { currentBuild.result == null || currentBuild.result == 'SUCCESS' }
            }
            steps {
                script {
                    // Use Docker Compose to build and deploy
                    sh 'docker-compose down || true'  // Stop any running containers
                    sh 'docker-compose up --build -d' // Build and start new containers
                }
            }
        }
    }

    post {
        failure {
            // Notify on failure
            echo 'Build or tests failed. No deployment performed.'
        }
        success {
            // Notify on success
            echo 'Build and tests passed! Deployment successful.'
        }
    }
}
