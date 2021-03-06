## ignite ps

List running VMs

### Synopsis


List all running VMs. By specifying the all flag (-a, --all),
also list VMs that are not currently running.
Using the -f (--filter) flag, you can give conditions VMs should fullfilled to be displayed.
You can filter on all the underlying fields of the VM struct, see the documentation:
https://ignite.readthedocs.io/en/stable/api/ignite_v1alpha2.html#VM.

Different operators can be used:
- "=" and "==" for the equal
- "!=" for the is not equal
- "=~" for the contains
- "!~" for the not contains

Non-exhaustive list of identifiers to apply filter on:
- the VM name
- CPUs usage
- Labels
- Image
- Kernel
- Memory

Example usage:
	$ ignite ps -f "{{.ObjectMeta.Name}}=my-vm2,{{.Spec.CPUs}}!=3,{{.Spec.Image.OCI}}=~weaveworks/ignite-ubuntu"

	$ ignite ps -f "{{.Spec.Memory}}=~1024,{{.Status.Running}}=true"


```
ignite ps [flags]
```

### Options

```
  -a, --all             Show all VMs, not just running ones
  -f, --filter string   Filter the VMs
  -h, --help            help for ps
```

### Options inherited from parent commands

```
      --log-level loglevel      Specify the loglevel for the program (default info)
      --network-plugin plugin   Network plugin to use. Available options are: [cni docker-bridge] (default cni)
  -q, --quiet                   The quiet mode allows for machine-parsable output by printing only IDs
      --runtime runtime         Container runtime to use. Available options are: [docker containerd] (default containerd)
```

### SEE ALSO

* [ignite](ignite.md)	 - ignite: easily run Firecracker VMs

