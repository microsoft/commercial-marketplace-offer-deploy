

# Azure Test Suite
The AzureTestSuite is the base test suite for all Azure tests. It provides a common setup and teardown for all tests.

## Setup requirements
There are two required environment variables that must be set before running the tests:
```
export TEST_AZURE_SUBSCRIPTION_ID=
export TEST_AZURE_LOCATION=
```

## Running the dry run scenario test suites

the `-test.v` flag is important for manually executing the tests, since it will print out confirmations of the test pipeline, including setup and teardown.

```sh
cd ./test/scenario/dryrun

# run all
go test -timeout 500s

# run each suite
go test -timeout 500s -run ^TestNameConflictTestSuite$ -test.v
go test -timeout 500s -run ^TestUnavailableResourceTestSuite$  -test.v
go test -timeout 120s -run ^TestDirectTemplateParamsTestSuite$  -test.v
go test -timeout 300s -run ^TestNestedReferenceTestSuite$  -test.v

# run a particular test
go test -timeout 500s -run ^TestNameConflictTestSuite$  -test.v -testify.m Test_Should_Fail_In_Different_Resource_Group
```
