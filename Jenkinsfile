pipeline {
    agent any

    environment {
        GO_ENV = 'production'                 // Set to your desired environment
        DB_HOST = 'db'              // Replace with your DB host
        POSTGRES_USER = 'postgres'  // Replace with your PostgreSQL user
        POSTGRES_PASSWORD = 'oOYyyha5lFkEyiWsy855'   // Replace with your PostgreSQL password
        POSTGRES_DB = 'products_db'    // Replace with your database name
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
                    // Separate the commands without '|| true' for Windows compatibility
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
