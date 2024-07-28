#!/usr/bin/env bash

export SOCIETY_VOTING_RESTART_ENABLED=yes

STATUS=153

while [[ "$STATUS" == "153" ]]; do
  $1
  STATUS=$?
done
