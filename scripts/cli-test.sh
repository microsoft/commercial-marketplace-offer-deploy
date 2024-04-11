
#!/bin/bash

echo "CLI tests"
echo "------------------"
echo ""


source ./scripts/cli-setup.sh

# args
# ------------------------------------------------------------

POSITIONAL_ARGS=()

while [[ $# -gt 0 ]]; do
case $1 in
    -t|--tests)
    tests="$2"
    shift # past argument
    shift # past value
    ;;
    -*|--*)
    echo "Unknown option $1"
    exit 1
    ;;
    *)
    POSITIONAL_ARGS+=("$1") # save positional arg
    shift # past argument
    ;;
esac
done
set -- "${POSITIONAL_ARGS[@]}" # restore positional parameters


function test_create_resources_archive() {
    echo "Create resources archive"
    modm util create-resources-archive -t ./templates -f src/ClientApp/ClientApp.csproj -o ./dist
}

function test_package_build() {
    echo "Build terraform complex app"
    echo "modm package build with a direct resources file path being used instead of what's in release"
    echo ""
    modm package build \
        --name "Complex terraform app" \
        --description "Simple Terraform application template that deploys a storage account" \
        --resources-file ./dist/resources.tar.gz \
        --vmi-reference-id /subscriptions/31e9f9a0-9fd2-4294-a0a3-0101246d9700/resourceGroups/modm-dev-vmi/providers/Microsoft.Compute/galleries/modm.dev.sig/images/modm/versions/0.1.155 \
        --main-template build/managedapp/terraform/complex/templates/main.tf \
        --create-ui-definition build/managedapp/terraform/complex/createUiDefinition.json \
        --out-dir ./dist

    echo "Build bicep simple app"
    echo "modm package build - with bicep simple app using a direct resources file"
    echo ""
    modm package build \
        --name "Simple bicep app" \
        --description "Simple bicep application template that deploys a storage account" \
        --resources-file ./dist/resources.tar.gz \
        --vmi-reference-id /subscriptions/31e9f9a0-9fd2-4294-a0a3-0101246d9700/resourceGroups/modm-dev-vmi/providers/Microsoft.Compute/galleries/modm.dev.sig/images/modm/versions/0.1.155 \
        --main-template build/managedapp/bicep/simple/templates/main.bicep \
        --create-ui-definition build/managedapp/bicep/simple/createUiDefinition.json \
        --out-dir ./dist
}

# if the tests are not specified, run all tests
if [ -z "$tests" ]; then
    echo ""
    echo "running all tests"
    test_create_resources_archive
    test_package_build
else
    echo "running tests: $tests"
    for test in $tests; do
        test_$test
    done
fi