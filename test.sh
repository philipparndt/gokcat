#!/usr/bin/env bash

topics=$(go run . -s dev topics)

for topic in $topics; do
  echo "Topic: $topic"
  go run . -s dev -t "$topic" --tail 5 | jq ". | length"
done

#go run . -s prod -t "user-service.command-response.userdata.1" --tail 5 | jq ". | length"
