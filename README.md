# Installing the Operator SDK
export RELEASE_VERSION=v0.15.1

mkdir ./bin

curl -LO https://github.com/operator-framework/operator-sdk/releases/download/${RELEASE_VERSION}/operator-sdk-${RELEASE_VERSION}-x86_64-linux-gnu

mv operator-sdk-${RELEASE_VERSION}-x86_64-linux-gnu ./bin/operator-sdk

chmos u+x ./bin/operator-sdk

export PATH=$PATH:$(pwd)/bin
 
# Set up environment
cd $GOPATH
export GO111MODULE=on

export OPERATOR_NAME="gramola-operator"
export API_VERSION="gramola.redhat.com/v1alpha1"

export PROJECT_NAME=${OPERATOR_NAME}-project

# Create an ${OPERATOR_NAME} project that defines the App CR.
mkdir -p $GOPATH/src/github.com/redhat

cd $GOPATH/src/github.com/redhat

operator-sdk new ${OPERATOR_NAME} --type=go --repo github.com/redhat/${OPERATOR_NAME}

cd ./${OPERATOR_NAME}

# Add a new API for the custom resource AppService

> IF error creating API ==> export GOROOT=$(go env GOROOT)

$ operator-sdk add api --api-version=${API_VERSION} --kind=AppService

# Add a new controller that watches for AppService

$ operator-sdk add controller --api-version=${API_VERSION} --kind=AppService


# Edit the CR

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

## Regenerate supporting code for your CRDs
operator-sdk generate k8s

## Regenerate your CRD OpenAPI definition [deprecated]
operator-sdk generate openapi

## Regenerate your CRD OpenAPI definition
operator-sdk generate crds

> **Build the latest openapi-gen from source!**
> `which ./bin/openapi-gen > /dev/null || go build -o ./bin/openapi-gen k8s.io/kube-openapi/cmd/openapi-gen`

## Run openapi-gen for each of your API group/version packages
./bin/openapi-gen --logtostderr=true -o "" -i ./pkg/apis/gramola/v1alpha1 -O zz_generated.openapi -p ./pkg/apis/gramola/v1alpha1 -h ./hack/boilerplate.go.txt -r "-"

# If you import a new module run this

$ go mod vendor

# List module versions
go list -m -versions gopkg.in/src-d/go-git.v4

# Set the username variable

$ export USERNAME=<username>

# Build and push the ${OPERATOR_NAME} image to a public registry such as quay.io
$ export OPERATOR_VERSION=0.0.1
$ operator-sdk build quay.io/${USERNAME}/${OPERATOR_NAME}:${OPERATOR_VERSION}

# Login to public registry such as quay.io
$ docker login quay.io

# Push image
$ docker push quay.io/${USERNAME}/${OPERATOR_NAME}:${OPERATOR_VERSION}

# Update the operator manifest to use the built image name (if you are performing these steps on OSX, see note below)
$ sed -i "s|REPLACE_IMAGE|quay.io/${USERNAME}/${OPERATOR_NAME}\:${OPERATOR_VERSION}|g" deploy/operator.yaml

# On OSX use:
$ sed -i "" "s|REPLACE_IMAGE|quay.io/${USERNAME}/${OPERATOR_NAME}\:${OPERATOR_VERSION}|g" deploy/operator.yaml

# Set/Create project

$ oc new-project ${PROJECT_NAME}

or

$ oc project ${PROJECT_NAME}

# Setup Service Account
$ oc apply -f deploy/service_account.yaml
# Setup RBAC
$ oc apply -f deploy/role.yaml
$ oc apply -f deploy/role_binding.yaml
# Setup the CRD
$ oc apply -f deploy/crds/gramola.redhat.com_appservices_crd.yaml

# Run locally
$ operator-sdk run --local --namespace ${PROJECT_NAME}

# Or Deploy the ${OPERATOR_NAME}
$ oc apply -f deploy/operator.yaml

# Create an AppService CR
# The default controller will watch for AppService objects and create a pod for each CR
$ oc apply -f deploy/crds/gramola.redhat.com_v1alpha1_appservice_cr.yaml

You should get this error:

```
oc apply -f deploy/crds/gramola.redhat.com_v1alpha1_appservice_cr.yaml 
The AppService "example-appservice" is invalid: 
* spec.enabled: Required value
* spec.initialized: Required value
```

Change the type, so that initialized is not required

```go
Initialized bool `json:"initialized,omitempty"`
```

Let's change the CR so that enabled is defined:

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

```sh
$ operator-sdk generate csv --csv-version 0.0.1
INFO[0000] Generating CSV manifest version 0.0.1        
WARN[0000] Required csv fields not filled in file deploy/olm-catalog/gramola-operator/0.0.1/gramola-operator.v0.0.1.clusterserviceversion.yaml:
	spec.keywords
	spec.maintainers
	spec.provider 
```

We have to provide keywords, maintainers, provider... and also an icon ;-)

```
yq d -i ./deploy/olm-catalog/gramola-operator/0.0.1/gramola-operator.v0.0.1.clusterserviceversion.yaml spec.provider
yq w -i ./deploy/olm-catalog/gramola-operator/0.0.1/gramola-operator.v0.0.1.clusterserviceversion.yaml spec.provider.name "ACME Inc."
yq w -i ./deploy/olm-catalog/gramola-operator/0.0.1/gramola-operator.v0.0.1.clusterserviceversion.yaml spec.keywords[+] "gramola"
yq w -i ./deploy/olm-catalog/gramola-operator/0.0.1/gramola-operator.v0.0.1.clusterserviceversion.yaml spec.keywords[+] "backend"
yq w -i ./deploy/olm-catalog/gramola-operator/0.0.1/gramola-operator.v0.0.1.clusterserviceversion.yaml spec.icon[+].base64data $(cat gramola.svg | base64)
yq w -i ./deploy/olm-catalog/gramola-operator/0.0.1/gramola-operator.v0.0.1.clusterserviceversion.yaml spec.icon[0].mediatype image/svg+xml
yq w -i ./deploy/olm-catalog/gramola-operator/0.0.1/gramola-operator.v0.0.1.clusterserviceversion.yaml spec.maintainers[+].email "admin@gramola.com"
yq w -i ./deploy/olm-catalog/gramola-operator/0.0.1/gramola-operator.v0.0.1.clusterserviceversion.yaml spec.maintainers[0].name "ACME Inc."

```



https://redhat-connect.gitbook.io/certified-operator-guide/ocp-deployment/openshift-deployment
