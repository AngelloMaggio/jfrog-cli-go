node {

    def architectures = [
        [pkg: 'jfrog-cli-linux-386', goos: 'linux', goarch: '386', fileExtention: ''],
        [pkg: 'jfrog-cli-linux-amd64', goos: 'linux', goarch: 'amd64', fileExtention: ''],
        [pkg: 'jfrog-cli-mac-386', goos: 'darwin', goarch: 'amd64', fileExtention: ''],
        [pkg: 'jfrog-cli-windows-amd64', goos: 'windows', goarch: 'amd64', fileExtention: '.exe']
    ]

    subject = 'jfrog'
    repo = 'jfrog-cli-go'
    sh 'rm -rf temp'
    sh 'mkdir temp'
    def goRoot = tool 'go-1.7.1'
    dir('temp'){
        cliWorkspace = pwd()
        withEnv(["GOROOT=$goRoot","GOPATH=${cliWorkspace}","PATH+GOROOT=${goRoot}/bin", "JFROG_CLI_OFFER_CONFIG=false"]) {

            stage 'Go get'
            sh 'go get -f -u github.com/AngelloMaggio/jfrog-cli-go/jfrog-cli/jfrog'
            if (BRANCH?.trim()) {
                dir("src/github.com/AngelloMaggio/jfrog-cli-go/jfrog-cli/jfrog") {
                    sh "git checkout $BRANCH"
                    sh 'go install'
                }
            }
            sh 'bin/jfrog --version > version'
            version = readFile('version').trim().split(" ")[2]
            print "publishing version: $version"
            for (int i = 0; i < architectures.size(); i++) {
                    def currentBuild = architectures[i]
                    stage "Build ${currentBuild.pkg}"
                    buildAndUpload(currentBuild.goos, currentBuild.goarch, currentBuild.pkg, currentBuild.fileExtention)
            }
        }
    }
}

def uploadToBintray(pkg, fileExtension) {
     sh """#!/bin/bash
           bin/jfrog bt u $cliWorkspace/jfrog$fileExtension $subject/$repo/$pkg/$version /$version/$pkg/ --user=$USER_NAME --key=$KEY
        """
}

def buildAndUpload(goos, goarch, pkg, fileExtension) {
    sh "env GOOS=$goos GOARCH=$goarch go build github.com/AngelloMaggio/jfrog-cli-go/jfrog-cli/jfrog"
    def extension = fileExtension == null ? '' : fileExtension
    uploadToBintray(pkg, extension)
    sh "rm jfrog$extension"
}
