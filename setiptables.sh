#!/bin/bash

iptables -t nat -I PREROUTING -p tcp --dport 80 -j REDIRECT --to-ports 8081
iptables -t nat -I PREROUTING -p tcp --dport 443 -j REDIRECT --to-ports 8080
echo "Rerouted Traffic"
