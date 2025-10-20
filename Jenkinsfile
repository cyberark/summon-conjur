#!/usr/bin/env groovy
// 'product-pipelines-shared-library' draws from DevOps/product-pipelines-shared-library repository.
// 'conjur-enterprise-sharedlib' draws from Conjur-Enterprise/jenkins-pipeline-library repository.
// Point to a branch of a shared library by appending @my-branch-name to the library name
@Library(['product-pipelines-shared-library', 'conjur-enterprise-sharedlib']) _

// Automated release, promotion and dependencies
properties([
  // Include the automated release parameters for the build
  release.addParams(),
  // Dependencies of the project that should trigger builds
  dependencies([])
])

// Performs release promotion.  No other stages will be run
if (params.MODE == "PROMOTE") {
  release.promote(params.VERSION_TO_PROMOTE) { infrapool, sourceVersion, targetVersion, assetDirectory ->
    // Any assets from sourceVersion Github release are available in assetDirectory
    // Any version number updates from sourceVersion to targetVersion occur here
    // Any publishing of targetVersion artifacts occur here
    // Anything added to assetDirectory will be attached to the Github Release
  }
  
  release.copyEnterpriseRelease(params.VERSION_TO_PROMOTE)
  return
}

pipeline {
  agent { label 'conjur-enterprise-common-agent' }

  options {
    timestamps()
    buildDiscarder(logRotator(numToKeepStr: '30'))
  }

  environment {
    // Sets the MODE to the specified or autocalculated value as appropriate
    MODE = release.canonicalizeMode()
  }

  triggers {
    cron(getDailyCronString())
    parameterizedCron("""
      ${getDailyCronString("%TEST_AZURE=true;TEST_GCP=true")}
      ${getWeeklyCronString("H(1-5)", "%MODE=RELEASE")}
    """)
  }

  parameters {
    booleanParam(name: 'TEST_AZURE', defaultValue: false, description: 'Run integration tests against Azure')
    booleanParam(name: 'TEST_GCP', defaultValue: false, description: 'Run integration tests against GCP')
  }

  stages {
    // Aborts any builds triggered by another project that wouldn't include any changes
    stage ("Skip build if triggering job didn't create a release") {
      when {
        expression {
          MODE == "SKIP"
        }
      }
      steps {
        script {
          currentBuild.result = 'ABORTED'
          error("Aborting build because this build was triggered from upstream, but no release was built")
        }
      }
    }
    
    stage('Scan for internal URLs') {
      steps {
        script {
          detectInternalUrls()
        }
      }
    }

    stage('Get InfraPool ExecutorV2 Agent') {
      steps {
        script {
          // Request ExecutorV2 agents for 1 hour(s)
          infrapool = getInfraPoolAgent.connected(type: "ExecutorV2", quantity: 1, duration: 1)[0]

          // Request additional executors for cloud specific tests
          if (params.TEST_AZURE) {
            INFRAPOOL_AZURE_EXECUTORV2_AGENT_0 = getInfraPoolAgent.connected(type: "AzureExecutorV2", quantity: 1, duration: 1)[0]
          }
          if (params.TEST_GCP){
            INFRAPOOL_GCP_EXECUTORV2_AGENT_0 = getInfraPoolAgent.connected(type: "GcpExecutorV2", quantity: 1, duration: 1)[0]
          }
        }
      }
    }

    stage('Get latest upstream dependencies') {
      steps {
        script {
          updatePrivateGoDependencies("${WORKSPACE}/go.mod")
          // Copy the vendor directory onto infrapool
          infrapool.agentPut from: "vendor", to: "${WORKSPACE}"
          infrapool.agentPut from: "go.*", to: "${WORKSPACE}"
        }
      }
    }

    // Generates a VERSION file based on the current build number and latest version in CHANGELOG.md
    stage('Validate changelog and set version') {
      steps {
        script {
          updateVersion(infrapool, "CHANGELOG.md", "${BUILD_NUMBER}")

          if (params.TEST_AZURE) {
            updateVersion(INFRAPOOL_AZURE_EXECUTORV2_AGENT_0, "CHANGELOG.md", "${BUILD_NUMBER}")
          }
          if (params.TEST_GCP) {
            updateVersion(INFRAPOOL_GCP_EXECUTORV2_AGENT_0, "CHANGELOG.md", "${BUILD_NUMBER}")
          }
        }
      }
    }

    stage('Run tests') {
      environment {
        INFRAPOOL_REGISTRY_URL = "registry.tld"
        INFRAPOOL_TEST_AWS=true
      }
      steps {
        script {
          infrapool.agentSh './bin/test.sh'
          infrapool.agentStash name: 'output-xml', includes: 'output/*.xml'
          unstash 'output-xml'
          junit 'output/junit.xml'
          cobertura autoUpdateHealth: false, autoUpdateStability: false, coberturaReportFile: 'output/coverage.xml', conditionalCoverageTargets: '70, 0, 70', failUnhealthy: true, failUnstable: false, lineCoverageTargets: '70, 0, 70', maxNumberOfBuilds: 0, methodCoverageTargets: '70, 0, 70', onlyStable: false, sourceEncoding: 'ASCII', zoomCoverageChart: false
          infrapool.agentSh 'cp output/c.out .'
          codacy action: 'reportCoverage', filePath: "output/coverage.xml"
        }
      }
    }

    stage('Run Azure tests') {
      when {
        expression { params.TEST_AZURE }
      }
      environment {
        REGISTRY_URL = "registry.tld"
        INFRAPOOL_TEST_AZURE=true
      }
      steps {
        script {
          INFRAPOOL_AZURE_EXECUTORV2_AGENT_0.agentSh 'summon ./bin/test.sh'
        }
      }
    }

    stage('Run GCP tests') {
      when {
        expression { params.TEST_GCP }
      }
      environment {
        REGISTRY_URL = "registry.tld"
        INFRAPOOL_TEST_GCP=true
        INFRAPOOL_CONJUR_AUTHN_LOGIN="test-app"
      }
      steps {
        script {
          INFRAPOOL_GCP_EXECUTORV2_AGENT_0.agentSh "./bin/get_gcp_token.sh host/test-app cucumber gcp"
          INFRAPOOL_GCP_EXECUTORV2_AGENT_0.agentStash name: 'token-out', includes: "gcp/*"
          GCP_TOKEN_STASHED = true
          infrapool.agentUnstash name: 'token-out'
          infrapool.agentSh "./bin/test.sh"
        }
      }
    }

    stage('Build Release Artifacts') {
      steps {
        script {
          infrapool.agentDir('./pristine-checkout') {
            // Go releaser requires a pristine checkout
            checkout scm

            // Copy the checkout content onto infrapool
            infrapool.agentPut from: "./", to: "."

            // Copy VERSION info into prisitine folder
            infrapool.agentSh "cp ../VERSION ./VERSION"

            infrapool.agentSh './build.sh --snapshot'
            infrapool.agentArchiveArtifacts artifacts: 'dist/goreleaser/'
          }
        }
      }
    }

    stage('Release') {
      when {
        expression {
          MODE == "RELEASE"
        }
      }
      steps {
        script {
          release(infrapool) { billOfMaterialsDirectory, assetDirectory, toolsDirectory ->
            // Publish release artifacts to all the appropriate locations
            // Copy any artifacts to assetDirectory to attach them to the Github release

            // Copy assets to be published in Github release.
            // Next step: https://teams.microsoft.com/l/message/19:6f977a4fd8824acbbd91603a796bc0cf@thread.skype/1720802784680?tenantId=dc5c35ed-5102-4908-9a31-244d3e0134c6&groupId=4ef75e39-cd4a-4b26-a225-b3833f31f1b2&parentMessageId=1720011926933&teamName=Secrets%20Manager%20HQ&channelName=Infrastructure&createdTime=1720802784680
            infrapool.agentSh "${toolsDirectory}/bin/copy_goreleaser_artifacts ${assetDirectory}"

            // Create Go application SBOM using the go.mod version for the golang container image
            infrapool.agentSh """export PATH="${toolsDirectory}/bin:${PATH}" && go-bom --tools "${toolsDirectory}" --go-mod ./go.mod --image "golang" --main "cmd/" --output "${billOfMaterialsDirectory}/go-app-bom.json" """
            // Create Go module SBOM
            infrapool.agentSh """export PATH="${toolsDirectory}/bin:${PATH}" && go-bom --tools "${toolsDirectory}" --go-mod ./go.mod --image "golang" --output "${billOfMaterialsDirectory}/go-mod-bom.json" """
          }
        }
      }
    }
  }

  post {
    always {
      script {
        releaseInfraPoolAgent(".infrapool/release_agents")
        // Resolve ownership issue before running infra post hook
        sh 'git config --global --add safe.directory ${PWD}'
        infraPostHook()
      }
    }
  }
}
