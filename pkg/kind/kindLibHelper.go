/*
Copyright 2022 Manuel Bovo.

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

package kind

import (
	"context"
	"strconv"
	"time"

	"net/url"

	"github.com/go-logr/logr"
	"github.com/mbovo/cluster-api-provider-kind/api/v1beta1"
	"github.com/mbovo/cluster-api-provider-kind/pkg/kubeconfig"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/log"
	v1alpha4Kind "sigs.k8s.io/kind/pkg/apis/config/v1alpha4"
	kindApiCluster "sigs.k8s.io/kind/pkg/cluster"
)

type KindLibHelper struct {
	Provider *kindApiCluster.Provider
	Config   *v1alpha4Kind.Cluster
}

func NewKindLibHelper(inlog logr.Logger, kindCluster *v1beta1.KindCluster, capiCluster *clusterv1.Cluster) *KindLibHelper {

	// This LoggerWrapper wraps logr.Logger in a Kind Logger fashion to make kind log to the controller logs
	logger := NewLoggerWrapper()
	provider := kindApiCluster.NewProvider(kindApiCluster.ProviderWithLogger(logger), kindApiCluster.ProviderWithDocker())
	config := newClusterConfig(kindCluster, capiCluster)
	return &KindLibHelper{Provider: provider, Config: config}
}

func (k *KindLibHelper) Exists(ctx context.Context, kindCluster *v1beta1.KindCluster) (bool, error) {
	clusterName := kindCluster.Name
	clusters, err := k.Provider.List()
	if err != nil {
		return false, nil
	}
	for _, cluster := range clusters {
		if cluster == clusterName {
			return true, nil
		}
	}
	return false, nil
}

func (k *KindLibHelper) Endpoint(ctx context.Context, kindCluster *v1beta1.KindCluster) (host string, port int, err error) {
	logger := log.FromContext(ctx)
	logger.Info("Getting Kind cluster endpoint", "cluster", kindCluster.Name)
	str, err := k.Provider.KubeConfig(kindCluster.Name, false)
	kubeCfg, err := kubeconfig.Decode([]byte(str))

	if kubeCfg != nil {
		url, err := url.Parse(kubeCfg.Clusters[0].Cluster.Server)
		if err != nil {
			return "", 0, err
		}
		host = url.Hostname()
		port, _ = strconv.Atoi(url.Port())
	}

	return
}

func (k *KindLibHelper) Create(ctx context.Context, kindCluster *v1beta1.KindCluster) (err error) {
	logger := log.FromContext(ctx)
	logger.Info("Creating Kind cluster", "cluster", kindCluster.Name)

	clusterCfg := kindApiCluster.CreateWithV1Alpha4Config(k.Config)
	waitForReady := kindApiCluster.CreateWithWaitForReady(120 * time.Second)
	kubeConfigPath := kindApiCluster.CreateWithKubeconfigPath("/tmp/kind-config-" + kindCluster.Name)

	err = k.Provider.Create(kindCluster.Name, clusterCfg, waitForReady, kubeConfigPath)

	return
}

func (k *KindLibHelper) Delete(ctx context.Context, kindCluster *v1beta1.KindCluster) (err error) {
	logger := log.FromContext(ctx)
	logger.Info("Deleting Kind cluster", "cluster", kindCluster.Name)

	err = k.Provider.Delete(kindCluster.Name, "/tmp/kind-config-"+kindCluster.Name)

	return
}

func newClusterConfig(kindCluster *v1beta1.KindCluster, capiCluster *clusterv1.Cluster) *v1alpha4Kind.Cluster {

	cfg := &v1alpha4Kind.Cluster{}
	cfg.Name = kindCluster.Name

	// adding nodes to the kind cluster
	for i := 0; i < int(kindCluster.Spec.ControlPlaneCount); i++ {
		cfg.Nodes = append(cfg.Nodes, v1alpha4Kind.Node{Role: v1alpha4Kind.ControlPlaneRole, Image: kindCluster.Spec.Image})
	}
	for i := 0; i < int(kindCluster.Spec.WorkerCount); i++ {
		cfg.Nodes = append(cfg.Nodes, v1alpha4Kind.Node{Role: v1alpha4Kind.WorkerRole, Image: kindCluster.Spec.Image})
	}

	return cfg
}
