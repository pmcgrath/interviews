#!/bin/sh

# This was just for starting, should now use the ucsvc_client go app as trying to parse http headers etc with curl is too messy 
svc_url=http://localhost:8090

echo '***** GET ALL USERS'
curl -v -X GET "$svc_url/users" -H 'Accept: application/json'

echo '***** CREATING USER'
curl -v -X POST "$svc_url/users" -H 'Content-Type: application/json' -H 'Accept: application/json' -d '{"name": "ted toe"}'

echo '***** GET ALL USERS'
curl -v -X GET "$svc_url/users" -H 'Accept: application/json'
