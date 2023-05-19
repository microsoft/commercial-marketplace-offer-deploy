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

func (t *dryRunTest) loadFakeValidators() {

	validatorThatReturnsError := ValidatorFunc(func(input DryRunValidationInput) (*sdk.DryRunError, error) {
		return nil, errors.New("error")
	})

	validatorThatReturnsDryRunError := ValidatorFunc(func(input DryRunValidationInput) (*sdk.DryRunError, error) {
		return &sdk.DryRunError{}, nil
	})
	t.validators = []DryRunValidator{validatorThatReturnsError, validatorThatReturnsDryRunError}
}

func newDryRunTest() *dryRunTest {
	return &dryRunTest{
		validators: []DryRunValidator{},
	}
}

func Test_validate_returns_runtime_error_if_any_validator_returns_error(t *testing.T) {
	test := newDryRunTest()
	test.loadFakeValidators()

	_, err := validate(test.validators, DryRunValidationInput{})
	assert.Error(t, err)
}

func Test_validate_still_returns_dryrun_error_even_with_runtime_error(t *testing.T) {
	test := newDryRunTest()
	test.loadFakeValidators()

	result, _ := validate(test.validators, DryRunValidationInput{})
	assert.Len(t, result.Errors, 1)
}
