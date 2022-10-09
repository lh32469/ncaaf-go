def registry = "registry.container-registry:5000"
def dockerHost = "tcp://dind.container-registry:2375"

project = "ncaaf"
branch = ""
image = ""
k8sYml = ""
secretsYml = ""

pipeline {

  options {
    // Discard everything except the last 10 builds
    buildDiscarder(logRotator(numToKeepStr: '10'))
    // Don't build the same branch concurrently
    disableConcurrentBuilds()

    // Cleanup orphaned branch Kubernetes namespace
    branchTearDownExecutor 'Cleanup'
  }

  agent {
    kubernetes {
      yamlFile 'buildPod.yml'
    }
  }

  stages {

    stage('Setup') {
      steps {
        container('golang') {
          script {
            branch = env.BRANCH_NAME.toLowerCase()
            println "Project/Branch = " + project + "/" + branch
            image = "${registry}/${project}-${branch}:$BUILD_NUMBER"
            println "Image = " + image

            def k8sFile = readFile "k8s.yml"
            println "Read file k8s.yml"

            if (branch == "master") {
              host = project
            } else if (branch == "main") {
              host = project
            } else {
              host = "${branch}.${project}"
            }

            println "Host = " + host
            def binding = [
                project: project,
                branch : branch,
                image  : image,
                host   : host
            ]

            def engine = new groovy.text.SimpleTemplateEngine()
            def template = engine.createTemplate(k8sFile).make(binding)
            k8sYml = template.toString()
            println "k8sYml template created"
          }


          script {
            // Need separate script because of SimpleTemplateEngine
            // NotSerializableException
            def secretsFile = readFile "secrets.yml"

            withCredentials([string(credentialsId: 'CFDB_TOKEN', variable: 'CFDB_TOKEN')]) {

              def binding = [
                  secret: CFDB_TOKEN
              ]

              def engine = new groovy.text.SimpleTemplateEngine()
              def template = engine.createTemplate(secretsFile).make(binding)
              secretsYml = template.toString()
            }
          }

        }
      }
    }

    stage('Compile') {
      steps {
        container('golang') {
          sh 'go build'
        }
      }
    }

    stage('Docker') {
      steps {
        container('docker') {
          sh 'docker -v'
          sh "docker build -t ${registry}/${project}:${branch} \
              --build-arg PROFILE=jenkins,${branch} \
              --label job.name=$JOB_NAME ."
          sh "docker push ${registry}/${project}:${branch}"
//          sh "docker push ${image}"
          sh 'docker image ls'
          // Cleanup image(s) once pushed
          sh "docker image prune -af \
              --filter label=job.name=$JOB_NAME"
          sh 'docker image ls'
        }
      }
    }

    stage('Kubernetes') {
      steps {
        container('kubectl') {

          script {
            def namespace = "${project}-${branch}" as String

            status = sh(
                returnStatus: true,
                script: "kubectl get namespace $namespace"
            )

            if (status == 0) {
              println "$namespace namespace exists"
              sh "kubectl -n ${namespace} rollout restart deployment/${project}"
            } else {
              sh "kubectl create namespace $namespace"
            }

            writeFile file: 'secrets-out.yml', text: secretsYml
            sh "kubectl -n ${namespace} apply -f secrets-out.yml"

            writeFile file: 'k8s-out.yml', text: k8sYml
            sh "kubectl -n ${namespace} apply -f k8s-out.yml"

          }
        }
      }
    }

  }


}

