#!/bin/bash

# Use to stop project.

ROOT=$(cd `dirname $0`; cd ..; pwd)
NAME=${ROOT##*/}
PIDFILE=$ROOT/run/run.pid

if [ -e $PIDFILE ]; then
    PID=$(cat $PIDFILE)
else
    PID=""
fi
if [ x$PID == x ]; then
    echo "$NAME is not running"
else
    kill $PID
    echo "$NAME stopped success"
fi
