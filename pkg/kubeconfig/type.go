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

type Config struct {
	// Clusters is a map of referenceable names to cluster configs
	Clusters []NamedCluster `yaml:"clusters,omitempty"`
	// Users is a map of referenceable names to user configs
	Users []NamedUser `yaml:"users,omitempty"`
	// Contexts is a map of referenceable names to context configs
	Contexts []NamedContext `yaml:"contexts,omitempty"`
	// CurrentContext is the name of the context that you would like to use by default
	CurrentContext string `yaml:"current-context,omitempty"`
	// OtherFields contains fields kind does not inspect or modify, these are
	// read purely for writing back
	OtherFields map[string]interface{} `yaml:",inline,omitempty"`
}

type Cluster struct {
	// Server is the address of the kubernetes cluster (https://hostname:port).
	Server string `yaml:"server,omitempty"`
	// OtherFields contains fields kind does not inspect or modify, these are
	// read purely for writing back
	OtherFields map[string]interface{} `yaml:",inline,omitempty"`
}

type NamedCluster struct {
	// Name is the nickname for this Cluster
	Name string `yaml:"name"`
	// Cluster holds the cluster information
	Cluster Cluster `yaml:"cluster"`
}

type NamedUser struct {
	// Name is the nickname for this User
	Name string `yaml:"name"`
	// User holds the user information
	// We do not touch this and merely write it back
	User map[string]interface{} `yaml:"user"`
}

type NamedContext struct {
	// Name is the nickname for this Context
	Name string `yaml:"name"`
	// Context holds the context information
	Context Context `yaml:"context"`
}

type Context struct {
	// Cluster is the name of the cluster for this context
	Cluster string `yaml:"cluster"`
	// User is the name of the User for this context
	User string `yaml:"user"`
	// OtherFields contains fields kind does not inspect or modify, these are
	// read purely for writing back
	OtherFields map[string]interface{} `yaml:",inline,omitempty"`
}
