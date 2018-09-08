# slow-csi

This CSI (Container Storage Interface) plugin does nothing but `time.Sleep()`. It simulates pending operations with configurable delays, useful for debugging, triggering timeouts, etc.

Delayable RPCs:
* `CreateVolume`
* `DeleteVolume`
* `NodeStageVolume`
* `NodeUnstageVolume`
* `NodePublishVolume`
* `NodeUnpublishVolume`

## Building

slow-csi plugin can be compiled in a form of a binary file or in a form of a Docker image. When compiled as a binary file, the result is stored in `_output/` directory with the name `slowplugin`. When compiled as an image, it's stored in the local Docker image store.

Building binary:
```bash
$ make slowplugin
```

Building Docker image:
```bash
$ make image
```

## Configuration

**Available command line arguments:**

Option | Default value | Description
------ | ------------- | -----------
`--endpoint` | `unix://tmp/csi.sock` | CSI endpoint, must be a UNIX socket
`--drivername` | `csi-slowplugin` | name of the driver (Kubernetes: `provisioner` field in StorageClass must correspond to this value)
`--nodeid` | _empty_ | This node's ID
`--defaultdelay` | _empty_ | Default delay for all delayable RPCs. Defaults to 0
`--delay` | _empty_ | Per-RPC delay configuration in format `RPC1=DELAY,RPC2=DELAY,...`

### Delay format

All delays are in seconds and may be in one of three formats:
* `n`: simple number
* `a..b`: random delay in an interval _[`a`, `b`)_, `b` must be greater than `a`
* `inf`: infinite delay, such RPC will never finish

Example: `--delay=CreateVolume=10..20,DeleteVolume=5,NodeUnstageVolume=inf`

## Deployment with Kubernetes

Requires Kubernetes 1.11

Your Kubernetes cluster must allow privileged pods (i.e. `--allow-privileged` flag must be set to true for both the API server and the kubelet). Moreover, as stated in the [mount propagation docs](https://kubernetes.io/docs/concepts/storage/volumes/#mount-propagation), the Docker daemon of the cluster nodes must allow shared mounts.

YAML manifests are located in `deploy/kubernetes`.

**Deploy RBACs for sidecar containers and node plugins:**

```bash
$ kubectl create -f csi-attacher-rbac.yaml
$ kubectl create -f csi-provisioner-rbac.yaml
$ kubectl create -f csi-nodeplugin-rbac.yaml
```

Those manifests deploy service accounts, cluster roles and cluster role bindings.

**Deploy CSI sidecar containers:**

```bash
$ kubectl create -f csi-slowplugin-attacher.yaml
$ kubectl create -f csi-slowplugin-provisioner.yaml
```

Deploys stateful sets for external-attacher and external-provisioner sidecar containers for slow-csi.

**Deploy slow-csi driver:**

```bash
$ kubectl create -f csi-slowplugin.yaml
```

Deploys a daemon set with two containers: CSI driver-registrar and the slow-csi driver. Here you can configure the delays.

Once the driver is successfully deployed, you can deploy the storage class, PVC and a deployment, no changes required. YAML manifests for those are located in `run/`:

```bash
$ kubectl create -f storageclass.yaml
$ kubectl create -f pvc.yaml
$ kubectl create -f deployment.yaml
```