#!/bin/bash
go build ./fwknop_tunnel.go 
mv ./fwknop_tunnel /usr/local/bin/
systemctl restart fwknop_tunnel.service

