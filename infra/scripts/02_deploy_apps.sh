#!/bin/bash

##################
### Apps Setup ###
##################

### Set variables

# Zookeeper
declare -A zookeeper
zookeeper["name"]="zookeeper"
zookeeper["namespace"]="kafka"
zookeeper["port"]=2181

# Kafka
declare -A kafka
kafka["name"]="kafka"
kafka["namespace"]="kafka"
kafka["port"]=9092
kafka["topic"]="tracing"

# Proxy
declare -A proxy
proxy["name"]="proxy"
proxy["imageName"]="proxy-go"
proxy["namespace"]="proxy"
proxy["port"]=8080

# First
declare -A first
first["name"]="first"
first["imageName"]="first-go"
first["namespace"]="first"
first["port"]=8080

# Second
declare -A second
second["name"]="second"
second["imageName"]="second-go"
second["namespace"]="second"
second["port"]=8080

# Third
declare -A third
third["name"]="third"
third["imageName"]="third-go"
third["namespace"]="third"
third["port"]=8080

#########

####################
### Build & Push ###
####################

# Zookeeper
docker build \
  --tag "${DOCKERHUB_NAME}/${zookeeper[name]}" \
  "../../apps/kafka/zookeeper/."
docker push "${DOCKERHUB_NAME}/${zookeeper[name]}"

# Kafka
docker build \
  --tag "${DOCKERHUB_NAME}/${kafka[name]}" \
  "../../apps/kafka/kafka/."
docker push "${DOCKERHUB_NAME}/${kafka[name]}"

# Proxy
docker build \
  --tag "${DOCKERHUB_NAME}/${proxy[imageName]}" \
  "../../apps/${proxy[name]}/."
docker push "${DOCKERHUB_NAME}/${proxy[imageName]}"

# First
docker build \
  --tag "${DOCKERHUB_NAME}/${first[imageName]}" \
  "../../apps/${first[name]}/."
docker push "${DOCKERHUB_NAME}/${first[imageName]}"

# Second
docker build \
  --tag "${DOCKERHUB_NAME}/${second[imageName]}" \
  "../../apps/${second[name]}/."
docker push "${DOCKERHUB_NAME}/${second[imageName]}"

# Third
docker build \
  --tag "${DOCKERHUB_NAME}/${third[imageName]}" \
  "../../apps/${third[name]}/."
docker push "${DOCKERHUB_NAME}/${third[imageName]}"

#######

#############
### Kafka ###
#############

# Zookeeper
echo "Deploying Zookeeper ..."

helm upgrade ${zookeeper[name]} \
  --install \
  --wait \
  --debug \
  --create-namespace \
  --namespace ${zookeeper[namespace]} \
  --set dockerhubName=$DOCKERHUB_NAME \
  --set name=${zookeeper[name]} \
  --set namespace=${zookeeper[namespace]} \
  --set port=${zookeeper[port]} \
  ../charts/zookeeper

# Kafka
echo "Deploying Kafka ..."

helm upgrade ${kafka[name]} \
  --install \
  --wait \
  --debug \
  --create-namespace \
  --namespace ${kafka[namespace]} \
  --set dockerhubName=$DOCKERHUB_NAME \
  --set name=${kafka[name]} \
  --set namespace=${kafka[namespace]} \
  --set port=${kafka[port]} \
  ../charts/kafka

# Topic
echo "Checking topic [${kafka[topic]}] ..."

topicExists=$(kubectl exec -n "${kafka[namespace]}" "${kafka[name]}-0" -it -- bash \
  /kafka/bin/kafka-topics.sh \
  --bootstrap-server "${kafka[name]}.${kafka[namespace]}.svc.cluster.local:${kafka[port]}" \
  --list \
  | grep ${kafka[topic]})

if [[ $topicExists == "" ]]; then

  echo " -> Topic does not exist. Creating ..."
  while :
  do
    isTopicCreated=$(kubectl exec -n "${kafka[namespace]}" "${kafka[name]}-0" -it -- bash \
      /kafka/bin/kafka-topics.sh \
      --bootstrap-server "${kafka[name]}.${kafka[namespace]}.svc.cluster.local:${kafka[port]}" \
      --create \
      --topic ${kafka[topic]} \
      2> /dev/null)

    if [[ $isTopicCreated == "" ]]; then
      echo " -> Kafka pods are not fully ready yet. Waiting ..."
      sleep 2
      continue
    fi

    echo -e " -> Topic is created successfully.\n"
    break

  done
else
  echo -e " -> Topic already exists.\n"
fi
#########

#############
### Proxy ###
#############
echo "Deploying proxy..."

helm upgrade ${proxy[name]} \
  --install \
  --wait \
  --debug \
  --create-namespace \
  --namespace ${proxy[namespace]} \
  --set dockerhubName=$DOCKERHUB_NAME \
  --set name=${proxy[name]} \
  --set imageName=${proxy[imageName]} \
  --set namespace=${proxy[namespace]} \
  --set port=${proxy[port]} \
  --set newRelicLicenseKey=$NEWRELIC_LICENSE_KEY \
  "../charts/${proxy[name]}"
#########

#################
### First app ###
#################
echo "Deploying first app..."

helm upgrade ${first[name]} \
  --install \
  --wait \
  --debug \
  --create-namespace \
  --namespace ${first[namespace]} \
  --set dockerhubName=$DOCKERHUB_NAME \
  --set name=${first[name]} \
  --set imageName=${first[imageName]} \
  --set namespace=${first[namespace]} \
  --set port=${first[port]} \
  "../charts/${first[name]}"
#########

#################
### Second app ###
##################
echo "Deploying second app..."

helm upgrade ${second[name]} \
  --install \
  --wait \
  --debug \
  --create-namespace \
  --namespace ${second[namespace]} \
  --set dockerhubName=$DOCKERHUB_NAME \
  --set name=${second[name]} \
  --set namespace=${second[namespace]} \
  --set port=${second[port]} \
  --set newRelicLicenseKey=$NEWRELIC_LICENSE_KEY \
  "../charts/${second[name]}"
#########

#################
### Third App ###
#################
echo "Deploying third app..."

helm upgrade ${third[name]} \
  --install \
  --wait \
  --debug \
  --create-namespace \
  --namespace ${third[namespace]} \
  --set dockerhubName=$DOCKERHUB_NAME \
  --set name=${third[name]} \
  --set namespace=${third[namespace]} \
  --set port=${third[port]} \
  --set newRelicLicenseKey=$NEWRELIC_LICENSE_KEY \
  "../charts/${third[name]}"
#########
