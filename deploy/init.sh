#!/usr/bin/env bash

#This scripts is to be run once per environments
# * enables the required APIs
# * create a service account and set the appropriate right to run the cloud functions
# * create the cloud tasks queues
#

echo "initializing GCP Project..."

#Cloud Functions
echo "gcloud services enable cloudfunctions.googleapis.com"
gcloud services enable cloudfunctions.googleapis.com

#Cloud Tasks
echo "gcloud services enable cloudtasks.googleapis.com"
gcloud services enable cloudtasks.googleapis.com

echo "Creating deploy service account"
gcloud iam service-accounts create github-deploy-action \
    --description="Service Account for the github deploy action" \
    --display-name="Deploy Action SA"
