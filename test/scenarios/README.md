

# Azure Test Suite
The AzureTestSuite is the base test suite for all Azure tests. It provides a common setup and teardown for all tests.

## Machine Setup
There are two required environment variables that must be set before running the tests:
```
export TEST_AZURE_SUBSCRIPTION_ID=
export TEST_AZURE_LOCATION=
```

## Running Dry Run Suite

```sh
# run the enture suite
go test -timeout 500s -run ^TestNameConflictTestSuite$ -test.v
go test -timeout 500s -run ^TestUnavailableResourceTestSuite$  -test.v
go test -timeout 120s -run ^TestDirectTemplateParamsTestSuite$  -test.v


# run a particular test
go test -timeout 500s -run ^TestNameConflictTestSuite$  -test.v -testify.m Test_Should_Fail_In_Different_Resource_Group
```
