#!/bin/bash

deployment_base="${1}"

if [[ -z $deployment_base ]]; then
	deployment_base="../deploy/kubernetes"
fi

cd "$deployment_base" || exit 1

objects=(csi-attacher-rbac csi-provisioner-rbac csi-nodeplugin-rbac csi-slowplugin-attacher csi-slowplugin-provisioner csi-slowplugin)

for obj in ${objects[@]}; do
	kubectl create -f "./$obj.yaml"
done
