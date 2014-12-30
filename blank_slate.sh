#!/bin/bash
mongo alexandria --eval "db.dropDatabase()" && go clean && go build && ./alexandria --answers answers.json
