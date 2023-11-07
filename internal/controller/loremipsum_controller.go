/*
Copyright 2023.

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
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apiv1 "github.com/broswen/loremipsum/api/v1"
)

// LoremIpsumReconciler reconciles a LoremIpsum object
type LoremIpsumReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

var LOREM_IPSUM = []string{
	"lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	"cursus risus at ultrices mi tempus imperdiet nulla malesuada.",
	"pellentesque pulvinar pellentesque habitant morbi tristique senectus et netus.",
	"dui faucibus in ornare quam viverra orci sagittis.",
	"sollicitudin aliquam ultrices sagittis orci a scelerisque purus semper eget. ",
	"sodales neque sodales ut etiam sit amet nisl purus in.",
}

//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=api.broswen.com,resources=loremipsums,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=api.broswen.com,resources=loremipsums/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=api.broswen.com,resources=loremipsums/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LoremIpsum object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *LoremIpsumReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	li := apiv1.LoremIpsum{}
	cm := v1.ConfigMap{}
	// get loremipsum
	if err := r.Client.Get(ctx, req.NamespacedName, &li); err != nil {
		// fail if error other than not-found
		if client.IgnoreNotFound(err) != nil {
			logger.Error(err, "failed to get loremipsum", "namespace", req.Namespace, "name", req.Name)
			return ctrl.Result{
				Requeue: true,
			}, err
		}
		cm = v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      req.Name,
				Namespace: req.Namespace,
			},
		}
		// delete child configmap if loremipsum is not found
		if err := r.Client.Delete(ctx, &cm); client.IgnoreNotFound(err) != nil {
			logger.Error(err, "failed to delete child configmap", "namespace", req.Namespace, "name", req.Name)
			return ctrl.Result{
				Requeue: true,
			}, err
		}
		// success
		return ctrl.Result{}, nil
	}

	// generate loremipsum content per spec
	content := make([]string, 0)
	for i := 0; i < li.Spec.Lines; i++ {
		if li.Spec.Capitalize {
			content = append(content, strings.ToUpper(LOREM_IPSUM[i]))
		} else {
			content = append(content, LOREM_IPSUM[i])
		}
	}
	data := strings.Join(content, "\n")

	// get configmap for this loremipsum
	if err := r.Client.Get(ctx, req.NamespacedName, &cm); err != nil {
		if client.IgnoreNotFound(err) != nil {
			logger.Error(err, "unable to get child ConfigMap", "namespace", req.Namespace, "name", req.Name)
			return ctrl.Result{
				Requeue: true,
			}, err
		}

		// if configmap is not found, create it
		cm = v1.ConfigMap{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ConfigMap",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      req.Name,
				Namespace: req.Namespace,
			},
			Data: map[string]string{"data": data},
		}
		if err := r.Client.Create(ctx, &cm); err != nil {
			logger.Error(err, "unable to create child ConfigMap", "namespace", req.Namespace, "name", req.Name)
			return ctrl.Result{
				Requeue: true,
			}, err
		}
		return ctrl.Result{}, nil
	}
	// if configmap exists, set new data and update it
	cm.Data["data"] = data
	if err := r.Client.Update(ctx, &cm); err != nil {
		logger.Error(err, "unable to update child ConfigMap", "namespace", req.Namespace, "name", req.Name)
		return ctrl.Result{
			Requeue: true,
		}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LoremIpsumReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1.LoremIpsum{}).
		Complete(r)
}
