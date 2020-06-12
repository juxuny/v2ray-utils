#!/bin/bash
CORE_PATH=${HOME}/v2ray
URL="https://nbsd.best"
DATE=`date +"%Y-%m-%d"`
IGNORE=过期,到期,剩余
IGNORE_ADDR=
VMESS_FILE=${CORE_PATH}/conf.d/vmess.list
v2ray-utils crawl --user xxxxxxxx --password xxxxxx --out ${VMESS_FILE}
#update-subscription -v -u ${URL} -t "${CORE_PATH}/conf.d/template.json" -o "/Users/juxuny/Library/Application Support/Mellow/config/auto_${DATE}.json" -ignore=${IGNORE} -ignore-addr=${IGNORE_ADDR} -l=${VMESS_FILE}
v2ray-utils merge --template=${CORE_PATH}/conf.d/template.json --vmess=${VMESS_FILE} --config="${HOME}/Library/Application Support/Mellow/config/auto_${DATE}.json"