#!/usr/bin/env groovy

def getRepoURL() {
    sh "git config --get remote.origin.url > .git/remote-url"
    return readFile(".git/remote-url").trim()
}

def getCommitSha() {
    sh "git rev-parse HEAD > .git/current-commit"
    return readFile(".git/current-commit").trim()
}

def updateGithubCommitStatus(build) {
    // workaround https://issues.jenkins-ci.org/browse/JENKINS-38674
    repoUrl = getRepoURL()
    commitSha = getCommitSha()
    echo repoUrl
    echo commitSha
    echo build.description
    echo "kumbirai"

    try {
        step([
                $class: 'GitHubCommitStatusSetter',
                reposSource: [$class: "ManuallyEnteredRepositorySource", url: "https://github.com/cyberark/summon-conjur"],
                commitShaSource: [$class: "ManuallyEnteredShaSource", sha: commitSha],
                errorHandlers: [[$class: 'ShallowAnyErrorHandler']],
                statusResultSource: [
                        $class: 'ConditionalStatusResultSource',
                        results: [
                                [$class: 'BetterThanOrEqualBuildResult', result: 'SUCCESS', state: 'SUCCESS', message: currentBuild.description],
                                [$class: 'BetterThanOrEqualBuildResult', result: 'FAILURE', state: 'FAILURE', message: currentBuild.description],
                                [$class: 'AnyBuildResult', state: 'FAILURE', message: 'Loophole']
                        ]
                ]
        ])


    } catch(e) {
        echo "failed because"
        echo e
    }

}


pipeline {
  agent { label 'executor-v2' }

  options {
    timestamps()
    buildDiscarder(logRotator(daysToKeepStr: '30'))
  }

  stages {
    // stage('Build Go binaries') {
    //   steps {
    //     sh './build.sh'
    //     archiveArtifacts artifacts: 'output/*', fingerprint: true
    //   }
    // }
    // stage('Run unit tests') {
    //   steps {
    //     sh './test.sh'
    //     junit 'output/junit.xml'
    //     sh 'sudo chown -R jenkins:jenkins .'  // bad docker mount creates unreadable files TODO fix this
    //   }
    // }

    // stage('Package distribution tarballs') {
    //   steps {
    //     sh './package.sh'
    //     archiveArtifacts artifacts: 'output/dist/*', fingerprint: true
    //   }
    // }

    stage('echo 2') {
      steps {
        sh 'echo 2'
      }
    }


  }

  post {
    always {
        updateGithubCommitStatus(currentBuild)
        // cleanupAndNotify(currentBuild.currentResult)
    }
  }
}
