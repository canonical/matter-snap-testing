#!/bin/bash -e

# clean
pid=$(ps aux | grep '[c]hip-all-clusters-minimal-app' | awk '{print $2}')
if [ -n "$pid" ]; then
    kill "$pid"
    echo "kill"
fi

rm -f ./chip-clusters-minimal-log.txt

# start
./chip-all-clusters-minimal-app >> chip-clusters-minimal-log.txt 2>&1 &

gnome-terminal -- bash -c 'go test -v -failfast -count 1 .; exec bash'

# todo: tear down

