pipeline {
  agent any
  stages {
    stage('Docker Build and Push') {
      steps {
        script {
          def conductorImage = docker.build("cjblink1/wordnet-world-conductor:${GIT_COMMIT}")
          conductorImage.push()
        }
      }
    }
    stage('Update Kubernetes Config') {
      when {
        branch 'master'
        not { changeRequest() }
      }
      steps {
        git 'git@github.com:wordnet-world/config.git'
        script {
          docker.image('bitnami/kubectl:1.14.1').inside('--entrypoint=\'\'') {
            sh 'cp manifests/wordnet-world-conductor-deployment.yaml manifests/wordnet-world-conductor-deployment.yaml.old'
            sh "kubectl set image -f manifests/wordnet-world-conductor-deployment.yaml.old --local wordnet-world-conductor=cjblink1/wordnet-world-conductor:${GIT_COMMIT} -o yaml > manifests/wordnet-world-conductor-deployment.yaml"
          }
        }
        sh 'git add manifests/wordnet-world-conductor-deployment.yaml'
        sh "git commit -m 'Updated conductor deployment to ${GIT_COMMIT}'"
        sh 'git push --set-upstream origin master'
      }
    }
  }
}