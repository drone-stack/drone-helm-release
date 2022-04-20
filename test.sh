#!/bin/bash

echo "多charts"
./release/darwin/arm64/plugin-darwin-arm64 --debug --context /Users/ysicing/Work/gitea/zcrop/pangu/charts/stable --multi
echo -e "\n多charts"
./release/darwin/arm64/plugin-darwin-arm64 --debug --context /Users/ysicing/Work/gitea/zcrop/pangu/charts/stable/ --multi
echo -e "\n多charts不合法"
./release/darwin/arm64/plugin-darwin-arm64 --debug --context /Users/ysicing/Work/gitea/zcrop/pangu/charts/stable/common --multi
echo -e "\n单charts"
./release/darwin/arm64/plugin-darwin-arm64 --debug --context /Users/ysicing/Work/gitea/zcrop/pangu/charts/stable/common
echo -e "\n单charts"
./release/darwin/arm64/plugin-darwin-arm64 --debug --context /Users/ysicing/Work/gitea/zcrop/pangu/charts/stable/common/
