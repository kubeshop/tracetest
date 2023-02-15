#!/bin/bash

trap 'echo caught interrupt and exiting;exit 1' INT

TIME_OUT=$1
CONDITION=$2
IF_TRUE=$3
IF_FALSE=$4

set -ex
timeout --foreground $TIME_OUT bash -c "until ${CONDITION}; do ${IF_FALSE}; sleep 1; done; ${IF_TRUE}"
