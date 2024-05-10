package middleware

import (
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestThrottle_CallsNextMiddleware(t *testing.T) {
	assert := assert.New(t)
	m := mockEchoContext{}
	next, called := createHandlerFuncWithCalledBoolean()

	em, close := Throttle(1, 1, 1)
	callable := em(next)
	callable(&m)

	assert.True(*called)
	close <- true
}

func TestThrottle_WhenNoTokensLeft_ExpectTooManyRequests(t *testing.T) {
	assert := assert.New(t)
	mc := newMockEchoContext(http.StatusOK)
	mc.request.Header = map[string][]string{
		apiKeyHeaderKey: {defaultApiKey1.String()},
	}
	next := createHandlerFuncReturning(nil)

	em, close := Throttle(0, 0, 0)
	callable := em(next)
	callable(mc)

	assert.Equal(http.StatusTooManyRequests, mc.reportedCode)
	close <- true
}

func TestThrottle_WhenWaitingForRefill_ExpectOk(t *testing.T) {
	assert := assert.New(t)
	mc := newMockEchoContext(http.StatusOK)
	mc.request.Header = map[string][]string{
		apiKeyHeaderKey: {defaultApiKey1.String()},
	}
	next := createHandlerFuncReturningCode(http.StatusOK)

	em, close := Throttle(0, 2, 2)
	callable := em(next)

	callable(mc)
	assert.Equal(http.StatusTooManyRequests, mc.reportedCode)

	time.Sleep(2 * time.Second)

	callable(mc)
	assert.Equal(http.StatusOK, mc.reportedCode)

	close <- true
}

func TestThrottle_PropagatesError(t *testing.T) {
	assert := assert.New(t)
	m := newMockEchoContext(http.StatusOK)
	next := createHandlerFuncReturning(errDefault)

	em, close := Throttle(1, 1, 1)

	callable := em(next)
	actual := callable(m)

	assert.Equal(errDefault, actual)
	close <- true
}

func TestThrottle_ConcurrentUse_ExpectFirstServed(t *testing.T) {
	assert := assert.New(t)

	em, close := Throttle(1, 1, 1)
	next := createHandlerFuncReturningCode(http.StatusOK)
	handler := em(next)

	client1 := newMockEchoContext(http.StatusOK)
	client2 := newMockEchoContext(http.StatusOK)

	var wg sync.WaitGroup
	wg.Add(2)

	c1 := func() {
		defer wg.Done()

		handler(client1)
		assert.Equal(http.StatusOK, client1.reportedCode)

		time.Sleep(1500 * time.Millisecond)

		handler(client1)
		assert.Equal(http.StatusTooManyRequests, client1.reportedCode)
	}

	c2 := func() {
		defer wg.Done()

		time.Sleep(1100 * time.Millisecond)

		handler(client2)
		assert.Equal(http.StatusOK, client2.reportedCode)
	}

	go c1()
	go c2()

	wg.Wait()

	close <- true
}
