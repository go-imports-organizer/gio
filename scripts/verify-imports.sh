#!/bin/bash

bad_files=$(go run main.go -l)
if [[ -n "${bad_files}" ]]; then
        echo "!!! gio needs to be run on the following files:"
        echo "${bad_files}"
        echo "Try running 'make gio'"
        exit 1
fi
