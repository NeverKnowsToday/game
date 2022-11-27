# !/bin/bash

#!/bin/bash

CurrentPath=$PWD
ProjectPath=$(dirname $(dirname "$PWD"))
ServerPath=$(dirname $(dirname "$PWD"))/server

set -x
docker-compose -f build.yaml  run  --rm game
mv $ServerPath/game $ProjectPath/
