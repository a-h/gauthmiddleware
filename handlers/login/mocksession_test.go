package login

import (
	"net/http"
)

type mockSession struct {
	validateResponse             bool
	validateEmailAddressResponse string
	validateError                error
	validateWasCalled            bool
	startWasCalled               bool
}

func (ms mockSession) Validate(r *http.Request) (isValid bool, emailAddress string, err error) {
	ms.validateWasCalled = true
	return ms.validateResponse, ms.validateEmailAddressResponse, ms.validateError
}

func (ms mockSession) Start(w http.ResponseWriter, r *http.Request, emailAddress string) error {
	ms.startWasCalled = true
	return nil
}
