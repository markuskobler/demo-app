



```
$ fly -t demo set-pipeline --pipeline=demo --config pipeline.yml --var private_key="$(cat .ci/id_deploy)"
```