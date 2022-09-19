pipeline {
    agent any

    tools {
        gradle 'Gradle'
    }

    stages {
        stage('build') {
            steps {
                echo 'building the application...'
            }
        }

        stage('test') {
            steps {
                echo 'testing the application...'
            }
        }

        stage('deployment') {
            steps {
                echo 'deploying the application...'
            }
        }
    }
}