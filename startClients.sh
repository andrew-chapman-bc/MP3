#!/bin/sh

for i in 1234 4567 8543 1432
do
  go run main.go --string $i &
done
go run controller.go
