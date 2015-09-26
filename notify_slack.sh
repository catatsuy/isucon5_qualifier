#!/bin/bash

exec tee >(while read line; do curl -s --data-binary '`'"$line"'`' 'https://teamfreesozai.slack.com/services/hooks/slackbot?token=kwiHIj3FF1a23MtmwW7Ii3Ji&channel=%23general' -o/dev/null; done)
