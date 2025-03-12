/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package layer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeVersion(t *testing.T) {
	result := NormalizeVersion("")
	assert.Equal(t, result, "0.0.0.0")

	result = NormalizeVersion("1.2.3.4")
	assert.Equal(t, result, "1.2.3.4")

	result = NormalizeVersion("01.02.03.04")
	assert.Equal(t, result, "1.2.3.4")

	result = NormalizeVersion("010.021.030.040")
	assert.Equal(t, result, "10.21.30.40")

	result = NormalizeVersion("010.02-1.0X30.0~40")
	assert.Equal(t, result, "10.21.30.40")

	result = NormalizeVersion("123.456-fwf78.12")
	assert.Equal(t, result, "123.45678.12.0")

	result = NormalizeVersion("123.456-fwf78.012")
	assert.Equal(t, result, "123.45678.12.0")
}

func TestIsValidVersion(t *testing.T) {
	result := IsValidComponentVersion("")
	assert.False(t, result)

	result = IsValidComponentVersion("1.2.3.4")
	assert.True(t, result)

	result = IsValidComponentVersion("01.02.3.4")
	assert.False(t, result)

	result = IsValidComponentVersion("1.02.3.4")
	assert.False(t, result)

	result = IsValidComponentVersion("1.2.03.4")
	assert.False(t, result)

	result = IsValidComponentVersion("1.2.3.04")
	assert.False(t, result)

	result = IsValidComponentVersion("-1.02.3.4")
	assert.False(t, result)

	result = IsValidComponentVersion("1.-2.3.4")
	assert.False(t, result)

	result = IsValidComponentVersion("1.2.-3.4")
	assert.False(t, result)

	result = IsValidComponentVersion("1.2.3.-4")
	assert.False(t, result)

	result = IsValidComponentVersion("010.021.030.040")
	assert.False(t, result)

	result = IsValidComponentVersion("010.02-1.0X30.0~40")
	assert.False(t, result)

	result = IsValidComponentVersion("123.456-fwf78.12")
	assert.False(t, result)

	result = IsValidComponentVersion("123.456-fwf78.012")
	assert.False(t, result)
}
