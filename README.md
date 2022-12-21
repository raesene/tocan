# tòcan - Token Request Kubeconfig creator

this is just a small program to use the Kubernetes Token Request API to create a kubeconfig file.

It connects to a cluster based on your currently active kubeconfig, and creates a kubeconfig with the token for the service account you specify (default: default). you can specify `expirationSeconds` (default: 3600), but the maximum allowed might be restricted by your distribution (E.G. EKS)

## Usage

If run with no parameters it will create a file called `default.kubeconfig` in the current directory which contains a token for the default service account in the default namespace.

If you want to create a kubeconfig for a different service account, you can specify the namespace and service account name with the `-namespace` and `-service-account` parameters.

## command line parameters

There are four command line parameters:

```
  -expiration-seconds int
        The expiration time of the token in seconds (default 3600)
  -namespace string
        The namespace to use for the token (default "default")
  -output-file string
        The name of the kubeconfig file to create
  -service-account string
        The service account to use for the token (default "default")
```




## Name

tòcan is scottish gaelic for token