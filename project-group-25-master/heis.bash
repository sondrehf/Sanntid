#!/bin/bash
echo Write id from 0:

read myID

echo Write simulator port or write 0 to run physical elevator:

read myPort 


function trap_ctrlc ()
{  
    myFlag=1
}

trap "trap_ctrlc" 2

myFlag=1

while :
do
    if [ "$myFlag" -eq 1 ]
    then
        myFlag=0
        if [ "$myPort" -eq 0 ]
        then 
            echo Running physical elevator
            go run main.go -id=$myID
        else
            echo Running simulator
            go run main.go -port=$myPort -id=$myID
        fi
    fi
 # loop inf
done
#done func