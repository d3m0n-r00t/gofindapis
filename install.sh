#!/bin/bash

git clone https://github.com/d3m0n-r00t/gofindapis.git
cd gofindapis
go build .
cp gofindapis /usr/bin/