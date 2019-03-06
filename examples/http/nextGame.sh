#!/bin/bash

##
##  Http GET to fetch the next 
##  scheduled game
##

curl -X GET \
  http://localhost:3000/api/nextgame