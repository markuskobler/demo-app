



```
$ fly -t demo set-pipeline --pipeline demo --config pipeline.yml \
  --var private_key="$(cat .ci/id_deploy)" --var docker_password=${DOCKER_PASSWORD}
```