#! /bin/bash

set -e

APP_URL=$(oc get routes/my-web-app-route -o jsonpath='{.spec.host}')
PAUSE_TIMEOUT_SEC=1

declare -A string_set

add_string() {
    local str=$1
    if [[ -v "string_set[$str]" ]]; then
        ((string_set[$str]++))
    else
        string_set[$str]=1
    fi
}

for pc in $(seq 1 100000); do
    # curl -s $APP_URL > /dev/null
    x=$(curl -s $APP_URL | awk -F ' ' '{print $5}')
    add_string $x
    clear
    echo "Number of requests served by:"
    for str in "${!string_set[@]}"; do
        echo "$str: ${string_set[$str]}"
    done
    sleep $PAUSE_TIMEOUT_SEC
done