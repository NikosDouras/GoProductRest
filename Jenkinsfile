pipeline {
    agent any

    environment {
        GO_ENV = 'production'                 // Set to your desired environment
        DB_HOST = 'your-db-host'              // Replace with your DB host
        POSTGRES_USER = 'your-postgres-user'  // Replace with your PostgreSQL user
        POSTGRES_PASSWORD = 'your-password'   // Replace with your PostgreSQL password
        POSTGRES_DB = 'your-database-name'    // Replace with your database name
        DB_PORT = '5432'                      // Replace with your DB port (if different)
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
                git branch: 'main', url: 'https://github.com/NikosDouras/GoProductRest.git'
            }
        }

        stage('Build') {
            steps {
                bat 'go mod download'
                bat 'go build -v ./...'
            }
        }

        stage('Test') {
            steps {
                bat 'go test -v ./...'
            }
        }

        stage('Deploy') {
            steps {
                script {
                    bat 'docker-compose down || true'  
                    bat 'docker-compose up --build -d' 
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
