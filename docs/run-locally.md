# Running Locally

To run the Docker image


## Test Harness

Build the test harness, which acts as a controller example
```bash
docker build . -t testharness:latest -f ./build/package/Dockerfile.testharness
```

### Calling the test harness

```
base_url=http://localhost:8280

curl $base_url/
curl $base_url/createdeployment

# now use the deployment id that's been returned:
deployment_id=<the deployment id>
curl $base_url/dryrun/$deployment_id

```