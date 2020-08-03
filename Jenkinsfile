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
        cobertura autoUpdateHealth: true, autoUpdateStability: true, coberturaReportFile: 'output/coverage.xml', conditionalCoverageTargets: '30, 0, 0', failUnhealthy: true, failUnstable: false, lineCoverageTargets: '30, 0, 0', maxNumberOfBuilds: 0, methodCoverageTargets: '30, 0, 0', onlyStable: false, sourceEncoding: 'ASCII', zoomCoverageChart: false
        sh 'cp output/c.out .'
        ccCoverage("gocov", "--prefix github.com/cyberark/summon-conjur")
      }
    }
  }

  post {
    always {
      cleanupAndNotify(currentBuild.currentResult)
    }
  }
}
