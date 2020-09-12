#!/bin/bash
set -e

source ./.env

curl --header "PRIVATE-TOKEN: $GITLAB_TOKEN" "https://gitlab.com/api/v4/projects/$PROJECT_ID/pipelines?per_page=$PER_PAGE&sort=asc" | jq '.[].id' > jobs.txt