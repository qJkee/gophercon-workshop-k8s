/*
Copyright 2021.

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
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	workshopv1alpha1 "github.com/workshop/mysql-operator/api/v1alpha1"
)

// CustomMysqlReconciler reconciles a CustomMysql object
type CustomMysqlReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=workshop.example.com,resources=custommysqls,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=workshop.example.com,resources=custommysqls/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=workshop.example.com,resources=custommysqls/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=pods;pods/exec;services,verbs=get;list;watch;create;update;delete;patch
//+kubebuilder:rbac:groups=apps,resources=deployments;replicasets;statefulsets,verbs=get;list;watch;create;update;delete;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CustomMysql object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *CustomMysqlReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("custommysql", req.NamespacedName)

	cr := &workshopv1alpha1.CustomMysql{}

	err := r.Client.Get(ctx, types.NamespacedName{
		Namespace: req.Namespace,
		Name:      req.Name,
	}, cr)

	if err != nil {
		return ctrl.Result{RequeueAfter: 5 * time.Second}, err
	}

	err = r.CreateService(cr)
	if err != nil {
		return ctrl.Result{RequeueAfter: 5 * time.Second}, err
	}

	err = r.CreateSfs(cr)
	if err != nil {
		return ctrl.Result{RequeueAfter: 5 * time.Second}, err
	}

	r.manageReplication(cr)

	return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CustomMysqlReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&workshopv1alpha1.CustomMysql{}).
		Complete(r)
}
