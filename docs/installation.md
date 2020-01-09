# Ignite installation guide

This guide describes the installation and uninstallation process of Ignite.

## System requirements

Ignite runs on most Intel, AMD or ARM (AArch64) based `linux/amd64` systems with `KVM` support.
See the full CPU support table in [dependencies.md](dependencies.md) for more information.

See [cloudprovider.md](cloudprovider.md) for guidance on running Ignite on various cloud providers and suitable instances that you could use.

**NOTE:** You do **not** need to install any "traditional" QEMU/KVM packages, as long as
there is virtualization support in the CPU and kernel it works.

See [dependencies.md](dependencies.md) for needed dependencies.

### Checking for KVM support

Please read [dependencies.md](dependencies.md) for the full reference, but if you quickly want
to check if your CPU and kernel supports virtualization, run these commands:

```console
$ lscpu | grep Virtualization
Virtualization:      VT-x

$ lsmod | grep kvm
kvm_intel             200704  0
kvm                   593920  1 kvm_intel
```

Alternatively, on Ubuntu-like systems there's a tool called `kvm-ok` in the `cpu-checker` package.
Check for KVM support using `kvm-ok`:

```console
$ sudo apt-get update && sudo apt-get install -y cpu-checker
...
$ kvm-ok
INFO: /dev/kvm exists
KVM acceleration can be used
```

With this kind of output, you're ready to go!

## Installing dependencies

Ignite has a few dependencies (read more in this [doc](dependencies.md)).
Install them on Ubuntu/CentOS like this:  
(Ignite does not depend on docker package version. If you already installed docker-ce, you don't need to replace it to docker.io.)

Ubuntu:

```bash
apt-get update && apt-get install -y --no-install-recommends dmsetup openssh-client git binutils
which containerd || apt-get install -y --no-install-recommends containerd
    # Install containerd if it's not present -- prevents breaking docker-ce installations
```

CentOS:

```bash
yum install -y e2fsprogs openssh-clients git
which containerd || ( yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo && yum install -y containerd.io )
    # Install containerd if it's not present
```

### CNI Plugins

Install the CNI binaries like this:

```shell
export CNI_VERSION=v0.8.2
export ARCH=$([ $(uname -m) = "x86_64" ] && echo amd64 || echo arm64)
mkdir -p /opt/cni/bin
curl -sSL https://github.com/containernetworking/plugins/releases/download/${CNI_VERSION}/cni-plugins-linux-${ARCH}-${CNI_VERSION}.tgz | tar -xz -C /opt/cni/bin
```

Note that the SSH and Git packages are optional; they are only needed if you use
the `ignite ssh` and/or `ignite gitops` commands.

## Downloading the binary

Ignite is a currently a single binary application. To install it,
download the binary from the [GitHub releases page](https://github.com/weaveworks/ignite/releases),
save it as `/usr/local/bin/ignite` and make it executable.

To install Ignite from the command line, follow these steps:

```bash
export VERSION=v0.6.3
export GOARCH=$(go env GOARCH 2>/dev/null || echo "amd64")

for binary in ignite ignited; do
    echo "Installing ${binary}..."
    curl -sfLo ${binary} https://github.com/weaveworks/ignite/releases/download/${VERSION}/${binary}-${GOARCH}
    chmod +x ${binary}
    sudo mv ${binary} /usr/local/bin
done
```

Ignite uses [semantic versioning](https://semver.org), select the version to be installed
by changing the `VERSION` environment variable.

## Verifying the installation

If the installation was successful, the `ignite` command should now be available:

```console
$ ignite version
Ignite version: version.Info{Major:"0", Minor:"6", GitVersion:"v0.6.3", GitCommit:"...", GitTreeState:"clean", BuildDate:"...", GoVersion:"...", Compiler:"gc", Platform:"linux/amd64"}
Firecracker version: v0.18.1
Runtime: containerd
```

Now you can continue with the [Getting Started Walkthrough](usage.md).

## Removing the installation

To completely remove the Ignite installation, execute the following as root:

```bash
# Force-remove all running VMs
ignite rm -f $(ignite ps -aq)
# Remove the data directory
rm -r /var/lib/firecracker
# Remove the ignite and ignited binaries
rm /usr/local/bin/ignite{,d}
```
