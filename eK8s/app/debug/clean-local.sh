#!/bin/bash
set -xe

cd mysql-initdbd && bash delete-vol.sh
bash mysql-pv/delete-vol.sh
bash redis-pv/delete-vol.sh
