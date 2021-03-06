# Introduction

After a conversation with my colleague [Tero](https://github.com/tahonen) we decided to prepare a session to highlight the benefits a developer can get from the Operators in general and from the [Operator Framework](https://github.com/operator-framework) in particular. For that session we depicted two demos: one to show how easy is to have a Kafka cluster on OpenShift (Red Hat's kubernetes distribution) and another one showing a custom operator that deployed and updated an application called **Gramola**.

This is the reposiory for the Operator created for that demo and effectively the demo itself.

## TL;DR

With this guide *you will learn how we created the operator to deploy version 0.0.1 of [Gramola]()* (a Java based application including an Angular UI, a Gateway and Events API and a PostgreSQL database). More importantly *you will also learn how we evolved the operator moving from version 0.0.1 to 0.0.2* (this included update the database schema and migrating data to the new schema). And even more importantly how we used the [Operator Lifecycle Manager](https://github.com/operator-framework/operator-lifecycle-manager) to do all this automatically.

The 2nd part of the guide explains how to run the demo, no need to code just enjoy deploying and upgrading our target application.

## Prerequisites

You need basic understanding of what an operator is to understand this guide
Additionally if you want to run the demo or create your own Operator you also need:

* [Go](https://golang.org/dl) 1.13.5+
* [Operator SDK](https://sdk.operatorframework.io/build/) v0.15.1+
* Free account in [Quay](https://quay.io) (this is needed to store the manifests that describe channels and versions of your operators)

Golang Based Operator SDK Installation
Follow the steps in the installation guide to learn how to install the Operator SDK CLI tool.

Additional Prerequisites https://sdk.operatorframework.io/docs/golang/installation/
git
go version v1.12+.
mercurial version 3.9+
docker version 17.03+.
kubectl version v1.11.3+.
Access to a Kubernetes v1.11.3+ cluster.

## Being grateful first
Parts of the code of this operator were borrowed from another [operator](https://github.com/mcouliba/openshift-workshop-operator) coded by my colleague [Madou](https://github.com/mcouliba) 

## About this guide

I have divided the guide in two parts:

1. The first part explains end to end how to create an Operator taking this operator as a starting point.
2. The second one explains how to run the demo consisting on deploying version 0.0.1 and then upgrade to 0.0.2 and see how the database schema is modified and data is migrated.

# Creting the Gramola Operator

## Installing the Operator SDK

In order to simplify and speed up the development of an operator we're going to install the [Operator SDK](https://sdk.operatorframework.io/build/).

```sh
export RELEASE_VERSION=v0.19.2
export OS=apple-darwin

mkdir ./bin

curl -LO https://github.com/operator-framework/operator-sdk/releases/download/${RELEASE_VERSION}/operator-sdk-${RELEASE_VERSION}-x86_64-${OS}

mv operator-sdk-${RELEASE_VERSION}-x86_64-${OS} ./bin/operator-sdk

chmod u+x ./bin/operator-sdk

export PATH=$(pwd)/bin:$PATH
```

## Log in to quay.io

You need the free account to use Quay as the repository for your operator manifests bundle. Then use those credentials to log in using docker/podman.

```sh
docker login quay.io
```

## Set up environment

Let's define some environment variables that will come handy later. Special attention to **GO111MODULE**.

```sh
export GO111MODULE=on

export OPERATOR_VENDOR="redhat" 
export APP_NAME="gramola"

export OPERATOR_NAME="${APP_NAME}-operator"
export OPERATOR_IMAGE="${APP_NAME}-operator-image"

export API_VERSION="${APP_NAME}.${OPERATOR_VENDOR}.com/v1alpha1"
export API_VERSION="v1alpha1"

export PROJECT_NAME=${OPERATOR_NAME}-project
``` 

## Create the scaffold for our operator

We're going to generate the Golang scaffold for our operator, just do as follows.

> **NOTE:** operator-sdk init generates a go.mod file to be used with Go modules. The --repo=<path> flag is required when creating a project outside of $GOPATH/src, as scaffolded files require a valid module path. Ensure you activate module support by running export GO111MODULE=on before using the SDK.

```sh
mkdir -p $GOPATH/src/github.com/${OPERATOR_VENDOR}

or 

mkdir -p ./${OPERATOR_VENDOR}/${OPERATOR_NAME}



cd $GOPATH/src/github.com/${OPERATOR_VENDOR}

or

cd ./${OPERATOR_VENDOR}/${OPERATOR_NAME}



operator-sdk new ${OPERATOR_NAME} --type=go --repo github.com/${OPERATOR_VENDOR}/${OPERATOR_NAME}

or

operator-sdk init --domain=${OPERATOR_VENDOR}.com --repo=github.com/${OPERATOR_VENDOR}/${OPERATOR_NAME}

```

## Add a new API for the custom resource AppService

> **NOTE:** IF error creating API ==> export GOROOT=$(go env GOROOT)

```sh
cd ./${OPERATOR_NAME}



operator-sdk add api --api-version=${API_VERSION} --kind=AppService

or

operator-sdk create api --group=${APP_NAME} --version=${API_VERSION} --kind=AppService
Create Resource [y/n]
y
Create Controller [y/n]
y
```

## [OLD] Add a new controller that watches for AppService


$ operator-sdk add controller --api-version=${API_VERSION} --kind=AppService


# Edit the CR [OLD]

code ./pkg/apis/gramola/<version>/<kind>_types.go

In this case: 

code ./pkg/apis/gramola/v1alpha1/appservice_types.go

  // AppServiceSpec defines the desired state of AppService
  type AppServiceSpec struct {
    // INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
    // Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
    // Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

    // Flags if the the AppService object is enabled or not
    Enabled bool `json:"enabled"`
    // Flags if the object has been initialized or not
    Initialized bool `json:"initialized"`
    // +kubebuilder:validation:Enum=Gramola,Gramophone,Phonograph
    Alias string `json:"alias,omitempty"``
  }

## Edit CR AppService

Go to `./pkg/apis/v1alpha1/appservice_types.go`

Find this:

```go
	// Foo is an example field of AppService. Edit AppService_types.go to remove/update
	Foo string `json:"foo,omitempty"`
```

And subsitute it with:

```go
	// +kubebuilder:validation:Minimum=0
	// Size is the size of the memcached deployment
	Size int32 `json:"size"
```

Go to

```go
// Nodes are the names of the memcached pods
	Nodes []string `json:"nodes"`
```

## Regenerate supporting code for your CRDs
operator-sdk generate k8s

## Regenerate your CRD OpenAPI definition [deprecated]
operator-sdk generate openapi

## Regenerate your CRD OpenAPI definition
operator-sdk generate crds

> **Build the latest openapi-gen from source!**
> `which ./bin/openapi-gen > /dev/null || go build -o ./bin/openapi-gen k8s.io/kube-openapi/cmd/openapi-gen`

## Run openapi-gen for each of your API group/version packages

```sh
./bin/openapi-gen --logtostderr=true -o "" -i ./pkg/apis/gramola/v1alpha1 -O zz_generated.openapi -p ./pkg/apis/gramola/v1alpha1 -h ./hack/boilerplate.go.txt -r "-"
```

# If you import a new modules run this

```sh
go mod vendor
```

# List module versions

```sh
go list -m -versions gopkg.in/src-d/go-git.v4
```

# Set the username variable

You should have an account in *quay.io*, if you don't please create it and use it here.

```sh
export USERNAME=<username>
```

# Build and push the ${OPERATOR_NAME} image to a public registry such as quay.io

```
export OPERATOR_VERSION=0.0.1
operator-sdk build quay.io/${USERNAME}/${OPERATOR_IMAGE}:${OPERATOR_VERSION}
```

# Login to public registry such as quay.io

```sh
docker login quay.io
```

# Push image

```sh
docker push quay.io/${USERNAME}/${OPERATOR_IMAGE}:${OPERATOR_VERSION}
```

# Update the operator manifest to use the built image name (if you are performing these steps on OSX, see note below)

The `operator-sdk` has generated a default `deployment` descriptor in `deploy/operator.yaml` but instead of pointing to a real image it contains a placeholder `REPLACE_IMAGE`, let's substitute it with the image we just built for image `0.0.1`.

```
sed -i "s|REPLACE_IMAGE|quay.io/${USERNAME}/${OPERATOR_IMAGE}\:${OPERATOR_VERSION}|g" deploy/operator.yaml
```

On OSX use:

```sh
sed -i "" "s|REPLACE_IMAGE|quay.io/${USERNAME}/${OPERATOR_IMAGE}\:${OPERATOR_VERSION}|g" deploy/operator.yaml
```

# Set/Create project

```sh
$ oc new-project ${PROJECT_NAME}
```

or

```sh
$ oc project ${PROJECT_NAME}
```

# Deploy basic elements to run the operator withouth OLM

Before we use OLM to deliver our operator we need to develop and test it locally.

1. Setup Service Account and RBAC

```sh
oc apply -f deploy/service_account.yaml
oc apply -f deploy/role.yaml
oc apply -f deploy/role_binding.yaml
```

2. Setup the CRD

```
oc apply -f deploy/crds/gramola.redhat.com_appservices_crd.yaml
```

# Run locally

The `operator-sdk` is prepared to run the code of our locally but within our kuberenetes cluster.

```sh
operator-sdk run --local --namespace ${PROJECT_NAME}
```

If there are no errors you should see something like this.

```sh
INFO[0000] Running the operator locally in namespace gramola-operator-project. 
{"level":"info","ts":1584114880.591337,"logger":"cmd","msg":"Operator Version: 0.0.1"}
{"level":"info","ts":1584114880.591379,"logger":"cmd","msg":"Go Version: go1.13.5"}
{"level":"info","ts":1584114880.591383,"logger":"cmd","msg":"Go OS/Arch: darwin/amd64"}
{"level":"info","ts":1584114880.5913868,"logger":"cmd","msg":"Version of operator-sdk: v0.15.1"}
{"level":"info","ts":1584114880.602292,"logger":"leader","msg":"Trying to become the leader."}
...
{"level":"info","ts":1584114899.0024478,"logger":"controller_appservice","msg":"Reconciling AppService","Request.Namespace":"gramola-operator-project","Request.Name":"example-appservice"}
```

You could also, alternativele, deploy the operator. But we're not going to do that for now.

```sh
oc apply -f deploy/operator.yaml
```

# Create an AppService CR

The default controller `pkg/controller/appservice/appservice_controller.go` will watch for `AppService` objects and create a pod for each CR.

Now in a different terminal window but *in the same path* let's create an example `AppService` object.

```sh
oc apply -f deploy/crds/gramola.redhat.com_v1alpha1_appservice_cr.yaml
```

You should get this error:

```sh 
The AppService "example-appservice" is invalid: 
* spec.enabled: Required value
* spec.initialized: Required value
```

Change the type [`AppService`](), so that initialized is not required

```go
Initialized bool `json:"initialized,omitempty"`
```

Let's change the CR so that `enabled` is defined:

```yaml
apiVersion: gramola.redhat.com/v1alpha1
kind: AppService
metadata:
  name: example-appservice
spec:
  enabled: true
```

Now if you try again... is thould work.

```
$ oc apply -f deploy/crds/gramola.redhat.com_v1alpha1_appservice_cr.yaml
appservice.gramola.redhat.com/example-appservice created
```

# Verify that a pod is created
$ oc get pod -l app=example-appservice
NAME                     READY     STATUS    RESTARTS   AGE
example-appservice-pod   1/1       Running   0          1m

# Test the new Resource Type
$ oc describe appservice example-appservice
Name:         example-appservice
Namespace:    gramola-operator-project
Labels:       <none>
Annotations:  kubectl.kubernetes.io/last-applied-configuration:
                {"apiVersion":"gramola.redhat.com/v1alpha1","kind":"AppService","metadata":{"annotations":{},"name":"example-appservice","namespace":"gram...
API Version:  gramola.redhat.com/v1alpha1
Kind:         AppService
Metadata:
  Creation Timestamp:  2020-03-12T10:30:09Z
  Generation:          1
  Resource Version:    920134
  Self Link:           /apis/gramola.redhat.com/v1alpha1/namespaces/gramola-operator-project/appservices/example-appservice
  UID:                 21dac31a-fc2f-4fb4-9d78-bc8d72531281
Spec:
  Enabled:  true
Events:     <none>

# Generate CSV 0.0.1

The `operator-sdk` tool will help us to create the CSV for our ClusterService `Gramola Operator` if you run it now, as follows, you'll see a **WARNING**.

> **NOTE 1:** Flag `--update-crds` commands `operator-sdk` to generate properties, definitions, etc. related to the CRDs your operator manages

> **NOTE 2:** go [here](https://github.com/operator-framework/operator-sdk/blob/master/doc/user/olm-catalog/generating-a-csv.md#csv-fields) to find out which properties are mandatory, optional, etc. in a CSV

```sh
$ operator-sdk generate csv --csv-version 0.0.1 --update-crds
INFO[0000] Generating CSV manifest version 0.0.1        
WARN[0000] Required csv fields not filled in file deploy/olm-catalog/gramola-operator/0.0.1/gramola-operator.v0.0.1.clusterserviceversion.yaml:
	spec.keywords
	spec.maintainers
	spec.provider 
```

We need to fix this **WARNING** but at the same time we'll also fix some others you would find along the way and why not add some useful information and even an icon.

We have to provide keywords, maintainers, provider... and also an icon ;-)

```
export CSV_PATH=./deploy/olm-catalog/gramola-operator/0.0.1/gramola-operator.v0.0.1.clusterserviceversion.yaml

yq w -i -s update_csv_instructions.yaml ${CSV_PATH}
```

All these changes (user defined) are reflected directly in the CSV generated for version 0.0.1 and will be kept if you re-generate the CSV for the same version over and over.

But there're still some changes to the CSV that need to be done in the type associated with our CRD `AppService`. This can be done adding some comments to the type code, comments starting with `+operator-sdk:gen-csv`.

For instance to add `spec.customresourcedefinitions.owned[].displayName` to the CSV indirectly we add `+operator-sdk:gen-csv:customresourcedefinitions.displayName="AppService"` as in this excerpt.

```go
...
// AppService is the Schema for the appservices API defines Gramola Backend Services
// +operator-sdk:gen-csv:customresourcedefinitions.displayName="AppService"
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=appservices,scope=Namespaced
type AppService struct {
...
```

Another example, imagine you want to add some UI related comments to surface some CRD realated data. For instance if the `AppService` is enabled and also the `Alias` chosen. Next example has this into account.

```go
  ...
  // Flags if the the AppService object is enabled or not
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Enabled"
  Enabled bool `json:"enabled"`
  ...
  // Different names for Gramola Service
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Alias"
	// +kubebuilder:validation:Enum=Gramola;Gramophone;Phonograph
  Alias string `json:"alias,omitempty"`
  ...
```

We have put all this together, so please substitute `./pkg/apis/gramola/v1alpha1/appservice_types.go` with these.

```go
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AppServiceSpec defines the desired state of AppService
type AppServiceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Flags if the the AppService object is enabled or not
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Enabled"
	Enabled bool `json:"enabled"`
	// Flags if the object has been initialized or not
	Initialized bool `json:"initialized,omitempty"`
	// Different names for Gramola Service
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Alias"
	// +kubebuilder:validation:Enum=Gramola;Gramophone;Phonograph
	Alias string `json:"alias,omitempty"`
}

// AppServiceStatus defines the observed state of AppService
type AppServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AppService is the Schema for the appservices API defines Gramola Backend Services
// +operator-sdk:gen-csv:customresourcedefinitions.displayName="AppService"
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=appservices,scope=Namespaced
type AppService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AppServiceSpec   `json:"spec,omitempty"`
	Status AppServiceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AppServiceList contains a list of AppService
type AppServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AppService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AppService{}, &AppServiceList{})
}

```


##```
##cat << EOF > update_instructions.yaml
##- command: update 
##  path: spec.customresourcedefinitions.owned
##  value:
##    - description: Gramola AppService definition
##      displayName: AppService
##      kind: AppService
##      name: appservices.gramola.redhat.com
##      version: v1alpha1
##EOF
##yq w -i -s update_instructions.yaml ${CSV_PATH}
##```




## Add /db and files...
TODO

## ENV VAR

```sh
export DB_SCRIPTS_BASE_DIR=$(pwd)
```

git add .
git commit -a -m "new"
git push origin master
git checkout -b 0.0.1
git push origin 0.0.1
git checkout master
git tag -a v0.0.1 -m "version 0.0.1"

# Installing operator-courier

The easiest way to build, validate and push Operator Artifacts is using the [`operator-courier`](https://github.com/operator-framework/operator-courier).

You can follow the instructions given [here](https://github.com/operator-framework/operator-courier) or use `virtualenv` to install it locally in your project:

1. Install `virtualenv` (unless already installed)

```sh
sudo pip install virtualenv
virtualenv -p python3 venv
printf "\nvenv/\n" >> .gitignore
source venv/bin/activate
```

2. Install `operator-courier`

```
(venv) $ pip3 install operator-courier
```

> **NOTE:** To leave the virtual env just deactivate it...
>
> `(venv) $ deactivate`

# Generate Quay token

We provide a script to generate the authetication token used later by `operator-courier` to push artifacts to *quay.io*.

```sh
$ sh gen_quayio_auth_token.sh
Username: jmanning
Password:
basic am1hbm5pbmc6ZXhhbXBsZXB3
```

Time to push your operator artifacts corresponding to CSV 0.0.1 to *quay.io*.

> Substitute <AUTH_TOKEN> with the token generated

```sh
operator-courier push ./deploy/olm-catalog/gramola-operator ${USERNAME} gramola-operator 0.0.1 "basic <AUTH_TOKEN>"
```

Test it

```
curl https://quay.io/cnr/api/v1/packages?namespace=cvicensa
[{"channels":null,"created_at":"2020-03-12T18:32:37","default":"0.0.1","manifests":["helm"],"name":"cvicensa/gramola-operator","namespace":"cvicensa","releases":["0.0.1"],"updated_at":"2020-03-12T18:32:37","visibility":"public"}]
```

# Go to Quay / Applications

Make your gramola-operator application public

# Linking the Operator Metadata from Quay into OpenShift
For OpenShift to become aware of the Quay application repository, an OperatorSource CR needs to be added to the cluster. Login to your OpenShift cluster as an admin (such as kubeadmin) and change to the openshift-marketplace project:

```
oc get opsrc -n openshift-marketplace
```

# Update the gramola-operatorsource to use your Quay USERNAME (if you are performing these steps on OSX, see note below)


```sh
$ sed -i "s|USERNAME|${USERNAME}|g" ./deploy/operator-source.yaml
```

On OSX use:

```sh
$ sed -i "" "s|USERNAME|${USERNAME}|g" ./deploy/operator-source.yaml
```

Now create your catalog resource:

```sh
oc apply -n openshift-marketplace -f ./deploy/operator-source.yaml
```

And check it has been created and refreshed properly.

```sh
oc get opsrc -n openshift-marketplace
NAME                  TYPE          ENDPOINT              REGISTRY              DISPLAYNAME           PUBLISHER   STATUS      MESSAGE                                       AGE
acme-operators        appregistry   https://quay.io/cnr   cvicensa              ACME Operators        ACME        Succeeded   The object has been successfully reconciled   64s
certified-operators   appregistry   https://quay.io/cnr   certified-operators   Certified Operators   Red Hat     Succeeded   The object has been successfully reconciled   2d6h
community-operators   appregistry   https://quay.io/cnr   community-operators   Community Operators   Red Hat     Succeeded   The object has been successfully reconciled   2d6h
redhat-operators      appregistry   https://quay.io/cnr   redhat-operators      Red Hat Operators     Red Hat     Succeeded   The object has been successfully reconciled   2d6h
```

If you go to the Operator Hub you should see `Other` category... etc.


# Preview your CSV

https://operatorhub.io/preview

# Deploy your operator


# Create some sample data



# Troubleshooting 
Have a look here

oc logs -f acme-operators-85cf48968d-9mgcg -n openshift-marketplace

Or here

oc logs -f catalog-operator-7fccd6877f-phh9p -n openshift-operator-lifecycle-manager 




https://redhat-connect.gitbook.io/certified-operator-guide/ocp-deployment/openshift-deployment
