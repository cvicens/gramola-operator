# Demo Script

This demo...

# Git clone your repo

No need to repeat this step if already cloned.

# Update USERNAME settings.sh to match your quay.io account

...

# Log in quay.io with your USERNAME

docker login ...

# Build version 0.0.1

Checkout branch 0.0.1

===> TODO

# Build version 0.0.2

Checkout branch 0.0.2

===> TODO

# Delete previous applications in your quai.io account

Go to [quay.io](https://quay.io) ... delete

# Delete previous AppService objects and uninstall Gramola Operator

...

# Push version 0.0.1

Change `currentCSV` to `v0.0.1` in the package manifest `./deploy/olm-catalog/gramola-operator/gramola-operator.package.yaml`

```yaml
channels:
- currentCSV: gramola-operator.v0.0.1
  name: alpha
defaultChannel: alpha
packageName: gramola-operator
```

Now we're ready to push version 0.0.1

./push-csv-0.0.1.sh

Check that version 0.0.1 has been uploaded to quay.io

Make it public! As you know this is only to make the demo easier!

# To speed up things... let's re-create the catalog source

```sh
oc delete -n openshift-marketplace -f ./deploy/operator-source.yaml
```

```sh
oc apply -n openshift-marketplace -f ./deploy/operator-source.yaml
```

# Deploy version 0.0.1 of our operator

Create a project and go to the operator hub -> Other

Install version 0.0.1

Create the example AppService CR

# Create some events

./create-sample-data.sh

# Deploy version 0.0.2

Change `currentCSV` to `v0.0.2` in the package manifest `./deploy/olm-catalog/gramola-operator/gramola-operator.package.yaml`

```yaml
channels:
- currentCSV: gramola-operator.v0.0.2
  name: alpha
defaultChannel: alpha
packageName: gramola-operator
```

Now we're ready to push version 0.0.2

./push-csv-0.0.2.sh

Check that version 0.0.2 has been uploaded to quay.io

# To speed up things... let's re-create the catalog source

```sh
oc delete -n openshift-marketplace -f ./deploy/operator-source.yaml
```

```sh
oc apply -n openshift-marketplace -f ./deploy/operator-source.yaml
```

 oc get opsrc -n openshift-marketplace
NAME                  TYPE          ENDPOINT              REGISTRY              DISPLAYNAME           PUBLISHER   STATUS      MESSAGE                                       AGE
acme-operators        appregistry   https://quay.io/cnr   cvicensa              ACME Operators        ACME        Succeeded   The object has been successfully reconciled   12s
certified-operators   appregistry   https://quay.io/cnr   certified-operators   Certified Operators   Red Hat     Succeeded   The object has been successfully reconciled   4d14h
community-operators   appregistry   https://quay.io/cnr   community-operators   Community Operators   Red Hat     Succeeded   The object has been successfully reconciled   4d14h
redhat-operators      appregistry   https://quay.io/cnr   redhat-operators      Red Hat Operators     Red Hat     Succeeded   The object has been successfully reconciled   4d14h

# Observe how the operator is updated

...

# Observe the events related to the AppService instance

...

# Test the application and see the new and ols interfaces working

... old tab
... new tab

So data was migrated along with the code!