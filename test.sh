#!/bin/bash

echo "多charts"
./release/darwin/arm64/plugin-darwin-arm64 --debug --hub http://192.168.0.115:9080 --username ysicing --password aituZie3eex5fiDongoShairiangae6o  --context /Users/ysicing/Work/gitea/zcrop/pangu/charts/stable --multi
echo -e "\n多charts"
./release/darwin/arm64/plugin-darwin-arm64 --debug --hub http://192.168.0.115:9080 --username ysicing --password aituZie3eex5fiDongoShairiangae6o --context /Users/ysicing/Work/gitea/zcrop/pangu/charts/stable/ --multi --force
# echo -e "\n多charts不合法"
# ./release/darwin/arm64/plugin-darwin-arm64 --debug --hub http://192.168.0.115:9080 --username ysicing --password aituZie3eex5fiDongoShairiangae6o --context /Users/ysicing/Work/gitea/zcrop/pangu/charts/stable/common --multi
# echo -e "\n单charts"
# ./release/darwin/arm64/plugin-darwin-arm64 --debug --hub http://192.168.0.115:9080 --username ysicing --password aituZie3eex5fiDongoShairiangae6o --context /Users/ysicing/Work/gitea/zcrop/pangu/charts/stable/common
# echo -e "\n单charts"
# ./release/darwin/arm64/plugin-darwin-arm64 --debug --hub http://192.168.0.115:9080 --username ysicing --password aituZie3eex5fiDongoShairiangae6o --context /Users/ysicing/Work/gitea/zcrop/pangu/charts/stable/common/
