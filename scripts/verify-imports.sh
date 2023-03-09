#!/bin/bash

bad_files=$(go run main.go -l)
if [[ -n "${bad_files}" ]]; then
        echo "!!! goio needs to be run on the following files:"
        echo "${bad_files}"
        echo "Try running 'make imports'"
        exit 1
fi
