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

package controllers

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/BartekTao/kubernetes-alertprobe-controller/api/v1alpha1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// AlertProbeReconciler reconciles a AlertProbe object
type AlertProbeReconciler struct {
	client.Client
	Log     logr.Logger
	Scheme  *runtime.Scheme
	cancels sync.Map
}

func NewAlertProbeReconciler(client client.Client, scheme *runtime.Scheme, log logr.Logger) *AlertProbeReconciler {
	return &AlertProbeReconciler{
		Client:  client,
		Scheme:  scheme,
		Log:     log,
		cancels: sync.Map{},
	}
}

//+kubebuilder:rbac:groups=probe.rextein.com,resources=alertprobes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=probe.rextein.com,resources=alertprobes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=probe.rextein.com,resources=alertprobes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AlertProbe object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *AlertProbeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	log := r.Log.WithValues("alertprobe", req.NamespacedName)

	var alertProbe v1alpha1.AlertProbe
	if err := r.Get(ctx, req.NamespacedName, &alertProbe); err != nil {
		if errors.IsNotFound(err) {
			// Cancel the goroutine if the AlertProbe was deleted
			if cancel, exists := r.cancels.Load(req.NamespacedName.String()); exists {
				cancel.(context.CancelFunc)()
				r.cancels.Delete(req.NamespacedName.String())
			}
			return ctrl.Result{}, nil
		}
		// handle other errors
		log.Error(err, "unable to fetch AlertProbe")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Cancel the existing goroutine if one already exists for the alertProbe
	if cancel, ok := r.cancels.Load(req.NamespacedName); ok {
		cancel.(context.CancelFunc)()
	}

	ctxProbe, cancel := context.WithCancel(ctx)

	r.cancels.Store(req.NamespacedName, cancel)

	go func() {
		ticker := time.NewTicker(time.Duration(alertProbe.Spec.PeriodSeconds) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctxProbe.Done():
				log.Info("stopping goroutine", "url", alertProbe.Spec.URL)
				return
			case <-ticker.C:
				log.Info("checking url", "url", alertProbe.Spec.URL)
				res, err := http.Get(alertProbe.Spec.URL)
				if err != nil {
					log.Error(err, "unable to send GET request")
					return
				}
				if res.StatusCode != 200 {
					notify("URL check failed for " + req.NamespacedName.String())
				}

				alertProbe.Status.LastCheckTime = metav1.Now()
				alertProbe.Status.LastCheckResult = res.Status
				if err := r.Status().Update(ctxProbe, &alertProbe); err != nil {
					log.Error(err, "unable to update AlertProbe status")
				}
				_ = res.Body.Close()
			}
		}
	}()

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AlertProbeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.AlertProbe{}).
		Complete(r)
}

func notify(msg string) {
	// implement your notification logic here
	println(msg)
}
