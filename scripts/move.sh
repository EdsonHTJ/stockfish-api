#!/bin/bash
( 
echo "setoption name Skill Level $2" ;
echo "position fen $3" ;
echo "go movetime 950" ;
sleep 1
) | $1