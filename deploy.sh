#!/bin/bash

DNS_DEPLOYMENT=$(kubectl get deploy -n kube-system -l k8s-app=kube-dns -o=custom-columns=NAME:.metadata.name | tail -n1 2>/dev/null)
[[ -z ${DNS_DEPLOYMENT} ]] && echo "could not detect DNS deployment for K8s cluster. Can not install istio-discovery" && exit 1

echo -e "Patching ${DNS_DEPLOYMENT} deployment"
if [[ $(kubectl patch deploy -n kube-system ${DNS_DEPLOYMENT} -p "$(cat kubernetes/deploy_patch.yaml)" >/dev/null 2>&1 )$? != 0 ]]; then
   echo -e "Error patching deployment ${DNS_DEPLOYMENT}" && exit 1
fi

echo -e "Patching ${DNS_DEPLOYMENT} service"
if [[ $(kubectl patch svc -n kube-system kube-dns -p "$(cat kubernetes/service_patch.yaml)" >/dev/null 2>&1 )$? != 0 ]]; then
   echo -e "Error patching service kube-dns" && exit 1
fi

echo -e "Patching system:${DNS_DEPLOYMENT} cluster-role"
if [[ $(kubectl patch clusterrole system:${DNS_DEPLOYMENT} -p "$(cat kubernetes/clusterrole.yaml)" >/dev/null 2>&1 )$? != 0 ]]; then
   echo -e "Error patching clusterrole system:${DNS_DEPLOYMENT}" && exit 1
fi



