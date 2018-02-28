



```
$ fly -t demo set-pipeline --pipeline demo --config pipeline.yml \
  --var private_key="$(cat .ci/id_deploy)" --var docker_password=${DOCKER_PASSWORD}

$ fly -t demo trigger-job -j demo/build

$ fly -t demo watch --job demo/build



$ git checkout execute

$ fly -t demo set-pipeline --pipeline demo --config pipeline.yml \
  --var private_key="$(cat .ci/id_deploy)" --var docker_password=${DOCKER_PASSWORD}

$ git push origin execute:master -f

$ fly -t demo execute --input demo=. --output out=tmp --config .ci/build.yml
$ fly -t demo execute --input demo=. --config .ci/test.yml

$ fly -t demo trigger-job -j demo/build
```