#!/bin/bash

TIMEFORMAT=%R
resultFolder=./results
fileResult=$resultFolder/results.csv
fileArch=$resultFolder/arch.txt

while getopts 'b:w:i:' opt; do
  case "$opt" in
    b)
      nBuffers=$((OPTARG))
      ;;
    w)
      nWorkers=$((OPTARG))
      ;;
    i)
      inputDirectory=$OPTARG
      ;;
    :)
      echo -e "Usage: $(basename $0) [-b bufferIterations] [-w MaxWorkers] [-i InputFilesDirectory]"
      exit 1
      ;;

    ?)
      echo -e "Usage: $(basename $0) [-b bufferIterations] [-w MaxWorkers] [-i InputFilesDirectory]"
      exit 1
      ;;
  esac
done
shift "$(($OPTIND -1))"

# CHECKS AND VALIDATES ARGUMENTS.
if [[ ! -v nBuffers ]] || [[ ! -v nWorkers ]] || [[ ! -v inputDirectory ]]  ;
then
    echo -e "Usage: $(basename $0) [-b bufferIterations] [-w MaxWorkers] [-i InputFilesDirectory]"
    exit 1
fi
if [[ $nBuffers < 10 ]]  ;
then
    echo -e "-b value must be equal of greater than 10."
    echo -e "Usage: $(basename $0) [-b bufferIterations] [-w MaxWorkers] [-i InputFilesDirectory]"
    exit 1
fi
if [[ $nBuffers < 0 ]]  ;
then
    echo -e "-w value must be equal of greater than 0."
    echo -e "Usage: $(basename $0) [-b bufferIterations] [-w MaxWorkers] [-i InputFilesDirectory]"
    exit 1
fi
if test ! -d $inputDirectory
then
    echo "'$inputDirectory' is not an directory."
    exit 1
fi

# SYSTEM ARCHITECTURE.
lscpu > $fileArch
printf "\n" >> $fileArch
lshw -c memory 2> /dev/null >> $fileArch

# OUTPUT FOLDER.
mkdir -p results

# LOOPING THROUGH INPUT FILES.
echo -n "" > $fileResult
for inputFile in ./input/* ; do

    # DEFINING INPUT FILE
    echo "${inputFile}," >> $fileResult;

    # LISTING GOROUTINES.
    echo -n "," >> $fileResult
    for workers in $(seq 1 $nWorkers); do  
        echo -n "$workers," >> $fileResult
    done
    echo >> $fileResult

    # SINGLE BUFFER RESULTS.
    echo -n "single buffer," >> $fileResult;
    for workers in $(seq 1 $nWorkers); do  
        ( time ../cmd/letter-counter/letter-counter -i $inputFile -o /tmp/saida.txt  -n $workers ) 2>&1 | tr ',' '.' | tr -d '\n' | cat  >> $fileResult;
        echo -n "," >> $fileResult;
    done

    printf "\n" >> $fileResult

    # BUFFERED RESULTS.
    for bf in $(seq 10 1 $nBuffers); do 

        bufferSize=$(( 2**bf ))
        echo -n "$bufferSize," >> $fileResult;

        for workers in $(seq 1 $nWorkers); do  
            ( time ../cmd/letter-counter/letter-counter -i $inputFile -o /tmp/saida.txt -buffered -b $bufferSize -n $workers ) 2>&1 | tr ',' '.' | tr -d '\n' | cat  >> $fileResult;
            echo -n "," >> $fileResult;
        done
        printf "\n" >> $fileResult
    done

    printf "\n" >> $fileResult
done

unset TIMEFORMAT;
