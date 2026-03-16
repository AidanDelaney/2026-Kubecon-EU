#!/usr/bin/env bash

. ../../demo-magic.sh

DEMO_PROMPT="${GREEN}➜ ${CYAN}\W ${COLOR_RESET}"

# hide the evidence
clear

# enters interactive mode and allows newly typed command to be executed
cmd

pei "mvn package"

pei "mvn exec:java -Dexec.mainClass=demo.LinearDemo"