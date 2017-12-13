#!/bin/bash
set -e
echo "Build Docker"
docker build -t gcr.io/hathbanger/dash-api:$TRAVIS_COMMIT .
docker tag gcr.io/hathbanger/dash-api:$TRAVIS_COMMIT gcr.io/hathbanger/dash-api:latest
echo "Authenticate Google Cloud Engine"
echo $GCLOUD_SERVICE_KEY | base64 --decode -i > ${HOME}/gcloud-service-key.json
gcloud auth activate-service-account --key-file ${HOME}/gcloud-service-key.json
echo "Configure Google Cloud Engine Environment Variables"
gcloud --quiet config set project $PROJECT_NAME
gcloud --quiet config set container/cluster $CLUSTER_NAME
gcloud --quiet config set compute/zone ${CLOUDSDK_COMPUTE_ZONE}
gcloud --quiet container clusters get-credentials $CLUSTER_NAME
echo "Push Docker Image"
gcloud docker -- push gcr.io/hathbanger/dash-api
yes | gcloud beta container images add-tag gcr.io/hathbanger/dash-api:$TRAVIS_COMMIT gcr.io/hathbanger/dash-api:latest
echo "Configure Kubernetes"
kubectl config view
kubectl config current-context
echo "Deploy"
kubectl set image deployment/dash-api dash-api=gcr.io/hathbanger/dash-api:$TRAVIS_COMMIT