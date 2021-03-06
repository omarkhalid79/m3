// Copyright (c) 2018 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTimeString(t *testing.T) {
	parsedTime, err := ParseTimeString("703354793")
	require.NoError(t, err)

	equalTimes := parsedTime.Equal(time.Unix(703354793, 0))
	assert.True(t, equalTimes)
}

func TestParseTimeStringWithMinMaxGoTime(t *testing.T) {
	parsedTime, err := ParseTimeString(minTimeFormatted)
	require.NoError(t, err)
	require.True(t, parsedTime.Equal(time.Unix(0, 0)))

	parsedTime, err = ParseTimeString(maxTimeFormatted)
	require.NoError(t, err)
	require.True(t, time.Now().Sub(parsedTime) < time.Minute)
}

func TestParseTimeStringLargeFloat(t *testing.T) {
	_, err := ParseTimeString("9999999999999.99999")
	require.NoError(t, err)
}

func TestParseTimeStringWithDefault(t *testing.T) {
	defaultTime := time.Now().Add(-1 * time.Minute)
	parsedTime, err := ParseTimeStringWithDefault("", defaultTime)
	require.NoError(t, err)
	require.True(t, defaultTime.Equal(parsedTime))
}

func TestParseDurationStringFloat(t *testing.T) {
	d, err := ParseDurationString("1595545968.4985256")
	require.NoError(t, err)
	v := time.Duration(1595545968.4985256 * float64(time.Second))
	require.Equal(t, v, d)
}

func TestParseDurationStringExtendedDurationString(t *testing.T) {
	d, err := ParseDurationString("2d")
	require.NoError(t, err)
	require.Equal(t, 2*24*time.Hour, d)
}
