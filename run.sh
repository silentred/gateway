#!/bin/bash

BASEDIR=$(dirname "$0")
NAME="entree"

start() {
    $BASEDIR/entree -vv=3 -log_dir=/data/logs/entree &
    echo $! > $BASEDIR/$NAME.pid
}

stop() {
    pid=`cat $BASEDIR/$NAME.pid` 
    kill -SIGTERM $pid & rm $BASEDIR/$NAME.pid
}

restart() {
    stop
    start
}

case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
esac

exit 0
