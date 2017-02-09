#!/usr/bin/env groovy

@Library('sec_ci_libs') _

if (env.BRANCH_NAME == "v0.5.3-zkfix") {
    // Rebuild main branch once a day
    properties([
        pipelineTriggers([cron('H H * * *')])
    ])
}

task_wrapper('mesos'){
    stage("Verify author") {
        def master_branches = ['v0.5.3-zkfix']

        user_is_authorized(master_branches)
    }

    stage('Cleanup workspace') {
        deleteDir()
    }

    stage('Checkout') {
        checkout scm
    }

    load 'Jenkinsfile-insecure.groovy'
}
