# Nested Resource Id

The purpose of this integration test is to verify that issues in nested resources will be caught by the call to dryrun as long as there are no calles to the reference() function.

## mainTemplate.bicep

The mainTemplate.bicep is the source of the test input.  In this template, you will find 2 references to a storage module.  The important thing to note is that all input to these modules come from parameters that live at the same level as the module itself.

"storage2" depends on "storage1", but no input on "storage2" comes from the output of "storage1".  If it did, dryrun would ingore this relationship when doing its evaluation.

## Modifying and Executing

To run the test, modify the mainTemplate.bicep file with whatever values you like in the variables.  For instance, if you would like to test in another region besides "eastus2", simply modify the location variable.

When you are doing making modifications, you can generate the ARM template by running the command
`az bicep build -f ./mainTemplate.bicep`.  This will create an ARM template called 'mainTemplate.json' in the current directory.  This 'mainTemplate.json' template is used in the integration test.
