# Deployment Resource

A concourse resource that checks an http(s) endpoint for a git commit ref. Can be used to trigger actions once a service had comeup cleanly and its commit has changed (used for triggering functional tests for example).

Output is expected in the following format
```json
{
  "...": "...",
  "GitCommit": "c0fb5318590e62f7ee935f21faab7968c43dc52b",
  "...": "..."
}
```


## Source Configuration

* `endpoint`: *Required.* The endpoint used to check for a git commit id.
    Example: `https://ci.dev.cargurus.com/Cars/_deployment`


## Example
```
job:
  - name: "functional-us"
    plan:
      - aggregate:
        - get: cg-main
        - get: na-deployment
          trigger: true
    ...


resources
  - name: "na-deployment"
    type: deployment-resource
    check_every: 15s
    source:
      endpoint: https://ci.dev.cargurus.com/Cars/_deployment


resource_types:
  - name: "deployment-resource"
    type: docker-image
    source:
      repository: docker-local.cargurus.com/platform/deployment-resource
      tag: latest
```


## Tests
```sh
go test ./...
```