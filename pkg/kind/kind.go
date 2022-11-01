package kind

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func ClusterExists(ctx context.Context, clusterName string) (bool, error) {
	logger := log.FromContext(ctx)
	b, err := exec.Command("kind", "get", "clusters").Output()
	if err != nil {
		return false, errors.Wrap(err, "failed to get kind clusters")
	}
	for _, line := range strings.Split(string(b), "\n") {
		if line == clusterName {
			logger.Info("Cluster already exists", "cluster", clusterName)
			return true, nil
		}
	}
	return false, nil
}

func ClusterEndpoint(ctx context.Context, clusterName string) (host string, port int, err error) {
	logger := log.FromContext(ctx)
	logger.Info("Getting Kind cluster endpoint", "cluster", clusterName)
	cmdStr := "kind get kubeconfig -n mycluster | grep server | cut -d \"/\" -f 3"
	b, err := exec.Command("bash", "-c", cmdStr).Output()
	if err != nil {
		errors.Wrap(err, "failed to get kind cluster endpoint")
	}
	s := strings.ReplaceAll(strings.TrimSpace(string(b)), "\n", "")
	host = strings.Split(s, ":")[0]
	port, err = strconv.Atoi(strings.Split(s, ":")[1])
	return
}

func ClusterCreate(ctx context.Context, clusterName string) (err error) {
	logger := log.FromContext(ctx)
	logger.Info("Creating Kind cluster", "cluster", clusterName)
	cmd := exec.Command("kind", "create", "cluster", "--kubeconfig", fmt.Sprintf("/tmp/%s.kubeconfig", clusterName), "--wait", "120s", "--name", clusterName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		errors.Wrap(err, "failed to create kind cluster")
	}
	return
}

func ClusterDelete(ctx context.Context, clusterName string) (err error) {
	logger := log.FromContext(ctx)
	logger.Info("Deleting Kind cluster", "cluster", clusterName)
	cmd := exec.Command("kind", "delete", "cluster", "--name", clusterName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		errors.Wrap(err, "failed to delete kind cluster")
	}
	return
}
