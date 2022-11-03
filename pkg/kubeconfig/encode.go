/*
Copyright 2019 The Kubernetes Authors.
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

package kubeconfig

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Encode encodes the cfg to yaml
func Encode(cfg *Config) ([]byte, error) {
	// NOTE: kubernetes's yaml library doesn't handle inline fields very well
	// so we're not using that to marshal
	encoded, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode KUBECONFIG: %v", err)
	}

	return encoded, nil
}

func Decode(in []byte) (*Config, error) {
	decoded := &Config{}
	if err := yaml.Unmarshal(in, decoded); err != nil {
		return nil, fmt.Errorf("failed to decode KUBECONFIG: %v", err)
	}
	return decoded, nil
}
