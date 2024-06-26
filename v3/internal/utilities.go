// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"time"
)

// FloatSecondsToDuration turns a float64 in seconds into a time.Duration.
func FloatSecondsToDuration(seconds float64) time.Duration {
	nanos := seconds * 1000 * 1000 * 1000
	return time.Duration(nanos) * time.Nanosecond
}

// CompactJSONString removes the whitespace from a JSON string.  This function
// will panic if the string provided is not valid JSON.  Thus is must only be
// used in testing code!
func CompactJSONString(js string) string {
	buf := new(bytes.Buffer)
	if err := json.Compact(buf, []byte(js)); err != nil {
		panic(fmt.Errorf("unable to compact JSON: %v", err))
	}
	return buf.String()
}

// HandlerName return name of a function.
func HandlerName(h interface{}) string {
	if h == nil {
		return ""
	}
	t := reflect.ValueOf(h).Type()
	if t.Kind() == reflect.Func {
		if pointer := runtime.FuncForPC(reflect.ValueOf(h).Pointer()); pointer != nil {
			return pointer.Name()
		}
	}
	return ""
}
