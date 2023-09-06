package onecache

import (
	"runtime"
	"testing"
)

// A simple assertion library for convenient tests.

func assert(t *testing.T, cond bool, format string, args ...interface{}) {
	if !cond {
		fail(t, format, args...)
	}
}

func assertNot(t *testing.T, cond bool, format string, args ...interface{}) {
	assert(t, !cond, format, args...)
}

func assertEquals(t *testing.T, actual interface{}, expected interface{}) {
	if actual != expected {
		fail(t, "Expected %#v to equal %#v", actual, expected)
	}
}

func assertNotEquals(t *testing.T, actual interface{}, expected interface{}) {
	if actual == expected {
		fail(t, "Expected %#v to not equal %#v", actual, expected)
	}
}

func assertNil(t *testing.T, val interface{}) {
	if val != nil {
		fail(t, "Expected nil, got: %#v", val)
	}
}

func assertNotNil(t *testing.T, val interface{}) {
	if val == nil {
		fail(t, "Expected not nil, got nil")
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		fail(t, "Error: %s", err.Error())
	}
}

func fail(t *testing.T, format string, args ...interface{}) {
	stack := getStackTrace(t)
	t.Fatalf(format+"\n\n"+stack+"\n", args...)
}

func getStackTrace(t *testing.T) string {
	buf := make([]byte, 10000)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}
