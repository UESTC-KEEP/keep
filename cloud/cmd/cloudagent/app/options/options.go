/*
Copyright 2021 The KEEP Authors.

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

package options

import (
	cloud2 "keep/constants/cloud"
	"keep/pkg/apis/compoenentconfig/keep/v1alpha1/cloud"
)

type CloudAgentOptions struct {
	ConfigFile string
}

func NewDefaultEdgeAgentOptions() *CloudAgentOptions {
	return &CloudAgentOptions{
		ConfigFile: cloud2.DefaultCloudConfigFile,
	}
}

func (o *CloudAgentOptions) Config() (*cloud.CloudAgentConfig, error) {
	cfg := cloud.NewDefaultEdgeAgentConfig()
	if err := cfg.Parse(o.ConfigFile); err != nil {
		return nil, err
	}
	return cfg, nil
}
