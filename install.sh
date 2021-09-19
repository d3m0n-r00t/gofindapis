#!/bin/bash

git clone https://github.com/d3m0n-r00t/gofindapis.git
cd gofindapis
go build .
sudo cp gofindapis /usr/bin/
# This part needs to be run as root since it requires root permission to copy this flle to /usr/bin/