#!/bin/bash

[ ! -z "$1" ] && (
    exec /bin/bash
) 

[ -z "${PLUGIN_DEBUG}" ] && set -e || set -ex

CHARTS_DIR=${PLUGIN_CONTEXT:-.}
CHARTS_URL=${PLUGIN_URL}
PLUGIN_USERNAME=${PLUGIN_USERNAME}
PLUGIN_PASSWORD=${PLUGIN_PASSWORD}

if [ -z "${PLUGIN_URL}" ]; then
    echo "url is required"
    exit 1
fi

[ -z "${PLUGIN_FORCE}" ] && PLUGIN_FORCE="" || PLUGIN_FORCE="--force"

cd ${CHARTS_DIR}
    helm package .
cd -

if [ -z "${PLUGIN_USERNAME}" ] && [ -z "${PLUGIN_PASSWORD}" ]; then
    helm cm-push ${CHARTS_DIR} ${CHARTS_URL} ${PLUGIN_FORCE}
else 
    helm cm-push ${CHARTS_DIR} ${CHARTS_URL} --username ${PLUGIN_USERNAME} --password ${PLUGIN_PASSWORD} ${PLUGIN_FORCE}
fi