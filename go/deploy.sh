#!/bin/bash

set -x

SSH_USER=isucon

echo "deploy start by $USER" | ../notify_slack.sh

for SSH_SERVER in isucon@test.asia-east1-c.isucon5-1060
do
  rsync -avz ./ $SSH_USER@$SSH_SERVER:/home/isucon/webapp/go/
  ssh -t $SSH_USER@$SSH_SERVER /bin/bash -c "/home/isucon/webapp/go/build.sh"
  ssh -t $SSH_USER@$SSH_SERVER pkill app
done

echo "deploy finished $USER" | ../notify_slack.sh
