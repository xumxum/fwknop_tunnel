#!/bin/bash

RC_FILE="/home/john/.fwknoprc"
PORTS="tcp/22"

#This will open ssh port 22, on server www.example.com, fwknop keys are in /home/john/.fwknoprc
fwknop -s -A $PORTS --rc-file=$RC_FILE --save-args-file=/tmp/fwknop.cmd --no-save-args -n www.example.com

