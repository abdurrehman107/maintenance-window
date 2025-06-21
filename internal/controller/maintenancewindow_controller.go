/*
Copyright 2025.

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

package controller

import (
	"context"
	"time"

	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	maintenanceoperatoriov1alpha1 "github.com/abdurrehman107/maintenance-window/api/v1alpha1"
)

// MaintenanceWindowReconciler reconciles a MaintenanceWindow object
type MaintenanceWindowReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=maintenanceoperator.io.maintenanceoperator.io,resources=maintenancewindows,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=maintenanceoperator.io.maintenanceoperator.io,resources=maintenancewindows/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=maintenanceoperator.io.maintenanceoperator.io,resources=maintenancewindows/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MaintenanceWindow object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.1/pkg/reconcile
func (r *MaintenanceWindowReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// fetch the maintenance window object
	var mw maintenanceoperatoriov1alpha1.MaintenanceWindow
	if err := r.Get(ctx, req.NamespacedName, &mw); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// extract the start, end and current time
	// var startTime, endTime metav1.Time // time will look something like this 2025-03-25T00:00:00Z
	var err error

	startTime, err := time.Parse(time.RFC3339, mw.Spec.StartTime)
	if err != nil {
		logger.Error(err, "unable to parse start time")
		return ctrl.Result{}, err
	}
	
	endTime, err := time.Parse(time.RFC3339, mw.Spec.EndTime)
	if err != nil {
		logger.Error(err, "unable to parse end time")
		return ctrl.Result{}, err
	}

	currentTime := time.Now()

	if !currentTime.Equal(startTime) && currentTime.Before(endTime) {
		// begin maintenance window change state to true
		mw.Status.State = maintenanceoperatoriov1alpha1.StateActive
	} else if !currentTime.Before(endTime) && !currentTime.Equal(endTime) {
		// change state to expired
		mw.Status.State = maintenanceoperatoriov1alpha1.StateExpired
	} else if currentTime.Before(startTime) {
		mw.Status.State = maintenanceoperatoriov1alpha1.StateInactive
	}

	// update the object in the cluster
	if err = r.Status().Update(ctx, &mw); err != nil {
		logger.Error(err, "unable to update the mw object")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MaintenanceWindowReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&maintenanceoperatoriov1alpha1.MaintenanceWindow{}).
		Named("maintenancewindow").
		Complete(r)
}
