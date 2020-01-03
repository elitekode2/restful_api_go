package log

import (
	"bytes"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	assert.NotNil(t, New())
}

func TestNewWithZap(t *testing.T) {
	zl, _ := zap.NewProduction()
	l := NewWithZap(zl)
	assert.NotNil(t, l)
}

func TestWithRequest(t *testing.T) {
	req := buildRequest("abc", "123")
	ctx := WithRequest(context.Background(), req)
	assert.Equal(t, "abc", ctx.Value(requestIDKey).(string))
	assert.Equal(t, "123", ctx.Value(correlationIDKey).(string))

	req = buildRequest("", "123")
	ctx = WithRequest(context.Background(), req)
	assert.NotEmpty(t, ctx.Value(requestIDKey).(string))
	assert.Equal(t, "123", ctx.Value(correlationIDKey).(string))
}

func Test_getCorrelationID(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", bytes.NewBufferString(""))
	assert.Empty(t, getCorrelationID(req))
	req.Header.Set("X-Correlation-ID", "test")
	assert.Equal(t, "test", getCorrelationID(req))
}

func Test_getRequestID(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", bytes.NewBufferString(""))
	assert.Empty(t, getRequestID(req))
	req.Header.Set("X-Request-ID", "test")
	assert.Equal(t, "test", getRequestID(req))
}

func Test_logger_With(t *testing.T) {
	l := New()
	l2 := l.With(nil)
	assert.True(t, reflect.DeepEqual(l2, l))

	req := buildRequest("abc", "123")
	ctx := WithRequest(context.Background(), req)
	l3 := l.With(ctx)
	assert.False(t, reflect.DeepEqual(l3, l2))
}

func buildRequest(requestID, correlationID string) *http.Request {
	req, _ := http.NewRequest("GET", "http://example.com", bytes.NewBufferString(""))
	if requestID != "" {
		req.Header.Set("X-Request-ID", requestID)
	}
	if correlationID != "" {
		req.Header.Set("X-Correlation-ID", correlationID)
	}
	return req
}

func TestNewForTest(t *testing.T) {
	logger, entries := NewForTest()
	assert.Equal(t, 0, entries.Len())
	logger.Info("msg 1")
	assert.Equal(t, 1, entries.Len())
	logger.Info("msg 2")
	logger.Info("msg 3")
	assert.Equal(t, 3, entries.Len())
	entries.TakeAll()
	assert.Equal(t, 0, entries.Len())
	logger.Info("msg 4")
	assert.Equal(t, 1, entries.Len())
}

func TestDBQuery(t *testing.T) {
	logger, entries := NewForTest()
	f := DBQuery(logger)
	f(context.Background(), time.Millisecond*3, "sql", nil, nil)
	if assert.Equal(t, 1, entries.Len()) {
		assert.Equal(t, "DB query successful", entries.All()[0].Message)
	}
	entries.TakeAll()

	f(context.Background(), time.Millisecond*3, "sql", nil, fmt.Errorf("test"))
	if assert.Equal(t, 1, entries.Len()) {
		assert.Equal(t, "DB query error: test", entries.All()[0].Message)
	}
}

func TestDBExec(t *testing.T) {
	logger, entries := NewForTest()
	f := DBExec(logger)
	f(context.Background(), time.Millisecond*3, "sql", nil, nil)
	if assert.Equal(t, 1, entries.Len()) {
		assert.Equal(t, "DB execution successful", entries.All()[0].Message)
	}
	entries.TakeAll()

	f(context.Background(), time.Millisecond*3, "sql", nil, fmt.Errorf("test"))
	if assert.Equal(t, 1, entries.Len()) {
		assert.Equal(t, "DB execution error: test", entries.All()[0].Message)
	}
}
