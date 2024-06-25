// SPDX-FileCopyrightText: 2022 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package candlelight

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenTID(t *testing.T) {
	assert := assert.New(t)
	tid := GenTID()
	assert.NotEmpty(tid)
}
