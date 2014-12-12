#!/bin/bash
[[ -z $GOROOT ]] && GOROOT=/usr/local/go
if [[ -f "$GOROOT/bin/go" ]]; then
    echo "goinstalled=true"
    echo "goroot=$GOROOT"
    
    # echo `go version`
    $GOROOT/bin/go version 2>/dev/null | awk '{ printf "goversion=%s\ngoarch=%s\n", $3, $4}'
    
    # echo `go env` variables
    $GOROOT/bin/go env 2>/dev/null | awk '{split($0,v,"=");printf("goenv_%s=%s\n", tolower(v[1]), v[2])}'
else
    echo "goinstalled=false"
fi