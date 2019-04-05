package gobetafaceapi

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

const (
	code    = "test.code"
	message = "test.message"
)

type ErrorTestSuite struct {
	suite.Suite
	err *Error
}

func (suite *ErrorTestSuite) SetupTest() {
	// use the native HTTPClient by default
	suite.err = NewError(code, message)
}

func (suite *ErrorTestSuite) TestError() {
	suite.Equal(message, suite.err.Error())
}

func (suite *ErrorTestSuite) TestCode() {
	suite.Equal(code, suite.err.Code())
}

// ************************************************************************************************
// ** Run Test Suite
// ************************************************************************************************

func TestErrorRunSuite(t *testing.T) {
	suite.Run(t, new(ErrorTestSuite))
}
