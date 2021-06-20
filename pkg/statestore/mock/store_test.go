// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.package storage

package mock_test

import (
	"testing"

	"github.com/yanhuangpai/ifi-bee/pkg/statestore/mock"
	"github.com/yanhuangpai/ifi-bee/pkg/statestore/test"
	"github.com/yanhuangpai/ifi-bee/pkg/storage"
)

func TestMockStateStore(t *testing.T) {
	test.Run(t, func(t *testing.T) storage.StateStorer {
		return mock.NewStateStore()
	})
}
