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
	assert.Equal(t, "0.0.0.0", result)

	result = NormalizeVersion("1.2.3.4")
	assert.Equal(t, "1.2.3.4", result)

	result = NormalizeVersion("01.02.03.04")
	assert.Equal(t, "1.2.3.4", result)

	result = NormalizeVersion("010.021.030.040")
	assert.Equal(t, "10.21.30.40", result)

	result = NormalizeVersion("23.01+dfsg-6")
	assert.Equal(t, "23.1.0.6", result)

	result = NormalizeVersion("6.10.18.1793")
	assert.Equal(t, "6.10.18.1793", result)

	result = NormalizeVersion("0.5.2+dfsg+~cs5.2.9-9")
	assert.Equal(t, "0.5.2.9", result)

	result = NormalizeVersion("2.41-6~deepin1+11+nmu1")
	assert.Equal(t, "2.41.0.6", result)

	result = NormalizeVersion("4:22.12.3-1")
	assert.Equal(t, "22.12.3.1", result)

	result = NormalizeVersion("4:22.12.3-1")
	assert.Equal(t, "22.12.3.1", result)

	result = NormalizeVersion("1:20230101~dfsg-3")
	assert.Equal(t, "20230101.0.0.3", result)

	result = NormalizeVersion("1:5.9~svn20110310-15")
	assert.Equal(t, "5.9.0.15", result)

	result = NormalizeVersion("3:6.04~git20190206.bf6db5b4+dfsg1-deepin1")
	assert.Equal(t, "6.4.6.1", result)

	result = NormalizeVersion("1:0.2.91~git20170110-4")
	assert.Equal(t, "0.2.91.4", result)

	result = NormalizeVersion("1:0.2.91~git20170110-4")
	assert.Equal(t, "0.2.91.4", result)

	result = NormalizeVersion("1:4.00~git30-7274cfa-1.1")
	assert.Equal(t, "4.0.0.1", result)

	result = NormalizeVersion("1:1.3.dfsg+really1.3.1-1deepin1")
	assert.Equal(t, "1.3.1.1", result)
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
