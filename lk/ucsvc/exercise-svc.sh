#!/bin/sh

# This was just for starting, should now use the ucsvc_client go app as trying to parse http headers etc with curl is too messy 
svc_url=https://localhost:8090
cert_file_path=../certs/server.crt
user_name=ted
password=toe
credentials=$user_name:$password

echo '***** GET ALL USERS'
curl -v --cacert $cert_file_path -u $credentials -X GET "$svc_url/users" -H 'Accept: application/json'

echo '***** CREATING USER'
curl -v --cacert $cert_file_path -u $credentials -X POST "$svc_url/users" -H 'Content-Type: application/json' -H 'Accept: application/json' -d '{"name": "larry murphy"}'

echo '***** GET ALL USERS'
curl -v --cacert $cert_file_path -u $credentials -X GET "$svc_url/users" -H 'Accept: application/json'
