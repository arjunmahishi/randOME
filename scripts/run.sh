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

echo "Starting randOME with frequency $FREQ"

if [ -n "$ADDR" ]; then
    echo "Remote write address: $ADDR"
    /app/randOME remote-write -f $FREQ --addr $ADDR
else
    /app/randOME emit -f $FREQ 
fi
