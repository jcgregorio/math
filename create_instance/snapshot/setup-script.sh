#!/bin/bash
#
# Script to set up a base image with just collectd, and pulld.
#
# This script is used on a temporary GCE instance. Just run it on a fresh
# Ubuntu 15.04 image and then capture a snapshot of the disk. Any image
# started with this snapshot as its image should be immediately setup to
# install applications via Push.
#
# For more details see ../../push/DESIGN.md.
sudo apt-get update
sudo apt-get --assume-yes install git
# Running "sudo apt-get --assume-yes upgrade" may upgrade the package
# gce-startup-scripts, which would cause systemd to restart gce-startup-scripts,
# which would kill this script because it is a child process of
# gce-startup-scripts.
#
# IMPORTANT: We are using a public Ubuntu image which has automatic updates
# enabled by default. Thus we are not running any commands to update packages.

sudo apt-get --assume-yes -o Dpkg::Options::="--force-confold" install collectd
sudo gsutil cp gs://mathinate-push/debs/pulld/pulld:jcgregorio@jcgregorio-glaptop-trusty:2015-11-21T21:57:51Z:2073ec3f12260757dae49f7ac44e900ba197d8bd.deb
sudo dpkg -i pulld.deb
sudo apt-get --assume-yes install --fix-broken

# Setup collectd.
sudo cat <<EOF > collectd.conf
FQDNLookup false
Interval 10

LoadPlugin "logfile"
<Plugin "logfile">
  LogLevel "info"
  File "/var/log/collectd.log"
  Timestamp true
</Plugin>

LoadPlugin syslog

<Plugin syslog>
        LogLevel info
</Plugin>

LoadPlugin battery
LoadPlugin cpu
LoadPlugin df
LoadPlugin disk
LoadPlugin entropy
LoadPlugin interface
LoadPlugin irq
LoadPlugin load
LoadPlugin memory
LoadPlugin processes
LoadPlugin swap
LoadPlugin users
LoadPlugin write_graphite

<Plugin write_graphite>
        <Carbon>
                Host "mathinate-monitoring"
                Port "2003"
                Prefix "collectd."
                StoreRates false
                AlwaysAppendDS false
                EscapeCharacter "_"
                Protocol "tcp"
        </Carbon>
</Plugin>
EOF
sudo install -D --verbose --backup=none --group=root --owner=root --mode=600 collectd.conf /etc/collectd/collectd.conf
sudo /etc/init.d/collectd restart
