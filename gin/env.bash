#!/bin/bash

PROJECT_NAME=$(pwd |awk -F '/' '{print $NF}')

echo "$PROJECT_NAME"

toUpper=$(echo -e "$PROJECT_NAME" | tr '[:lower:]' '[:upper:]')

echo "$toUpper"

APPNAME="${toUpper}_TEST"

echo "$APPNAME"

OUT="$(pwd)/conf"

echo "$OUT"

export "$toUpper"="$OUT"