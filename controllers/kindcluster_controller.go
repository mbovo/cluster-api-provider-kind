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

	"github.com/go-logr/logr"
	infrastructurev1beta1 "github.com/mbovo/cluster-api-provider-kind/api/v1beta1"
	"github.com/mbovo/cluster-api-provider-kind/pkg/kind"
	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// KindClusterReconciler reconciles a KindCluster object
type KindClusterReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	Recorder   record.EventRecorder
	patcher    *patch.Helper
	kindHelper kind.KindHelper
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

	kindCluster := &infrastructurev1beta1.KindCluster{}
	if err := r.Client.Get(ctx, req.NamespacedName, kindCluster); err != nil {
		logger.Info(fmt.Sprintf("Failed to get KindCluster resource '%s/%s'.", req.NamespacedName.Namespace, req.NamespacedName.Name))
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

	// set up the patch helper
	r.patcher, err = patch.NewHelper(kindCluster, r.Client)
	if err != nil {
		logger.Error(err, "cannot create patch helper")
		return reconcile.Result{}, err
	}
	// set up kind helper
	r.kindHelper = kind.NewKindLibHelper(logger, kindCluster, cluster)

	// Handle deleted clusters
	if !kindCluster.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, kindCluster)
	}

	// Handle non-deleted clusters
	return r.reconcileNormal(ctx, kindCluster)
}

func (r *KindClusterReconciler) reconcileNormal(ctx context.Context, kindCluster *infrastructurev1beta1.KindCluster) (reconcile.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling KindCluster")
	// setting up the finalizer immediatlye
	ok := controllerutil.AddFinalizer(kindCluster, infrastructurev1beta1.KindClusterFinalizer)
	if !ok {
		logger.Info("Finalizer already exists")
	}

	if err := r.patcher.Patch(ctx, kindCluster); err != nil {
		return reconcile.Result{}, err
	}

	logger.Info("Search for already existing cluster", "cluster", kindCluster.ObjectMeta.Name)
	ok, err := r.kindHelper.Exists(ctx, kindCluster)
	if err != nil {
		return reconcile.Result{}, err
	}

	if !ok {
		logger.Info("Creating cluster", "cluster", kindCluster.ObjectMeta.Name)
		err := r.kindHelper.Create(ctx, kindCluster)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	host, port, err := r.kindHelper.Endpoint(ctx, kindCluster)
	if err != nil {
		return reconcile.Result{}, err
	}

	logger.Info("Setting controlPlaneEndpoint", "host", host, "port", port)

	kindCluster.Spec.ControlPlaneEndpoint.Host = host
	kindCluster.Spec.ControlPlaneEndpoint.Port = int32(port)
	kindCluster.Status.Ready = true
	kindCluster.Status.Nodes = kindCluster.Spec.WorkerCount + kindCluster.Spec.ControlPlaneCount
	r.patcher.Patch(ctx, kindCluster)

	return reconcile.Result{}, nil
}

func (r *KindClusterReconciler) reconcileDelete(ctx context.Context, kindCluster *infrastructurev1beta1.KindCluster) (reconcile.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Deleting KindCluster")

	// delete logic

	if ok, err := r.kindHelper.Exists(ctx, kindCluster); ok == false || err != nil {
		logger.Info("Cluster does not exist, skip", "cluster", kindCluster.ObjectMeta.Name)
	}

	err := r.kindHelper.Delete(ctx, kindCluster)
	if err != nil {
		return reconcile.Result{}, err
	}

	controllerutil.RemoveFinalizer(kindCluster, infrastructurev1beta1.KindClusterFinalizer)
	if err := r.patcher.Patch(ctx, kindCluster); err != nil {
		return reconcile.Result{}, err
	}

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

	controller, err := ctrl.NewControllerManagedBy(mgr).
		For(&infrastructurev1beta1.KindCluster{}).
		WithEventFilter(predicates.ResourceNotPaused(log)).
		WithEventFilter(predicates.ResourceIsNotExternallyManaged(log)).
		WithEventFilter(filterStatusChanges()).
		Build(r)

	if err != nil {
		return errors.Wrap(err, "failed setting up with a controller manager")
	}

	// Watch for changes to Cluster resources to enque referenced KindCluster resources for reconciliation.
	return controller.Watch(
		&source.Kind{Type: &clusterv1.Cluster{}},
		handler.EnqueueRequestsFromMapFunc(r.mapClusterToKindCluster(ctx, log)),
		predicates.ClusterUnpaused(log))
}

// Map a KindCluster to the Cluster that is referencing it.
func (r *KindClusterReconciler) mapClusterToKindCluster(ctx context.Context, log logr.Logger) handler.MapFunc {
	return func(o client.Object) []ctrl.Request {

		c, ok := o.(*clusterv1.Cluster)
		if !ok {
			log.Error(errors.Errorf("expected a Cluster but got a %T", o), "failed to get KindCluster for Cluster")
			return nil
		}
		log = log.WithValues("func", "mapClusterToKindCluster", "cluster", c.Name, "namespace", c.Namespace)

		if !c.ObjectMeta.DeletionTimestamp.IsZero() {
			log.V(3).Info("Ignoring deleting Cluster")
			return nil
		}

		if c.Spec.InfrastructureRef.GroupVersionKind().Kind != "KindCluster" {
			log.V(3).Info("Cluster Has an InfrastructureRef of a wrong type, skipping resource")
			return nil
		}

		kindCluster := &infrastructurev1beta1.KindCluster{}
		namspacedName := types.NamespacedName{Namespace: c.Spec.InfrastructureRef.Namespace, Name: c.Spec.InfrastructureRef.Name}

		if err := r.Client.Get(ctx, namspacedName, kindCluster); err != nil {
			log.V(3).Info("InfrastructureRef field not set yet, skipping resource")
			return nil
		}

		if annotations.IsExternallyManaged(kindCluster) {
			log.V(3).Info("Ignoring externally managed KindCluster")
			return nil
		}

		log.V(3).Info("Enqueueing request for KindCluster", c.Spec.InfrastructureRef.Name)
		return []ctrl.Request{
			{NamespacedName: client.ObjectKey{Namespace: c.Spec.InfrastructureRef.Namespace, Name: c.Spec.InfrastructureRef.Name}},
		}
	}
}
