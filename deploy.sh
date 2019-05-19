#!/bin/bash

DNS_DEPLOYMENT=$(kubectl get deploy -n kube-system -l k8s-app=kube-dns -o=custom-columns=NAME:.metadata.name | tail -n1 2>/dev/null)
if [ -z ${DNS_DEPLOYMENT} ]; then
    echo "could not detect DNS deployment for K8s cluster. Can not install istio-discovery"
    return
fi
echo patching "${DNS_DEPLOYMENT}"
kubectl patch deploy -n kube-system ${DNS_DEPLOYMENT} -p "`<kubernetes/deploy_patch.yaml`"
kubectl patch svc -n kube-system kube-dns -p "`<kubernetes/service_patch.yaml`"
kubectl patch clusterrole system:${DNS_DEPLOYMENT} -p "`sed -e 's/#DEPLOYMENT#/'"${DNS_DEPLOYMENT}"'/g' kubernetes/clusterrole.yaml`"
