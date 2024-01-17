#!/usr/bin/env bash

BEFORE=$(git diff | wc -l)

swag fmt -d ./app/controller

AFTER=$(git diff | wc -l)

if [ "$BEFORE" != "$AFTER" ]; then
  exit 1
fi
