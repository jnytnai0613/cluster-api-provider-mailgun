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
	"fmt"

	"github.com/BoltApp/mailgun-go/v4"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	infrav1betav1 "github.com/jnytnai0613/cluster-api-provider-mailgun/api/v1beta1"
)

// MailgunClusterReconciler reconciles a MailgunCluster object
type MailgunClusterReconciler struct {
	client.Client
	Mailgun   mailgun.Mailgun
	Recipient string
}

//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=mailgunclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=mailgunclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=mailgunclusters/finalizers,verbs=update
//+kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.6/pkg/reconcile
func (r *MailgunClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var mgCluster infrav1betav1.MailgunCluster
	if err := r.Get(ctx, req.NamespacedName, &mgCluster); err != nil {
		// 	import apierrors "k8s.io/apimachinery/pkg/api/errors"
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		logger.Error(err, "couldn't get obj")
		return ctrl.Result{}, err
	}

	/*
		cluster, err := util.GetOwnerCluster(ctx, r.Client, mgCluster.ObjectMeta)
		if err != nil {
			logger.Error(err, fmt.Sprintf("couldn't get ownercluster %q", mgCluster.GetName()))
			return ctrl.Result{}, err

		}
	*/

	if mgCluster.Status.MessageID != nil {
		// We already sent a message, so skip reconciliation
		return ctrl.Result{}, nil
	}

	//subject := fmt.Sprintf("[%s] New Cluster %s requested", mgCluster.Spec.Priority, cluster.Name)
	subject := fmt.Sprintf("[%s] New Cluster requested", mgCluster.Spec.Priority)
	body := fmt.Sprintf("Hello! One cluster please.\n\n%s\n", mgCluster.Spec.Request)
	requester := fmt.Sprintf("Cluster API Sandbox Provider <mailgun@%s>", mgCluster.Spec.Requester)
	msg := r.Mailgun.NewMessage(requester, subject, body, r.Recipient)
	_, msgID, err := r.Mailgun.Send(ctx, msg)
	if err != nil {
		logger.Error(err, fmt.Sprintf("couldn't send message %q", msgID))
		return ctrl.Result{}, err
	}

	helper, err := patch.NewHelper(&mgCluster, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}
	mgCluster.Status.MessageID = &msgID
	if err := helper.Patch(ctx, &mgCluster); err != nil {
		logger.Error(err, fmt.Sprintf("couldn't patch cluster %q", mgCluster.Name))
		return ctrl.Result{}, err
	}
	logger.Info(fmt.Sprintf("Message ID %s", *mgCluster.Status.MessageID))

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MailgunClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1betav1.MailgunCluster{}).
		Complete(r)
}
