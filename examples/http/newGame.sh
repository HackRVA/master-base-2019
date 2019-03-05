#!/bin/bash

##
##  Http post to create a new game
##  AbsStart must be 10 digit 
##  unix time in the future
##

curl -X POST \
  http://localhost:3000/api/newgame \
  -H 'Content-Type: application/json' \
  -d '{
    "AbsStart": 1551746973,
    "Duration": 13,
    "Variant": 1
}'