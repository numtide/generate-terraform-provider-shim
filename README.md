# Generate terraform provider shim

This utility will generate terraform provider shims for the community providers available in the terraform repository.

## Usage

running

```sh
$ generate-provider-shim draganm/terraform-provider-linuxbox
```

will generate provider shims for linuxbox provider in the `.terraform/plugins/` directory.


## Why shims and not binaries?

... TODO: explanation here ...