#!/bin/bash
./rebuild.sh
./fwknop_tunnel --version
INSTALL_PATH="/usr/local/bin"

echo "Moving binary to $INSTALL_PATH"
mv ./fwknop_tunnel $INSTALL_PATH

echo "Restarting fwknop_tunnel.service"
systemctl restart fwknop_tunnel.service

