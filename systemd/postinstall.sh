#!/bin/sh

set -e

if [ ! -f "/etc/sabapingd/sabapingd.yaml" ]; then
    cp /etc/sabapingd/sabapingd.yaml.sample /etc/sabapingd/sabapingd.yaml
fi

if [ -d /run/systemd/system ]; then
    systemctl daemon-reload
fi
