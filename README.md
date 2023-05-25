# cluster-api-provider-mailgun
## Description
Provider to run mailgun on cluster resources.<br>
Cluster API is used to implement Cluster Provider.<br> 
Although Machines should also be implemented, only Clusters are used in this study in order to learn how to implement them.<br>
This Provider sends mail to the destination specified in the Configmap, depending on the Requester specified in the .spec field of the MailgunCluster resource.

## Getting Started
### Prerequisites
- It is implemented by the following versions of the tool.
    - Kubebuilder("v3.10.0")
    - Golang("v1.19.9")
    - Kubernetes("v1.26.3")
    - make("v4.3")
    - gcc("v11.3.0")
- When using the latest version v1.4.2 of the Cluster API, the following versions of apimachinery, clinet-go, and controller-runtime must be used.
    - k8s.io/apimachinery v0.26.5
	- k8s.io/client-go v0.26.5
	- sigs.k8s.io/cluster-api v1.4.2
	- sigs.k8s.io/controller-runtime v0.14.6
- You must have already registered an account with Mailgun.<br>
You do not need to create a domain name. In this case, we will use the Sandbox domain.
### Running on the cluster
1. Setting Environment Variables<br>
The configmap and secret resources are set using this environment variable.
```shell
export MAILGUN_DOMAIN="<Your Mailgun Sandbox Domain>"
export MAILGUN_API_KEY="<Your Mailgun Private API key>"
export MAIL_RECIPIENT="<Mail recipient address>"
```

2. Build and push your image to the location specified by `IMG`:
```sh
make docker-build docker-push IMG=<Your Registry>/cluster-api-provider-mailgun:tag
```

3. Deploy the controller to the cluster with the image specified by `IMG`:
```sh
make deploy IMG=<some-registry>/cluster-api-provider-mailgun:tag
```

4. Edit Custom Resources
Replace the .spec.requester in MailgunCluster with the Sandbox domain.
```
cat << EOT > config/samples/infrastructure_v1beta1_mailguncluster.yaml
apiVersion: cluster.x-k8s.io/v1beta1
kind: Cluster
metadata:
  name: hello-mailgun
  namespace: cluster-api-provider-mailgun-system
spec:
  clusterNetwork:
    pods:
      cidrBlocks: ["192.168.0.0/16"]
  infrastructureRef:
    apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
    kind: MailgunCluster
    name: hello-mailgun
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: MailgunCluster
metadata:
  name: hello-mailgun
  namespace: cluster-api-provider-mailgun-system
  labels:
    app.kubernetes.io/name: mailguncluster
    app.kubernetes.io/instance: mailguncluster-sample
    app.kubernetes.io/part-of: cluster-api-provider-mailgun
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/created-by: cluster-api-provider-mailgun
spec:
  priority: "ExtremelyUrgent"
  request: "Please make me a cluster, with sugar on top?"
  requester: "<Your Domain>"
EOT
```

5. Install Instances of Custom Resources:
```sh
kubectl apply -f config/samples/
```

6. Check to see if the mail has been delivered to the address specified in the MAIL_RECIPIENT environment variable.<br>
It is also possible to check from Mailgun's Logs as follows.
![Mailgun Log](https://github.com/jnytnai0613/cluster-api-provider-mailgun/blob/Update-README/docs/Mailgun_log.png)

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller from the cluster:

```sh
make undeploy
```

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/),
which provide a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
export MAILGUN_DOMAIN="<Your Mailgun Sandbox Domain>"
export MAILGUN_API_KEY="<Your Mailgun Private API key>"
export MAIL_RECIPIENT="<Mail recipient address>"
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

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

