#!/bin/bash
find /home/webserver/http-rest-api-starline/import_date/out.csv* -exec gunzip -c {} > /home/webserver/http-rest-api-starline/import_date/out.csv \;