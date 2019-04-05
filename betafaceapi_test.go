package gobetafaceapi

import (
	"testing"

	"github.com/enriquebris/gohttpclient"
	"github.com/stretchr/testify/suite"
)

const (
	key               = "apiKEY"
	secret            = "apiSECRET"
	getMediaUUID      = "8380b14d-1ba1-4355-b37c-3525a6f28620"
	sendMediaFileURI  = "https://domain.com/image.jpg"
	sendMediaFilename = "image.jpg"
)

type NativeHTTPClientTestSuite struct {
	suite.Suite
	client *BasicClient
}

func (suite *NativeHTTPClientTestSuite) SetupTest() {
	// use the native HTTPClient by default
	suite.client = NewNativeHTTPClient(key, secret)
}

// ************************************************************************************************
// ** initialize
// ************************************************************************************************

// TestInitialize tests initialize(). It verifies the apiKey/apiSecret and the default endpoints.
func (suite *NativeHTTPClientTestSuite) TestInitialize() {
	suite.client.initialize("key", "secret")

	// checking key/secret
	suite.Equal("key", suite.client.apiKey, "exepcted apiKey: 'key'")
	suite.Equal("secret", suite.client.apiSecret, "exepcted apiKey: 'key'")

	// checking default endpoints
	suite.Equal(GetMediaEndpoint, suite.client.getMediaURL, "Default value for sendMediaURL must be %v", SendMediaEndpoint)
	suite.Equal(SendMediaEndpoint, suite.client.sendMediaURL, "Default value for sendMediaURL must be %v", SendMediaEndpoint)
}

// ************************************************************************************************
// ** checkHTTPClient
// ************************************************************************************************

// TestCheckHTTPClient tests checkHTTPClient
func (suite *NativeHTTPClientTestSuite) TestCheckHTTPClient() {
	// create an instance of BasicClient (without HTTPClient)
	suite.client = &BasicClient{}
	suite.client.initialize(key, secret)

	err := suite.client.checkHTTPClient()

	// checking error
	suite.Error(err, "error expected")
	// error type ==> client.Error
	customError, ok := err.(*Error)
	suite.True(ok, "Error type must be clientv2.Error")
	// custom error type
	suite.Equal(ErrorCodeNotImplemented, customError.Code(), "Custom error code must be %v", ErrorCodeNotImplemented)
}

// ************************************************************************************************
// ** setGetMediaEndpointURL
// ************************************************************************************************

// TestSetGetMediaEndpointURL tests setGetMediaEndpointURL
func (suite *NativeHTTPClientTestSuite) TestSetGetMediaEndpointURL() {
	// updated value
	suite.client.setSendMediaEndpointURL("https://abc.com")
	suite.Equal("https://abc.com", suite.client.sendMediaURL, "incorrect sendMediaURL")
}

// ************************************************************************************************
// ** GetMedia
// ************************************************************************************************

// TestGetMediaWithNoHTTPClient calls GetMedia without a valid HTTPClient (using the BasicClient implementation, which is
// the Abstractor -Bridge pattern-), so it should fail. HTTPClient is the Implementor.
func (suite *NativeHTTPClientTestSuite) TestGetMediaWithNoHTTPClient() {
	suite.client = &BasicClient{}
	suite.client.initialize(key, secret)

	_, _, _, err := suite.client.GetMedia(getMediaUUID)
	suite.Error(err, "GetMedia must return error if no HTTPClient(Implementor) is defined")

	// error type ==> client.Error
	customError, ok := err.(*Error)
	suite.True(ok, "Error type must be client.Error")
	// custom error type
	suite.Equal(ErrorCodeNotImplemented, customError.Code(), "Custom error code must be %v", ErrorCodeNotImplemented)
}

// TestGetMediaWithEmptyString tests GetMedia passing an empty string as mediaUUID
func (suite *NativeHTTPClientTestSuite) TestGetMediaWithEmptyString() {
	_, _, _, err := suite.client.GetMedia("")
	suite.Error(err, "GetMedia must return error if an empty string is provided as mediaUUID")

	// error type ==> client.Error
	customError, ok := err.(*Error)
	suite.True(ok, "Error type must be client.Error")
	// custom error type
	suite.Equal(ErrorCodeArguments, customError.Code(), "Custom error code must be %v", ErrorCodeArguments)
}

// ************************************************************************************************
// ** setSendMediaEndpointURL
// ************************************************************************************************

// TestSetSendMediaEndpointURL tests setSendMediaEndpointURL
func (suite *NativeHTTPClientTestSuite) TestSetSendMediaEndpointURL() {
	// updated value
	suite.client.setSendMediaEndpointURL("https://abc.com")
	suite.Equal("https://abc.com", suite.client.sendMediaURL, "incorrect sendMediaURL")
}

// ************************************************************************************************
// ** SendMedia
// ************************************************************************************************

// TestSendMediaWithNoHTTPClient calls SendMedia without a valid HTTPClient (using the BasicClient implementation, which is
// the Abstractor -Bridge pattern-), so it should fail. HTTPClient is the Implementor.
func (suite *NativeHTTPClientTestSuite) TestSendMediaWithNoHTTPClient() {
	suite.client = &BasicClient{}
	suite.client.initialize(key, secret)

	flags := DetectionFlags{Classifiers: true}
	_, _, _, err := suite.client.SendMedia(sendMediaFileURI, flags, []string{}, sendMediaFilename)
	suite.Error(err, "SendMedia must return error if no HTTPClient(Implementor) is defined")

	// error type ==> client.Error
	customError, ok := err.(*Error)
	suite.True(ok, "Error type must be client.Error")
	// custom error type
	suite.Equal(ErrorCodeNotImplemented, customError.Code(), "Custom error code must be %v", ErrorCodeNotImplemented)
}

// TestSendMediaInvalidEndpointURL tests SendMedia with an invalid endpoint URL
func (suite *NativeHTTPClientTestSuite) TestSendMediaInvalidEndpointURL() {
	// wrong URL can't be parsed and will return error
	suite.client.sendMediaURL = "abc :// domain . com / image . jpg"
	_, _, _, err := suite.client.SendMedia(sendMediaFileURI, DetectionFlags{Classifiers: true}, []string{}, sendMediaFilename)

	// checking error
	suite.Error(err, "SendMedia must return error if fileURI could not be parsed")

	// error type ==> client.Error
	customError, ok := err.(*Error)
	suite.True(ok, "Error type must be client.Error")
	// custom error type
	suite.Equal(ErrorCodeHTTPRequest, customError.Code(), "Custom error code must be %v", ErrorCodeHTTPRequest)
}

// TestSendMediaStatus500 tests SendMedia after a 500 error
func (suite *NativeHTTPClientTestSuite) TestSendMediaStatus500() {
	// get the httpclient
	httpClient := suite.client.getHTTPClient()
	// start the fake endpoint (which mimics the original endpoint)
	httpTestServer := httpClient.NewTestServer(func(w gohttpclient.ResponseWriter, req gohttpclient.Request) {
		w.SetStatusCode(500)
		w.Print(sendMediaPOST500ResponseBody)
	})
	defer httpTestServer.Close()

	// make the request
	suite.client.setSendMediaEndpointURL(httpTestServer.GetURL())
	httpResp, sendMediaResp, errorResp, _ := suite.client.SendMedia("https://testDomain.com/image.jpg", DetectionFlags{}, []string{}, "image.jpg")

	// checking response code
	suite.Equal(500, httpResp.GetStatusCode(), "status code must be 500")
	// checking response body
	suite.Equal(sendMediaPOST500ResponseBody, httpResp.GetBody(), "incorrect response body")
	// checking ErrorResponse
	suite.Equal(ErrorResponse{Code: 0, Description: "string"}, errorResp, "incorrect ErrorResponse")
	// checking SendMediaResponse
	suite.Nil(sendMediaResp, "unexpected SendMediaResponse")
}

// TestSendMediaInvalidJSON tests SendMedia sending back an invalid JSON as the response's body
func (suite *NativeHTTPClientTestSuite) TestSendMediaInvalidJSON() {
	// get the httpclient
	httpClient := suite.client.getHTTPClient()

	// start the fake endpoint (which mimics the original endpoint)
	httpTestServer := httpClient.NewTestServer(func(w gohttpclient.ResponseWriter, req gohttpclient.Request) {
		w.Print(`{{"abc":`)
	})
	defer httpTestServer.Close()

	// make the request
	suite.client.setSendMediaEndpointURL(httpTestServer.GetURL())
	_, _, _, err := suite.client.SendMedia(sendMediaFileURI, DetectionFlags{Classifiers: true}, []string{}, sendMediaFilename)

	// checking error
	suite.Error(err, "error expected")

	// error type ==> clientv2.Error
	customError, ok := err.(*Error)
	suite.True(ok, "Error type must be clientv2.Error")
	// custom error type
	suite.Equal(ErrorCodeJSONUnmarshal, customError.Code(), "Custom error code must be %v", ErrorCodeJSONUnmarshal)
}

// TestSendMediaStatus500 tests SendMedia after a 200 status + valid json response
func (suite *NativeHTTPClientTestSuite) TestSendMediaStatus200() {
	var (
		fileURI          = "https://cdn.vox-cdn.com/thumbor/3V8wxIEwW8-JjMu-dX7lwcVQWd0=/0x0:1000x563/1200x800/filters:focal(420x202:580x362)/cdn.vox-cdn.com/uploads/chorus_image/image/60350569/killing_eve_review.0.jpg"
		originalFilename = "killing_eve_review.0.jpg"
	)

	// get the httpclient
	httpClient := suite.client.getHTTPClient()

	// start the fake endpoint (which mimics the original endpoint)
	httpTestServer := httpClient.NewTestServer(func(w gohttpclient.ResponseWriter, req gohttpclient.Request) {
		w.SetStatusCode(200)
		w.Print(sendMediaPOST200ResponseBody)
	})
	defer httpTestServer.Close()

	// make the request
	suite.client.setSendMediaEndpointURL(httpTestServer.GetURL())
	httpResp, resp, errResp, err := suite.client.SendMedia(fileURI, DetectionFlags{Classifiers: true}, []string{}, originalFilename)

	// checking no error
	suite.NoError(err, "no error expected")

	// checking no ErrorResponse
	suite.Equal(ErrorResponse{}, errResp, "no ErrorResponse expected")

	// checking http response code
	suite.Equal(200, httpResp.GetStatusCode(), "status code must be 200")
	// checking http response body
	suite.Equal(sendMediaPOST200ResponseBody, httpResp.GetBody(), "incorrect response body")

	// checking SendMediaResponse struct values
	suite.Equal(originalFilename, resp.Media.OriginalFilename, "incorrect original filename")
	// checking MediaUUID
	suite.Equal("8380b14d-1ba1-4355-b37c-3525a6f28620", resp.Media.MediaUUID, "incorrect mediaUUID")
	// checking total faces
	suite.Equal(1, len(resp.Media.Faces), "incorrect number of faces")
	// checking FaceUUID
	suite.Equal("3b21a241-50a5-11e9-9287-0cc47a6c4dbd", resp.Media.Faces[0].FaceUUID)
}

// ************************************************************************************************
// ** Run Test Suite
// ************************************************************************************************

func TestClientRunSuite(t *testing.T) {
	suite.Run(t, new(NativeHTTPClientTestSuite))
}
