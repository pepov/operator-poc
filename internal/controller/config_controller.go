/*
Copyright 2023 Pepov.

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
	"fmt"
	"reflect"

	"github.com/go-logr/logr"
	"github.com/pepov/operator-poc/api/v1beta1/applyconfigurations/api/v1beta1"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	whateverv1beta1 "github.com/pepov/operator-poc/api/v1beta1"
)

// ConfigReconciler reconciles a Config object
type ConfigReconciler struct {
	client.Client
	Logger logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=whatever.example.org,resources=configs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=whatever.example.org,resources=configs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=whatever.example.org,resources=configs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Config object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *ConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	c := &whateverv1beta1.Config{}
	if err := r.Client.Get(ctx, req.NamespacedName, c); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	actual := func() (interface{}, error) {
		if actual, err := v1beta1.ExtractConfig(c, "whatever-operator"); err != nil {
			return nil, errors.Wrap(err, "failed to extract apply config from original object")
		} else {
			// broken in ExtractConfig
			actual.WithAPIVersion("whatever.example.org/v1beta1")
			return actual, nil
		}
	}

	desired := func() (interface{}, error) {
		if desired, err := v1beta1.ExtractConfig(c, "whatever-operator"); err != nil {
			return nil, errors.Wrap(err, "failed to extract apply config from original object")
		} else {
			// broken in ExtractConfig
			desired.WithAPIVersion("whatever.example.org/v1beta1")
			desired.WithSpec(v1beta1.ConfigSpec().WithFoo("set by controller"))
			return desired, nil
		}
	}

	err := applyIfChanged(actual, desired, executeApply(ctx, r.Logger, r.Client, "whatever-operator"))

	return ctrl.Result{}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&whateverv1beta1.Config{}).
		Complete(r)
}

// executeApply is fully generic
func executeApply(ctx context.Context, l logr.Logger, r client.Client, fieldManager string) func(ac interface{}) error {
	return func(ac interface{}) error {
		var err error
		var u map[string]interface{}
		if u, err = runtime.DefaultUnstructuredConverter.ToUnstructured(ac); err != nil {
			return errors.Wrap(err, fmt.Sprintf("unable to convert apply config %+v to unstructured", ac))
		}

		o := &unstructured.Unstructured{
			Object: u,
		}

		if err := r.Patch(ctx, o, client.Apply, client.FieldOwner(fieldManager), client.ForceOwnership); err != nil {
			return errors.Wrap(err, fmt.Sprintf("patch failed for object %+v", o))
		}

		l.Info("object updated", "object", o)

		return nil

	}
}

func applyIfChanged(
	actual func() (any, error),
	desired func() (any, error),
	mutateFn func(interface{}) error) error {

	var err error
	var a, d any
	if a, err = actual(); err != nil {
		return err
	}

	if d, err = desired(); err != nil {
		return err
	}

	if reflect.DeepEqual(a, d) {
		return nil
	}

	return mutateFn(d)
}
