#!/bin/bash
#awk -F; '{ print q $1, $5 q}' | sed -r 's/^.{10}//' /home/eugenearch/Code/github.com/eugenefoxx/test/csv/awk_csv/out.csv
# awk -v q="'" --field-separator ';' '{print q $3 "=" $1 q}' /home/eugenearch/Code/github.com/eugenefoxx/test/csv/awk_csv/out.csv
#awk 'BEGIN { FS = ";" } ; { print $1, $2 }' | sed 's/00000000000//g' /home/eugenearch/Code/github.com/eugenefoxx/test/csv/awk_csv/out.csv > /home/eugenearch/Code/github.com/eugenefoxx/test/csv/awk_csv/output1.csv
#awk 'BEGIN { FS = ";" } ; { print $1, $2 }' | sed 's/0000000000//g' /home/eugenearch/Code/github.com/eugenefoxx/test/csv/awk_csv/output1.csv > /home/eugenearch/Code/github.com/eugenefoxx/test/csv/awk_csv/output2.csv
sed 's/00000000000//g' /home/eugenearch/Code/github.com/eugenefoxx/http-rest-api-starline/import_date/out.csv > /home/eugenearch/Code/github.com/eugenefoxx/http-rest-api-starline/import_date/output1.csv
sed 's/0000000000//g' /home/eugenearch/Code/github.com/eugenefoxx/http-rest-api-starline/import_date/output1.csv > /home/eugenearch/Code/github.com/eugenefoxx/http-rest-api-starline/import_date/output2.csv