#!/bin/bash

hRepo=$1
sRepo=$2
action=$3
basename=`basename "$0"`

if [[ $hRepo == "" ]] || [[ $sRepo == "" ]] || [[ $action == "" ]];then
    echo "Parameters error"
    echo -e "\033[40;37m Usage: $basename \$harborRepo \$swrRepo {on|off|del|check}  [ dntact server options ] \033[0m" 
    exit 1
fi

actArry=('on' 'off' 'del' 'check')
echo ${actArry[@]} | grep $action > /dev/null
if [[ $? -ne 0 ]];then
    echo "Parameters error"
    echo "Usage: $basename \$harborRepo \$swrRepo {on|off|del|check}  [ dntact server options ]"    
else
    res=`curl -H'verification-key:swr_dnt_parse' -X POST -d "harborRepo=$hRepo&swrRepo=$sRepo&action=$action" 'http://10.13.43.158:80/service/repoaction' -s`
    echo "$res" | jq .
fi
