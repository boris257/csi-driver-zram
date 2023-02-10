module github.com/boris257/csi-driver-zram

go 1.18

require (
	github.com/container-storage-interface/spec v1.7.0
	github.com/golang/protobuf v1.5.2
	github.com/hashicorp/go-multierror v1.1.1
	github.com/kubernetes-csi/csi-lib-utils v0.12.0
	github.com/stretchr/testify v1.8.1
	golang.org/x/net v0.5.0
	golang.org/x/sys v0.4.0
	google.golang.org/grpc v1.52.0
	k8s.io/apimachinery v0.26.0
	k8s.io/klog/v2 v2.80.1
	k8s.io/mount-utils v0.25.6
	k8s.io/utils v0.0.0-20221128185143-99ec85e7a448
	sigs.k8s.io/yaml v1.3.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/moby/sys/mountinfo v0.6.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	google.golang.org/genproto v0.0.0-20221118155620-16455021b5e6 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.25.6
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.25.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.25.6
	k8s.io/apiserver => k8s.io/apiserver v0.25.6
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.25.6
	k8s.io/client-go => k8s.io/client-go v0.25.6
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.25.6
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.25.6
	k8s.io/code-generator => k8s.io/code-generator v0.25.6
	k8s.io/component-base => k8s.io/component-base v0.25.6
	k8s.io/component-helpers => k8s.io/component-helpers v0.25.6
	k8s.io/controller-manager => k8s.io/controller-manager v0.25.6
	k8s.io/cri-api => k8s.io/cri-api v0.25.6
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.25.6
	k8s.io/gengo => k8s.io/gengo v0.0.0-20200114144118-36b2048a9120
	k8s.io/heapster => k8s.io/heapster v1.2.0-beta.1
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.25.6
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.25.6
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.25.6
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.25.6
	k8s.io/kubectl => k8s.io/kubectl v0.25.6
	k8s.io/kubelet => k8s.io/kubelet v0.25.6
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.25.6
	k8s.io/metrics => k8s.io/metrics v0.25.6
	k8s.io/mount-utils => k8s.io/mount-utils v0.25.6
	k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.25.6
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.25.6
	k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.25.6
	k8s.io/sample-controller => k8s.io/sample-controller v0.25.6
	k8s.io/system-validators => k8s.io/system-validators v1.0.4
	sigs.k8s.io/yaml => sigs.k8s.io/yaml v1.2.0
)
