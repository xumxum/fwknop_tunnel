#!/bin/bash

RC_FILE="/home/john/.fwknoprc"
PORTS="tcp/655"

#This will temporarily redirect port 443(https) from this source IP, on server www.example.com to port 655(tinc), fwknop keys are in /home/john/.fwknoprc
#Also the SPA is sent on TCP port 443, instead of the fwknop's default UDP port 62201.
#Great way to hide a VPN behind only port 443
fwknop -s -p 443 -P tcp -A $PORTS --rc-file=$RC_FILE --save-args-file=/tmp/fwknop.cmd --no-save-args -n www.example.com

