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

package v1

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	maintenancecustomiov1 "github.com/abdurrehman107/maintenance-window.git/api/v1"
)

// nolint:unused
// log is for logging in this package.
var maintenancewindowlog = logf.Log.WithName("maintenancewindow-resource")

// SetupMaintenanceWindowWebhookWithManager registers the webhook for MaintenanceWindow in the manager.
func SetupMaintenanceWindowWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&maintenancecustomiov1.MaintenanceWindow{}).
		WithValidator(&MaintenanceWindowCustomValidator{}).
		WithDefaulter(&MaintenanceWindowCustomDefaulter{
			// Active: false,
		}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-maintenance-custom-io-maintenance-window-io-v1-maintenancewindow,mutating=true,failurePolicy=fail,sideEffects=None,groups=maintenance.custom.io.maintenance-window.io,resources=maintenancewindows,verbs=create;update,versions=v1,name=mmaintenancewindow-v1.kb.io,admissionReviewVersions=v1

// MaintenanceWindowCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind MaintenanceWindow when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type MaintenanceWindowCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
	Active bool
}

var _ webhook.CustomDefaulter = &MaintenanceWindowCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind MaintenanceWindow.
func (d *MaintenanceWindowCustomDefaulter) Default(ctx context.Context, obj runtime.Object) error {
	maintenancewindow, ok := obj.(*maintenancecustomiov1.MaintenanceWindow)

	if !ok {
		return fmt.Errorf("expected an MaintenanceWindow object but got %T", obj)
	}
	maintenancewindowlog.Info("Defaulting for MaintenanceWindow", "name", maintenancewindow.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-maintenance-custom-io-maintenance-window-io-v1-maintenancewindow,mutating=false,failurePolicy=fail,sideEffects=None,groups=maintenance.custom.io.maintenance-window.io,resources=maintenancewindows,verbs=create;update,versions=v1,name=vmaintenancewindow-v1.kb.io,admissionReviewVersions=v1

// MaintenanceWindowCustomValidator struct is responsible for validating the MaintenanceWindow resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type MaintenanceWindowCustomValidator struct {
	//TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &MaintenanceWindowCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type MaintenanceWindow.
func (v *MaintenanceWindowCustomValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	maintenancewindow, ok := obj.(*maintenancecustomiov1.MaintenanceWindow)
	if !ok {
		return nil, fmt.Errorf("expected a MaintenanceWindow object but got %T", obj)
	}
	maintenancewindowlog.Info("Validation for MaintenanceWindow upon creation", "name", maintenancewindow.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type MaintenanceWindow.
func (v *MaintenanceWindowCustomValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	maintenancewindow, ok := newObj.(*maintenancecustomiov1.MaintenanceWindow)
	if !ok {
		return nil, fmt.Errorf("expected a MaintenanceWindow object for the newObj but got %T", newObj)
	}
	maintenancewindowlog.Info("Validation for MaintenanceWindow upon update", "name", maintenancewindow.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type MaintenanceWindow.
func (v *MaintenanceWindowCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	maintenancewindow, ok := obj.(*maintenancecustomiov1.MaintenanceWindow)
	if !ok {
		return nil, fmt.Errorf("expected a MaintenanceWindow object but got %T", obj)
	}
	maintenancewindowlog.Info("Validation for MaintenanceWindow upon deletion", "name", maintenancewindow.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
