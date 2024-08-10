package middleware

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var defaultUuid = uuid.MustParse("08ce96a3-3430-48a8-a3b2-b1c987a207ca")

func TestResponseEnvelope_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	r := responseEnvelope{
		RequestId: defaultUuid,
		Status:    "SUCCESS",
		Details:   json.RawMessage([]byte(`{"Field":32}`)),
	}

	out, err := json.Marshal(r)

	assert.Nil(err)
	expectedJson := `
	{
		"requestId": "08ce96a3-3430-48a8-a3b2-b1c987a207ca",
		"status": "SUCCESS",
		"details": {
			"Field": 32
		}
	}`
	assert.JSONEq(expectedJson, string(out))
}

func TestEnvelopeResponseWriter_UsesProvidedWriter(t *testing.T) {
	assert := assert.New(t)
	m := &mockResponseWriter{}

	erw := new(m, defaultUuid, nil)

	assert.Equal(m, erw.writer)
}

func TestEnvelopeResponseWriter_AutomaticallySetsSuccessStatusWhenNoStatusIsUsed(t *testing.T) {
	assert := assert.New(t)
	m := &mockResponseWriter{}

	erw := new(m, defaultUuid, nil)
	erw.Write(sampleJsonData)

	expectedJson := `
	{
		"requestId": "08ce96a3-3430-48a8-a3b2-b1c987a207ca",
		"status": "SUCCESS",
		"details": {
			"value": 12
		}
	}`
	assert.JSONEq(expectedJson, string(m.data))
}

func TestEnvelopeResponseWriter_UsesProvidedRequestId(t *testing.T) {
	assert := assert.New(t)
	m := &mockResponseWriter{}

	erw := new(m, defaultUuid, nil)
	erw.Write(sampleJsonData)

	expectedJson := `
	{
		"requestId": "08ce96a3-3430-48a8-a3b2-b1c987a207ca",
		"status": "SUCCESS",
		"details": {
			"value": 12
		}
	}`
	assert.JSONEq(expectedJson, string(m.data))
}

func TestEnvelopeResponseWriter_ForwardsProvidedWriterHeaders(t *testing.T) {
	assert := assert.New(t)
	m := &mockResponseWriter{
		header: http.Header{
			"key": []string{"val1", "val2"},
		},
	}

	erw := new(m, defaultUuid, nil)

	actual := erw.Header()

	assert.Equal(1, m.headerCalled)
	assert.Equal(m.header, actual)
}

func TestEnvelopeResponseWriter_WriteHeaderForwardsCallToProvidedWriter(t *testing.T) {
	assert := assert.New(t)
	m := &mockResponseWriter{}

	erw := new(m, defaultUuid, nil)

	erw.WriteHeader(http.StatusUnauthorized)

	assert.Equal(1, m.writeHeaderCalled)
	assert.Equal(http.StatusUnauthorized, m.code)
}

func TestEnvelopeResponseWriter_WriteHeaderUpdatesResponseEnvelopeStatus(t *testing.T) {
	assert := assert.New(t)
	m := &mockResponseWriter{}

	erw := new(m, defaultUuid, nil)

	erw.WriteHeader(http.StatusUnauthorized)
	assert.Equal("ERROR", erw.response.Status)

	erw.WriteHeader(http.StatusAccepted)
	assert.Equal("SUCCESS", erw.response.Status)
}

var sampleJsonData = []byte(`{"value":12}`)

func TestEnvelopeResponseWriter_WriteForwardsCallToProvidedWriter(t *testing.T) {
	assert := assert.New(t)
	m := &mockResponseWriter{}

	erw := new(m, defaultUuid, nil)

	erw.Write(sampleJsonData)

	assert.Equal(1, m.writeCalled)
}

var matcherStr = `{"requestId":"08ce96a3-3430-48a8-a3b2-b1c987a207ca","status":"${STATUS}","details":{"value":12}}`
var pattern = `${STATUS}`

func TestEnvelopeResponseWriter_WriteWrapsSuccessDataWithEnvelope(t *testing.T) {
	assert := assert.New(t)
	m := &mockResponseWriter{}

	erw := new(m, defaultUuid, nil)

	erw.WriteHeader(http.StatusAccepted)
	erw.Write(sampleJsonData)
	matcher := regexp.MustCompile(strings.ReplaceAll(matcherStr, pattern, "SUCCESS"))
	assert.True(matcher.MatchString(string(m.data)))
}

func TestEnvelopeResponseWriter_WriteWrapsFailureDataWithEnvelope(t *testing.T) {
	assert := assert.New(t)
	m := &mockResponseWriter{}

	erw := new(m, defaultUuid, nil)

	erw.WriteHeader(http.StatusUnprocessableEntity)
	erw.Write(sampleJsonData)
	matcher := regexp.MustCompile(strings.ReplaceAll(matcherStr, pattern, "ERROR"))
	assert.True(matcher.MatchString(string(m.data)))
}

func TestEnvelopeResponseWriter_UsesProvidedDataWhenNotJsonData(t *testing.T) {
	assert := assert.New(t)
	m := &mockResponseWriter{}
	data := []byte("some-data")

	erw := new(m, defaultUuid, nil)

	erw.Write(data)

	assert.Equal(data, m.data)
}
