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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	maintenanceoperatoriov1alpha1 "github.com/abdurrehman107/maintenance-window/api/v1alpha1"
)

var _ = Describe("MaintenanceWindow Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "mw-test-resource"
		const namespace = "maintenance-window-system"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: namespace,
		}

		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: namespace},
		}

		var testMaintenanceWindow = &maintenanceoperatoriov1alpha1.MaintenanceWindow{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "maintenanceoperator.io.maintenanceoperator.io/v1alpha1",
				Kind:       "MaintenanceWindow",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      resourceName,
				Namespace: namespace,
			},
			Spec: maintenanceoperatoriov1alpha1.MaintenanceWindowSpec{
				StartTime: "2025-06-19T01:00:00Z",
				EndTime:   "2025-06-25T03:00:00Z",
			},
		}

		BeforeEach(func() {
			By("creating the custom resource and namespace for the Kind MaintenanceWindow")
			_ = k8sClient.Create(ctx, ns)

			err := k8sClient.Get(ctx, typeNamespacedName, testMaintenanceWindow)

			if err != nil && errors.IsNotFound(err) {
				Expect(k8sClient.Create(ctx, testMaintenanceWindow)).To(Succeed())
			}
		})

		AfterEach(func() {
			// TODO(user): Cleanup logic after each test, like removing the resource instance.
			resource := &maintenanceoperatoriov1alpha1.MaintenanceWindow{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance MaintenanceWindow")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})
		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &MaintenanceWindowReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
