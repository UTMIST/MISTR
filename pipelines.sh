#!/bin/sh
set -e

. ./.env

curl --header "PRIVATE-TOKEN: $GITLAB_TOKEN" "https://gitlab.com/api/v4/projects/$PROJECT_ID/pipelines?page=$1&per_page=$PER_PAGE&sort=asc" | jq '.[].id' > jobs.txt