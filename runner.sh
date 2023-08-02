#!/bin/bash

TIMEFORMAT=%E
nWorkers=30
nBuffers=20
fileResult=results.csv

echo -n "," > $fileResult 

for workers in $(seq 1 $nWorkers); do  
    echo -n "$workers," >> $fileResult
done
echo >> $fileResult

for bf in $(seq 10 1 $nBuffers); do 

    bufferSize=$(( 2**bf ))
    echo -n "$bufferSize," >> $fileResult;

    for workers in $(seq 1 $nWorkers); do  
        ( time ./main -i entrada3.txt -o saida.txt -b $bufferSize -n $workers ) 2>&1 | tr ',' '.' | tr -d '\n' | cat  >> $fileResult;
        echo -n "," >> $fileResult;
    done
    printf "\n" >> $fileResult;
done

unset TIMEFORMAT;
