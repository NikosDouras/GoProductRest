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
                sh 'go test ./...'
            }
        }

        stage('Deploy') {
            when {
                // Only proceed if the build and tests pass
                expression { currentBuild.result == null || currentBuild.result == 'SUCCESS' }
            }
            steps {
                script {
                    // Example Docker deployment
                    // Build Docker image
                    sh 'docker build -t goproductrest:latest .'

                    // Stop the old container (if exists)
                    sh 'docker stop goproductrest || true'

                    // Remove the old container (if exists)
                    sh 'docker rm goproductrest || true'

                    // Run the new container
                    sh 'docker run -d --name goproductrest -p 8080:8080 goproductrest:latest'
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
