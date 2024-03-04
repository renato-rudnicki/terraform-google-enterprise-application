// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bootstrap

import (
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/gcloud"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/tft"
	"github.com/stretchr/testify/assert"
)

func TestBootstrap(t *testing.T) {

	bootstrap := tft.NewTFBlueprintTest(t,
		tft.WithTFDir("../../../1-bootstrap"),
	)

	bootstrap.DefineVerify(func(assert *assert.Assertions) {
		bootstrap.DefaultVerify(assert)

		projectID := bootstrap.GetStringOutput("project_id")
		gcloudArgsBucket := gcloud.WithCommonArgs([]string{"--project", projectID, "--json"})

		bucketPrefix := []string{
			"mt-state",
			"fs-state",
			"af-state",
		}
		for _, bucket := range bucketPrefix {
			urlBucket := fmt.Sprintf("https://www.googleapis.com/storage/v1/b/%s-%s", bucket, projectID)
			opBucket := gcloud.Run(t, fmt.Sprintf("storage ls --buckets gs://%s-%s", bucket, projectID), gcloudArgsBucket).Array()
			//t.Log(opBucket)
			assert.Equal(urlBucket, opBucket[0].Get("metadata.selfLink").String(), fmt.Sprintf("The bucket name should be %s.", urlBucket))
			assert.True(opBucket[0].Exists(), "Bucket %s should exist.", urlBucket)
		}

		repos := []string{
			"eab-applicationfactory",
			"eab-fleetscope",
			"eab-multitenant",
		}
		for _, repo := range repos {
			url := fmt.Sprintf("https://source.developers.google.com/p/%s/r/%s", projectID, repo)
			repoOP := gcloud.Runf(t, "source repos describe %s --project %s", repo, projectID)
			assert.Equal(url, repoOP.Get("url").String(), "source repo %s should have url %s", repo, url)
			t.Log(repoOP)
		}
	})
	bootstrap.Test()
}

