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
          steps { parseChangelog() }
        }
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

    stage('Build Release Artifacts') {
      when {
        not {
          tag "v*"
        }
      }

      steps {
        sh './build.sh --snapshot'
        archiveArtifacts 'dist/goreleaser/'
      }
    }

    stage('Build Release Artifacts and Create Pre Release') {
      // Only run this stage when triggered by a tag
      when { tag "v*" }

      steps {
        dir('./pristine-checkout') {
          // Go releaser requires a pristine checkout
          checkout scm

          // Create draft release
          sh "summon --yaml 'GITHUB_TOKEN: !var github/users/conjur-jenkins/api-token' ./build.sh"
          archiveArtifacts 'dist/goreleaser/'
        }
      }
    }

  }

  post {
    always {
      cleanupAndNotify(currentBuild.currentResult)
    }
  }
}
