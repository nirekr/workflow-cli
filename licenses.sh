#!/bin/bash

go get -u github.com/pmezard/licenses

for i in `cat glide.yaml | grep package | awk '{ print $3 }'` 
do 
   licenses $i
done
