#!/bin/bash

args=`find dataset -type f | xargs`

time bash go/concurrent-0/run.sh $args
