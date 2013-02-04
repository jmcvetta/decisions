#!/bin/bash
#-------------------------------------------------------------------------------
#
# Deploy Application to Heroku
#
#-------------------------------------------------------------------------------

set -x

heroku create --buildpack git://github.com/kr/heroku-buildpack-go.git
heroku addons:add mongolab:starter
heroku addons:add neo4j:try

git push heroku master
