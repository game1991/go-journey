#!/bin/bash
set -eu

while [ "$i" -lt 6 ]; do
  echo "$i"
  ((i++)) || true
done