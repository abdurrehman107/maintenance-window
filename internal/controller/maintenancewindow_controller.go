/*
Copyright 2024.

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
	// "time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	maintenancecustomiov1 "github.com/abdurrehman107/maintenance-window.git/api/v1"
)

// MaintenanceWindowReconciler reconciles a MaintenanceWindow object
type MaintenanceWindowReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=maintenance.custom.io.maintenence-window.io,resources=maintenancewindows,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=maintenance.custom.io.maintenence-window.io,resources=maintenancewindows/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=maintenance.custom.io.maintenence-window.io,resources=maintenancewindows/finalizers,verbs=update
// +kubebuilder:rbac:groups=maintenance.custom.io.maintenence-window.io,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=maintenance.custom.io.maintenence-window.io,resources=deployments/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MaintenanceWindow object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *MaintenanceWindowReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	l.Info("Initiated logger")
	var maintenanceWindow maintenancecustomiov1.MaintenanceWindow

	// Populate the maintenanceWindow object
	if err := r.Get(ctx, req.NamespacedName, &maintenanceWindow); err != nil {
		l.Error(err, "Unable to get the objects for maintenance window")
		return ctrl.Result{}, err
	}

	// Set the Status.Active to false upon detection of our object 
	maintenanceWindow.Status.Active = false

	// Fetch startTime and endTime
	var startTime *metav1.Time
	var endTime *metav1.Time
	var currentTime metav1.Time
	startTime = &maintenanceWindow.Spec.StartTime
	endTime = &maintenanceWindow.Spec.EndTime
	currentTime = metav1.Now() // Fetch current time
	
	// scheduleMaintenance() begins a maintenance window whenever called
	var scheduledMaintenance = func(maintenanceWindow *maintenancecustomiov1.MaintenanceWindow) {
		maintenanceWindow.Status.Active = true
	}
	// unscheduleMaintenance() ends a maintenance window
	var unscheduleMaintenance = func(maintenceWindow *maintenancecustomiov1.MaintenanceWindow) {
		maintenceWindow.Status.Active = false
	}
	// Validate currentTime is not more than startTime and endTime.
	//
	// 1. Scheudle a maintenance window if current time is less than startTime
	// 2. If current time is greater than startTime and less than the endTime then immediately implement the maintenance window.
	// 3. If it is greater than the start and endTime both then add the respective object to the list of completed windows // delete it
	if currentTime.Before(startTime) {
		// schedule the maintenance window object to start at the respective startTime (to be addressed later)
		return ctrl.Result{RequeueAfter: startTime.Rfc3339Copy().Sub(currentTime.UTC())}, nil
	} else if currentTime.After(startTime.Time) && currentTime.Before(endTime) {
		// start the maintenance window right away
		scheduledMaintenance(&maintenanceWindow)
		return ctrl.Result{RequeueAfter: currentTime.Rfc3339Copy().Sub(endTime.UTC())}, nil
	} else if currentTime.After(endTime.Time) {
		unscheduleMaintenance(&maintenanceWindow)
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MaintenanceWindowReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&maintenancecustomiov1.MaintenanceWindow{}).
		Named("maintenancewindow").
		Complete(r)
}
