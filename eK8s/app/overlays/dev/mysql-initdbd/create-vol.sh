#!/bin/bash
set -xe

mkdir -p /var/lib/mysqlx/initdbd-dev
chmod 777 /var/lib/mysqlx/initdbd-dev

cp -r setup /var/lib/mysqlx/initdbd-dev
