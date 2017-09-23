#!/usr/bin/env bash
BASEDIR=$(dirname "$0")
$BASEDIR/usock2wsock -u $1 -r $2 &
sleep 2
$BASEDIR/telnet -8 -u /tmp/telnetBYwebsocket.$2.sock
