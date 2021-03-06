/*
Copyright 2020 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
Script expects the following env variables:
 - UPDATE_TARGET=<string>: optional - if unset/absent, default option is "fs"; valid options are:
   - "fs"  - update only local filesystem repo files [default]
   - "gh"  - update only remote GitHub repo files and create PR (if one does not exist already)
   - "all" - update local and remote repo files and create PR (if one does not exist already)
 - GITHUB_TOKEN=<string>: GitHub [personal] access token
   - note: GITHUB_TOKEN is required if UPDATE_TARGET is "gh" or "all"
*/

package main

import (
	"context"
	"time"

	"golang.org/x/mod/semver"
	"k8s.io/klog/v2"

	"k8s.io/minikube/hack/update"
)

const (
	// default context timeout
	cxTimeout = 300 * time.Second
)

var (
	schema = map[string]update.Item{
		"Makefile": {
			Replace: map[string]string{
				`GOLINT_VERSION \?= v1.*`: `GOLINT_VERSION ?= {{.StableVersion}}`,
			},
		},
	}

	// PR data
	prBranchPrefix = "update-golint-version_" // will be appended with first 7 characters of the PR commit SHA
	prTitle        = `update go lint version: {stable: "{{.StableVersion}}"}`
)

// Data holds stable gopogh version in semver format.
type Data struct {
	StableVersion string `json:"stableVersion"`
}

func main() {
	// set a context with defined timeout
	ctx, cancel := context.WithTimeout(context.Background(), cxTimeout)
	defer cancel()

	// get Golang stable version
	stable, err := golintVersion(ctx, "golangci", "golangci-lint")
	if err != nil {
		klog.Fatalf("Unable to get Golang stable version: %v", err)
	}
	data := Data{StableVersion: stable}
	klog.Infof("Golang stable version: %s", data.StableVersion)

	update.Apply(ctx, schema, data, prBranchPrefix, prTitle, 12247)
}

//
// golintVersions returns stable version in semver format.
func golintVersion(ctx context.Context, owner, repo string) (stable string, err error) {
	// get Kubernetes versions from GitHub Releases
	stable, _, err = update.GHReleases(ctx, owner, repo)
	if err != nil || !semver.IsValid(stable) {
		return "", err
	}
	return stable, nil
}
