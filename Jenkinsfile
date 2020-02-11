#!/usr/bin/env groovy

pipeline {
  agent { label 'executor-v2' }

  options {
    timestamps()
    buildDiscarder(logRotator(daysToKeepStr: '30'))
  }

  triggers {
    cron(getDailyCronString())
  }

  stages {
    stage('Validate') {
      parallel {
        stage('Changelog') {
          steps { sh './bin/parse-changelog.sh' }
        }
      }
    }

    stage('Build artifacts') {
      steps {
        sh './build.sh'
        archiveArtifacts artifacts: "dist/*.tar.gz,dist/*.zip,dist/*.rb,dist/*.deb,dist/*.rpm,dist/*.txt", fingerprint: true
      }
    }

    stage('Run unit tests') {
      steps {
        sh './bin/test.sh'
        junit 'output/junit.xml'
      }
    }
  }

  post {
    always {
      cleanupAndNotify(currentBuild.currentResult)
    }
  }
}
