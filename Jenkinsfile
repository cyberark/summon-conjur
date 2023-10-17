#!/usr/bin/env groovy
@Library("product-pipelines-shared-library") _

pipeline {
  agent { label 'conjur-enterprise-common-agent' }

  options {
    timestamps()
    buildDiscarder(logRotator(daysToKeepStr: '30'))
  }

  triggers {
    cron(getDailyCronString())
  }

  stages {
    stage('Get InfraPool ExecutorV2 Agent') {
      steps {
        script {
          // Request ExecutorV2 agents for 1 hour(s)
          INFRAPOOL_EXECUTORV2_AGENT_0 = getInfraPoolAgent.connected(type: "ExecutorV2", quantity: 1, duration: 1)[0]
        }
      }
    }

    stage('Get latest upstream dependencies') {
      steps {
        script {
          withCredentials([usernamePassword(credentialsId: 'jenkins_ci_token', usernameVariable: 'GITHUB_USER', passwordVariable: 'TOKEN')]) {
            sh './bin/updateGoDependencies.sh -g "${WORKSPACE}/go.mod"'
          }
          // Copy the vendor directory onto infrapool
          INFRAPOOL_EXECUTORV2_AGENT_0.agentPut from: "vendor", to: "${WORKSPACE}"
          INFRAPOOL_EXECUTORV2_AGENT_0.agentPut from: "go.*", to: "${WORKSPACE}"
        }
      }
    }

    stage('Validate') {
      parallel {
        stage('Changelog') {
          steps { 
            parseChangelog(INFRAPOOL_EXECUTORV2_AGENT_0)
          }
        }
      }
    }

    stage('Run unit tests') {
      steps {
        script {
          INFRAPOOL_EXECUTORV2_AGENT_0.agentSh './bin/test.sh'
          INFRAPOOL_EXECUTORV2_AGENT_0.agentStash name: 'output-xml', includes: 'output/*.xml'
          unstash 'output-xml'
          junit 'output/junit.xml'
          cobertura autoUpdateHealth: true, autoUpdateStability: true, coberturaReportFile: 'output/coverage.xml', conditionalCoverageTargets: '30, 0, 0', failUnhealthy: true, failUnstable: false, lineCoverageTargets: '30, 0, 0', maxNumberOfBuilds: 0, methodCoverageTargets: '30, 0, 0', onlyStable: false, sourceEncoding: 'ASCII', zoomCoverageChart: false
          INFRAPOOL_EXECUTORV2_AGENT_0.agentSh 'cp output/c.out .'
          codacy action: 'reportCoverage', filePath: "output/coverage.xml"
        }
      }
    }

    stage('Build Release Artifacts') {
      when {
        not { buildingTag() }
      }

      steps {
        script {
          INFRAPOOL_EXECUTORV2_AGENT_0.agentSh './build.sh --snapshot'
          INFRAPOOL_EXECUTORV2_AGENT_0.agentArchiveArtifacts artifacts: 'dist/goreleaser/'
        }
      }
    }

    stage('Build Release Artifacts and Create Pre Release') {
      // Only run this stage when triggered by a tag
      when { buildingTag() }

      steps {
        script {
          INFRAPOOL_EXECUTORV2_AGENT_0.agentDir('./pristine-checkout') {
            // Go releaser requires a pristine checkout
            checkout scm

            // Copy the checkout content onto infrapool
            INFRAPOOL_EXECUTORV2_AGENT_0.agentPut from: "./", to: "."

            // Create draft release
            INFRAPOOL_EXECUTORV2_AGENT_0.agentSh 'summon --yaml "GITHUB_TOKEN: !var github/users/conjur-jenkins/api-token" ./build.sh'
          }
        }
      }
    }
  }

  post {
    always {
      script {
        releaseInfraPoolAgent(".infrapool/release_agents")
      }
    }
  }
}
