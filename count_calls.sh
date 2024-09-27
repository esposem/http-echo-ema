#! /bin/bash

set -e

APP_URL=$(oc get routes/http-echo-service -o jsonpath='{.spec.host}')

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
    curl -s $APP_URL > /dev/null
    # x=$(curl -s $APP_URL | awk -F '-' '{print $NF}')
    # add_string $x
    # clear
    # for str in "${!string_set[@]}"; do
    #     echo "# of pod $str replies: ${string_set[$str]}"
    # done
done