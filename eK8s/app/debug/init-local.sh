#!/bin/bash
set -xe

cd mysql-initdbd && bash create-vol.sh
bash mysql-pv/create-vol.sh
bash redis-pv/create-vol.sh
