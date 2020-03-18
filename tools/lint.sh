#!/bin/bash
set -e

go get -u golang.org/x/lint/golint

for d in $(go list ./... | grep -v /vendor/); do
    res=$(golint -min_confidence 0.8 $d)

    if [[ $res != '' ]]; then
        echo "$res"
        exit 1
    fi
done
