// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api_test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/yanhuangpai/ifi-bee/pkg/api"
	"github.com/yanhuangpai/ifi-bee/pkg/jsonhttp"
	"github.com/yanhuangpai/ifi-bee/pkg/jsonhttp/jsonhttptest"
	"github.com/yanhuangpai/ifi-bee/pkg/logging"
	statestore "github.com/yanhuangpai/ifi-bee/pkg/statestore/mock"
	"github.com/yanhuangpai/ifi-bee/pkg/storage/mock"
	"github.com/yanhuangpai/ifi-bee/pkg/swarm"
	"github.com/yanhuangpai/ifi-bee/pkg/tags"
	"github.com/yanhuangpai/ifi-bee/pkg/traversal"
)

func TestPinBzzHandler(t *testing.T) {
	var (
		dirUploadResource     = "/dirs"
		pinBzzResource        = "/pin/bzz"
		pinBzzAddressResource = func(addr string) string { return pinBzzResource + "/" + addr }
		pinChunksResource     = "/pin/chunks"

		mockStorer       = mock.NewStorer()
		mockStatestore   = statestore.NewStateStore()
		traversalService = traversal.NewService(mockStorer)
		logger           = logging.New(ioutil.Discard, 0)
		client, _, _     = newTestServer(t, testServerOptions{
			Storer:    mockStorer,
			Traversal: traversalService,
			Tags:      tags.NewTags(mockStatestore, logger),
		})
	)

	t.Run("pin-bzz-1", func(t *testing.T) {
		files := []f{
			{
				data: []byte("<h1>Swarm"),
				name: "index.html",
				dir:  "",
			},
		}

		tarReader := tarFiles(t, files)

		rootHash := "a85aaea6a34a5c7127a3546196f2111f866fe369c6d6562ed5d3313a99388c03"

		// verify directory tar upload response
		jsonhttptest.Request(t, client, http.MethodPost, dirUploadResource, http.StatusOK,
			jsonhttptest.WithRequestBody(tarReader),
			jsonhttptest.WithRequestHeader("Content-Type", api.ContentTypeTar),
			jsonhttptest.WithExpectedJSONResponse(api.FileUploadResponse{
				Reference: swarm.MustParseHexAddress(rootHash),
			}),
		)

		jsonhttptest.Request(t, client, http.MethodPost, pinBzzAddressResource(rootHash), http.StatusOK,
			jsonhttptest.WithExpectedJSONResponse(jsonhttp.StatusResponse{
				Message: http.StatusText(http.StatusOK),
				Code:    http.StatusOK,
			}),
		)

		expectedChunkCount := 7

		// get the reference as everytime it will change because of random encryption key
		var resp api.ListPinnedChunksResponse

		jsonhttptest.Request(t, client, http.MethodGet, pinChunksResource, http.StatusOK,
			jsonhttptest.WithUnmarshalJSONResponse(&resp),
		)

		if expectedChunkCount != len(resp.Chunks) {
			t.Fatalf("expected to find %d pinned chunks, got %d", expectedChunkCount, len(resp.Chunks))
		}
	})

	t.Run("unpin-bzz-1", func(t *testing.T) {
		rootHash := "a85aaea6a34a5c7127a3546196f2111f866fe369c6d6562ed5d3313a99388c03"

		jsonhttptest.Request(t, client, http.MethodDelete, pinBzzAddressResource(rootHash), http.StatusOK,
			jsonhttptest.WithExpectedJSONResponse(jsonhttp.StatusResponse{
				Message: http.StatusText(http.StatusOK),
				Code:    http.StatusOK,
			}),
		)

		jsonhttptest.Request(t, client, http.MethodGet, pinChunksResource, http.StatusOK,
			jsonhttptest.WithExpectedJSONResponse(api.ListPinnedChunksResponse{
				Chunks: []api.PinnedChunk{},
			}),
		)
	})

}
