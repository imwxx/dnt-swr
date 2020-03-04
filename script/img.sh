#!/bin/bash

pusher=$1
harborUri=$2
swrUri=$3

echo "pull start: `date "+%Y-%m-%d %H:%M:%S"`"
function PushLocal() {
    docker pull $harborUri
    if [[ $? -ne 0 ]];then
        echo "pull img failed"
        exit 1
    else
        echo "pull end: `date "+%Y-%m-%d %H:%M:%S"`"
        echo "docker tag $harborUri $swrUri"
        docker tag $harborUri $swrUri
        docker push $swrUri
        echo "push end: `date "+%Y-%m-%d %H:%M:%S"`"
        docker rmi $swrUri
    fi
}

function PushBuilder() {
    echo "docker -H tcp://$pusher:2375 tag $harborUri $swrUri 2> /dev/null"
    docker -H tcp://$pusher:2375 tag $harborUri $swrUri 2> /dev/null
    if [[ $? -ne 0 ]];then
        PushLocal
    else
       echo "push start: `date "+%Y-%m-%d %H:%M:%S"`"
       docker -H tcp://$pusher:2375 push $swrUri
       echo "push end: `date "+%Y-%m-%d %H:%M:%S"`"
       docker -H tcp://$pusher:2375 rmi $swrUri
    fi
}
 
if [[ $pusher == 'local' ]];then
    PushLocal
else
    PushBuilder
fi
