#!/bin/bash

ROOT=$(cd `dirname $0`; cd ..; pwd)
SRC_DIR="$ROOT/src/"
BASENAME="veronica"

SCRIPTNAME="${0##*/}"
SCRIPTNAME="${SCRIPTNAME##[KS][0-9][0-9]}"

TARGET_NAME=$1

usage() {
    echo "Usage: $SCRIPTNAME project_name"
    echo "project_name: some cool names"
}

replace_name() {
    for list in `ls $1`
    do
        list_p=$1/$list
        if [ -f $list_p ]; then
            go_file=$(ls $list_p|grep ".go$")
            if [ "$go_file"x != ""x ]; then      # replace
                echo "initial file: $list_p"
                sed -i "s/$BASENAME/$TARGET_NAME/g" $list_p
            fi
        else
            replace_name $list_p
        fi
    done
}

if [ $# != 1 ]; then
    usage
    exit 1
fi

run=`mv $SRC_DIR$BASENAME $SRC_DIR$1`
if [ $? == 1 ]; then
    echo "project has been initialized"
    exit 1
fi

replace_name $SRC_DIR$1
