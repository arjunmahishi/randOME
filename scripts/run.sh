#!/bin/sh

set -x

# args: 
#   $1: env var name
#   $2: default value
getEnv() {
  val=$(printenv | grep "^$1=" | awk -F'=' '{print $2}')
  if [ -n "$val" ]; then
    echo $val
    return
  fi

  echo $2
}

ADDR=$(getEnv "REMOTE_WRITE_ADDR")
FREQ=$(getEnv "FREQUENCY" "5s")
CONFIG=$(getEnv "CONFIG_PATH")

echo "Starting randOME with frequency $FREQ"

CONFIG_ARG=""
if [ -n "$CONFIG" ]; then
    echo "Using config file: $CONFIG"
    CONFIG_ARG="-c $CONFIG"
fi

if [ -n "$ADDR" ]; then
    echo "Remote write address: $ADDR"
    /app/randOME remote-write -f $FREQ --addr "$ADDR" $CONFIG_ARG
else
    /app/randOME emit -f $FREQ $CONFIG_ARG
fi
