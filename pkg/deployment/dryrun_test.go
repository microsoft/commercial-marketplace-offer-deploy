package deployment

import (
	"errors"
	"testing"

	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"github.com/stretchr/testify/assert"
)

type dryRunTest struct {
	validators []DryRunValidator
}

func (t *dryRunTest) loadValidatorsWithRuntimeError() {

	validatorThatReturnsError := ValidatorFunc(func(input DryRunValidationInput) (*sdk.DryRunError, error) {
		return nil, errors.New("error")
	})

	validatorThatReturnsDryRunError := ValidatorFunc(func(input DryRunValidationInput) (*sdk.DryRunError, error) {
		return &sdk.DryRunError{}, nil
	})
	t.validators = []DryRunValidator{validatorThatReturnsError, validatorThatReturnsDryRunError}
}

func (t *dryRunTest) loadSuccessValidators() {
	validatorThatReturnsDryRunError := ValidatorFunc(func(input DryRunValidationInput) (*sdk.DryRunError, error) {
		return nil, nil
	})
	t.validators = []DryRunValidator{validatorThatReturnsDryRunError, validatorThatReturnsDryRunError}
}

func newDryRunTest() *dryRunTest {
	return &dryRunTest{
		validators: []DryRunValidator{},
	}
}

func Test_validate_returns_runtime_error_if_any_validator_returns_error(t *testing.T) {
	test := newDryRunTest()
	test.loadValidatorsWithRuntimeError()

	_, err := validate(test.validators, DryRunValidationInput{})
	assert.Error(t, err)
}

func Test_validate_still_returns_dryrun_error_even_with_runtime_error(t *testing.T) {
	test := newDryRunTest()
	test.loadValidatorsWithRuntimeError()

	result, _ := validate(test.validators, DryRunValidationInput{})
	assert.Len(t, result.Errors, 1)
}

func Test_validate_returns_success_status_if_no_dryrun_errors(t *testing.T) {
	test := newDryRunTest()
	test.loadSuccessValidators()

	result, _ := validate(test.validators, DryRunValidationInput{})

	assert.Equal(t, sdk.StatusSuccess.String(), result.Status)
	assert.Len(t, result.Errors, 0)
}

func Test_validate_returns_failed_status_if_dryrun_errors(t *testing.T) {
	test := newDryRunTest()
	test.loadValidatorsWithRuntimeError()

	result, _ := validate(test.validators, DryRunValidationInput{})

	assert.Equal(t, sdk.StatusFailed.String(), result.Status)
	assert.Len(t, result.Errors, 1)
}
