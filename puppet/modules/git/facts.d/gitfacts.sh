#!/bin/bash
export PATH=$PATH:/usr/local/bin
which git >/dev/null 2>&1
if [[ $? -eq 0 ]]; then
    echo "gitinstalled=true"
    git --version 2>/dev/null | awk '{ printf("gitversion=%s\n", $3) }'
else
    echo "gitinstalled=false"
fi