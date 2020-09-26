# !/bin/sh

for i in $(find $(pwd)/test/integration/ -name *.go)
	do 
		go test $i
	done;