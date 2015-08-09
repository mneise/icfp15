#!/bin/bash
for i in {0..24}
do
    result="$(./play_icfp2015 -d=false -f p$i.json)"
    echo $result > body.json
    printf "Trying to submit solution for problem $i\n"
    curl --user :EuNgF0Z5GnOGCHVcKtd/iLIsqzE4Kyk0nSLPzC9pukY= -X POST -H "Content-Type: application/json" -d @body.json https://davar.icfpcontest.org/teams/260/solutions
    printf "\n"
done
