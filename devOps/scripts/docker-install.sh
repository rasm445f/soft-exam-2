#!/bin/bash

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# Configure firewall to enable Docker Swarm ports
yes | ufw enable

ufw allow 22/tcp
ufw allow 2376/tcp
ufw allow 2377/tcp
ufw allow 7946/tcp
ufw allow 7946/udp
ufw allow 4789/udp
ufw allow 8080/tcp

ufw reload
