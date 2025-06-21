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

package v1

import (
	"context"
	"fmt"
	"net/http"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	maintenanceoperatoriov1alpha1 "github.com/abdurrehman107/maintenance-window/api/v1alpha1"
)

// nolint:unused
// log is for logging in this package.
var deploymentlog = logf.Log.WithName("deployment-resource")

// SetupDeploymentWebhookWithManager registers the webhook for Deployment in the manager.
func SetupDeploymentWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&appsv1.Deployment{}).
		WithValidator(&DeploymentCustomValidator{Client: mgr.GetClient()}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-apps-v1-deployment,mutating=false,failurePolicy=fail,sideEffects=None,groups=apps,resources=deployments,verbs=create;update,versions=v1,name=vdeployment-v1.kb.io,admissionReviewVersions=v1

// DeploymentCustomValidator struct is responsible for validating the Deployment resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type DeploymentCustomValidator struct {
	Client client.Client
}

var _ webhook.CustomValidator = &DeploymentCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type Deployment.
func (v *DeploymentCustomValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	deployment, ok := obj.(*appsv1.Deployment)
	if !ok {
		return nil, fmt.Errorf("expected a Deployment object but got %T", obj)
	}
	deploymentlog.Info("Validation for Deployment upon creation", "name", deployment.GetName())

	// TODO(user): fill in your validation logic upon object creation.
	var mwList maintenanceoperatoriov1alpha1.MaintenanceWindowList
	if err := v.Client.List(ctx, &mwList); err != nil {
		return nil, fmt.Errorf("unable to get maintenance window: %v", err)
	}
	if len(mwList.Items) == 0 {
		return nil, nil
	}
	for _, mw := range mwList.Items {
		if mw.Status.State == maintenanceoperatoriov1alpha1.StateActive {
			deploymentlog.Info(fmt.Sprintf("blocked by maintenance window %q until %s", mw.Name, mw.Spec.EndTime))
			return admission.Warnings{fmt.Sprintf("blocked by maintenance window %q until %s", mw.Name, mw.Spec.EndTime)}, fmt.Errorf("blocked by maintenance window %q until %s", mw.Name, mw.Spec.EndTime)
		}
	}
	// if mw.Status.State == maintenanceoperatoriov1alpha1.StateActive {
	// 	deploymentlog.Info(fmt.Sprintf("blocked by maintenance window %q until %s", mw.Name, mw.Spec.EndTime))
	// 	return admission.Warnings{fmt.Sprintf("blocked by maintenance window %q until %s", mw.Name, mw.Spec.EndTime)}, fmt.Errorf("blocked by maintenance window %q until %s", mw.Name, mw.Spec.EndTime)
	// }

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type Deployment.
func (v *DeploymentCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	deployment, ok := newObj.(*appsv1.Deployment)
	if !ok {
		return nil, fmt.Errorf("expected a Deployment object for the newObj but got %T", newObj)
	}
	deploymentlog.Info("Validation for Deployment upon update", "name", deployment.GetName())

	// TODO(user): fill in your validation logic upon object update.
	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type Deployment.
func (v *DeploymentCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	deployment, ok := obj.(*appsv1.Deployment)
	if !ok {
		return nil, fmt.Errorf("expected a Deployment object but got %T", obj)
	}
	deploymentlog.Info("Validation for Deployment upon deletion", "name", deployment.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}

func (w *DeploymentCustomValidator) Handle(ctx context.Context) admission.Response {
	var mwList maintenanceoperatoriov1alpha1.MaintenanceWindowList

	if err := w.Client.List(ctx, &mwList,
		client.MatchingFields{"status.state": "active"}); err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	// if there is a maintenance window active, deny the deployment
	if len(mwList.Items) > 0 {
		msg := fmt.Sprintf("blocked by maintenance window %q until %s", mwList.Items[0].Name, mwList.Items[0].Spec.EndTime)
		return admission.Denied(msg)
	}

	return admission.Allowed("no maintenance window active")
}
