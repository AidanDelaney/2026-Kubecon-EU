#!/usr/bin/env bash

. ../../demo-magic.sh

export TMP_DIR=~/tmp

DEMO_PROMPT="${GREEN}➜ ${CYAN}\W ${COLOR_RESET}"

# hide the evidence
clear

# enters interactive mode and allows newly typed command to be executed
cmd

pei "pack build example --verbose --builder cuda-java-builder"
