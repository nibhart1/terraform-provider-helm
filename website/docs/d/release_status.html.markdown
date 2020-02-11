---
layout: "helm"
page_title: "helm: helm_release_status"
sidebar_current: "docs-helm-release-status"
description: |-

---


# Data Source: release_status
Displays the status of the named release.

release_status describes the  status of pods,services,ingress and persistence volume claim running  inside the release deployed in a kubernetes cluster. It consist of  list of info of the running pods ,list of info of the running services,list of info of the running ingress,list of info of the persistence volume claim.  

# Example
``` 
resource "helm_release" "test" {
    name      = "test"
    namespace = "default"
    chart     = "stable/mariadb"
    version   = "0.6.2"

    wait = true
}
    
data "helm_release_status" "test" {
    name      = "test"
    revision   = 1
}
``` 

# Requirements


* You must have Terraform installed
* Helm should be configured properly


# Argument Reference
The following arguments are supported:


 * name  - (Required) Name of the Release
 * revision  - (Required) No of times release got revised or updated
 

# Attribute Reference

The output of release_status is the status of deployed release ,namespace of release , list status of individual resources and the count of each resources.

* namespace - namespace of the release.

* chart_status - status of the chart.

Each individual resource list supports:
 
 
* `pods`
  * `name`  - name of the pod.
  * `age` - age of the pod.


* `service`
  * `name`  - name of the service.
  * `age` - age of the service.

* `ingress`
  * `name`  - name of the ingress.
  * `age` - age of the ingress.

* `pvc`
  * `name`  - name of the persistence volume claim.
  * `age` - age of the persistence volume claim.

* `total_pods` - total number of pods existing
* `total_services` - total number of service existing
* `total_ingress` - total number of ingress existing
* `total_pvcs` - total number of persistence volume claim existing