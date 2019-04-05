package gobetafaceapi

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/enriquebris/gohttpclient"
)

// ************************************************************************************************
// ** Betaface interface
// ************************************************************************************************

// Betafaceapi methods
type Betafaceapi interface {
	// HTTPClient methods
	getHTTPClient() gohttpclient.HTTPClient
	setHTTPClient(httpClient gohttpclient.HTTPClient)

	// GetMedia methods
	setGetMediaEndpointURL(url string)
	GetMedia(mediaUUID string) (gohttpclient.HTTPResponse, *Media, ErrorResponse, error)
	// SendMedia methods
	setSendMediaEndpointURL(url string)
	SendMedia(fileURI string, flags DetectionFlags, recognizeTargets []string, originalFilename string) (gohttpclient.HTTPResponse, *ExtendedMedia, ErrorResponse, error)
	//SendPerson(facesUUIDs []string, personID string)
}

// ************************************************************************************************
// ** NativeHTTPClient == (Betafaceapi + NativeHTTPClient) implementation
// ************************************************************************************************

// NewNativeHTTPClient builds and returns a *BasicClient using a NewNativeHTTPClient as HTTPClient (Implementor, Bridge pattern)
func NewNativeHTTPClient(apiKey string, apiSecret string) *BasicClient {
	ret := &BasicClient{}
	ret.initialize(apiKey, apiSecret)
	ret.setHTTPClient(gohttpclient.NewNativeHTTPClient())

	return ret
}

// ************************************************************************************************
// ** BasicClient == Basic Betafaceapi client. This class implements Betafaceapi interface.
// ************************************************************************************************

// BasicClient is a basic implementation of Betafaceapi interface.
// It does not contain a HTTPClient (Implementor for the Bridge pattern).
// It is the Abstractor (Bridge pattern)
type BasicClient struct {
	apiKey     string
	apiSecret  string
	httpClient gohttpclient.HTTPClient
	// endpoints URL
	getMediaURL  string
	sendMediaURL string
}

// initialize initializes the BasicClient class
func (st *BasicClient) initialize(apiKey string, apiSecret string) {
	st.apiKey = apiKey
	st.apiSecret = apiSecret

	// endpoints
	st.setGetMediaEndpointURL(GetMediaEndpoint)
	st.setSendMediaEndpointURL(SendMediaEndpoint)
}

// checkHTTPClient checks whether the HTTPClient was provided and return error if it was not
func (st *BasicClient) checkHTTPClient() error {
	if st.httpClient == nil {
		return NewError(ErrorCodeNotImplemented, "Missing Implementor (HTTPClient) to complete the Bridge pattern")
	}
	return nil
}

// getHTTPClient returns the HTTPClient (the Implementor for the Bridge pattern)
func (st *BasicClient) getHTTPClient() gohttpclient.HTTPClient {
	return st.httpClient
}

// setHTTPClient sets the HTTPClient (the Implementor for the Bridge pattern)
func (st *BasicClient) setHTTPClient(httpClient gohttpclient.HTTPClient) {
	st.httpClient = httpClient
}

// setGetMediaEndpointURL sets the GetMedia endpoint URL
func (st *BasicClient) setGetMediaEndpointURL(url string) {
	st.getMediaURL = url
}

// GetMedia gets and returns media information
// Possible responses:
//	200 ==> the media resource was successfully retrieved. Media data comes into *MediaResponse
//	400 ==> bad request
//	404 ==> the media does not exist
//	500 ==> server error. Error info comes into ErrorResponse
//	For all above cases, the HTTP response comes into HTTPResponse
//
// Endpoint: [GET] /v2/media
func (st *BasicClient) GetMedia(mediaUUID string) (gohttpclient.HTTPResponse, *Media, ErrorResponse, error) {
	// verifies that this class contains an implementor (HTTPClient)
	if err := st.checkHTTPClient(); err != nil {
		return nil, nil, ErrorResponse{}, err
	}

	// verifies that the mediaUUID is not an empty string
	if strings.TrimSpace(mediaUUID) == "" {
		return nil, nil, ErrorResponse{}, NewError(ErrorCodeArguments, "mediaUUID can't be an empty string")
	}

	// querystring
	url := fmt.Sprintf(GetMediaEndpoint, st.apiKey, mediaUUID)

	// build the http request
	st.getHTTPClient().Reset()
	st.getHTTPClient().SetMethod(GetMediaMethod)
	st.getHTTPClient().SetURL(url)
	st.getHTTPClient().AddHeader("Content-Type", "application/json")

	// send the request
	resp, err := st.getHTTPClient().Do()
	if err != nil {
		return nil, nil, ErrorResponse{}, NewError(ErrorCodeHTTPRequest, err.Error())
	}

	// decode json based on the http response status
	var (
		mediaResp    Media
		ptrMediaResp *Media = nil
		mediaError   ErrorResponse
	)

	switch resp.GetStatusCode() {
	case 200:
		err = json.Unmarshal([]byte(resp.GetBody()), &mediaResp)
		ptrMediaResp = &mediaResp
	case 400:
	case 404:
	case 500:
		err = json.Unmarshal([]byte(resp.GetBody()), &mediaError)
	default:
	}
	// verify possible unmarshal error
	if err != nil {
		return resp, nil, ErrorResponse{}, NewError(ErrorCodeJSONUnmarshal, err.Error())
	}

	return resp, ptrMediaResp, mediaError, nil
}

// setSendMediaEndpointURL sets the SendMedia endpoint URL
func (st *BasicClient) setSendMediaEndpointURL(url string) {
	st.sendMediaURL = url
}

// SendMedia uploads media file using file URI or file content as BASE64 encoded string.
// Note: all recognizeTargets have to be previously registered, otherwise the API would return error 500.
// Endpoint: [POST] /v2/media
func (st *BasicClient) SendMedia(fileURI string, flags DetectionFlags, recognizeTargets []string, originalFilename string) (gohttpclient.HTTPResponse, *ExtendedMedia, ErrorResponse, error) {
	// verifies that this class contains an implementor (HTTPClient)
	if err := st.checkHTTPClient(); err != nil {
		return nil, nil, ErrorResponse{}, err
	}

	// build the payload object
	payload := SendMediaRequest{
		ApiKey:           st.apiKey,
		FileURI:          fileURI,
		OriginalFilename: originalFilename,
		DetectionFlags:   flags.String(),
		RecognizeTargets: recognizeTargets,
	}

	// json payload
	payloadJSON, err := json.Marshal(&payload)
	if err != nil {
		return nil, nil, ErrorResponse{}, NewError(ErrorCodeJSONMarshal, err.Error())
	}

	// build the http request
	st.getHTTPClient().Reset()
	st.getHTTPClient().SetMethod(SendMediaMethod)
	st.getHTTPClient().SetURL(st.sendMediaURL)
	st.getHTTPClient().AddHeader("Content-Type", "application/json")
	st.getHTTPClient().SetPayload(string(payloadJSON))

	// send the request
	resp, err := st.getHTTPClient().Do()
	if err != nil {
		return nil, nil, ErrorResponse{}, NewError(ErrorCodeHTTPRequest, err.Error())
	}

	// decode json based on the http response status
	var (
		sendMediaResp    ExtendedMedia
		ptrSendMediaResp *ExtendedMedia
		sendMediaError   ErrorResponse
		tmp              interface{}
	)
	// decide whether unmarshal SendMediaResponse / ErrorResponse
	if resp.GetStatusCode() == 200 {
		tmp = &sendMediaResp
		ptrSendMediaResp = &sendMediaResp
	} else {
		tmp = &sendMediaError
		ptrSendMediaResp = nil
	}

	body := resp.GetBody()
	// json to struct
	if err := json.Unmarshal([]byte(body), tmp); err != nil {
		return resp, nil, ErrorResponse{}, NewError(ErrorCodeJSONUnmarshal, err.Error())
	}

	return resp, ptrSendMediaResp, sendMediaError, nil
}
