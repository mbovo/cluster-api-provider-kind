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

package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/mbovo/cluster-api-provider-kind/api/v1alpha1"
	infrastructurev1alpha1 "github.com/mbovo/cluster-api-provider-kind/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// KindClusterReconciler reconciles a KindCluster object
type KindClusterReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=kindclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=kindclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=kindclusters/finalizers,verbs=update
//+kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the KindCluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *KindClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, err error) {
	logger := log.FromContext(ctx)

	kindCluster := &v1alpha1.KindCluster{}
	if err := r.Client.Get(ctx, req.NamespacedName, kindCluster); err != nil {
		logger.Info(fmt.Sprintf("Failed to get KindClusteter resource '%s/%s'.", req.NamespacedName.Namespace, req.NamespacedName.Name))
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	// Fetch the Cluster:
	cluster, err := util.GetOwnerCluster(ctx, r.Client, kindCluster.ObjectMeta)
	if err != nil {
		return reconcile.Result{}, err
	}
	if cluster == nil {
		logger.Info("Cluster Controller has not yet set OwnerRef")
		return reconcile.Result{Requeue: true, RequeueAfter: 5 * time.Second}, nil
	}

	if annotations.IsPaused(cluster, kindCluster) {
		r.Recorder.Eventf(kindCluster, corev1.EventTypeNormal, "ClusterPaused", "Cluster is paused")
		logger.Info("KindCluster or linked Cluster is marked as paused, will not reconcile")
		return reconcile.Result{}, nil
	}

	// Handle deleted clusters
	if !kindCluster.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, kindCluster)
	}

	// Handle non-deleted clusters
	return r.reconcileNormal(ctx, kindCluster)
}

func (r *KindClusterReconciler) reconcileNormal(ctx context.Context, kindCluster *v1alpha1.KindCluster) (reconcile.Result, error) {
	logger := log.FromContext(ctx)
	_, err := patch.NewHelper(kindCluster, r.Client)
	if err != nil {
		log.Log.Error(err, "cannot create patch helper")
		return reconcile.Result{}, err
	}
	logger.Info("Reconciling KindCluster")

	return reconcile.Result{}, nil
}

func (r *KindClusterReconciler) reconcileDelete(ctx context.Context, kindCluster *v1alpha1.KindCluster) (reconcile.Result, error) {
	logger := log.FromContext(ctx)
	_, err := patch.NewHelper(kindCluster, r.Client)
	if err != nil {
		log.Log.Error(err, "cannot create patch helper")
		return reconcile.Result{}, err
	}
	logger.Info("Reconciling deleted KindCluster")

	return reconcile.Result{}, nil
}

// Filter out change of status subresource
func filterStatusChanges() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			// Ignore custom resource updates when metadata.Generation does not change
			// i.e. ignore change of evrything except "spec" field (https://github.com/kubernetes/kubernetes/issues/67428#issuecomment-456125985)
			return e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration()
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *KindClusterReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager) error {
	log := ctrl.LoggerFrom(ctx)

	return ctrl.NewControllerManagedBy(mgr).
		For(&infrastructurev1alpha1.KindCluster{}).
		WithEventFilter(predicates.ResourceNotPaused(log)).
		WithEventFilter(predicates.ResourceIsNotExternallyManaged(log)).
		WithEventFilter(filterStatusChanges()).
		Complete(r)
}
