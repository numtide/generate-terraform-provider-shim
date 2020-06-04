# Generate terraform provider shim

Terraform doesn't provide a mechanism to install non-official plugins. Instead, the user has a few [locations to place the pre-compiled binary](https://www.terraform.io/docs/extend/how-terraform-works.html#plugin-locations).

Instead of checking the binaries into the repository, this project generates bash trampoline scripts that fetch those plugins on demand. This is a liteweight solution to "install" the providers at the moment where terraform starts accessing them.

## Usage

running

```sh
$ generate-provider-shim draganm/terraform-provider-linuxbox
```

will generate provider shims for linuxbox provider in the `.terraform/plugins/` directory.


## Related projects

* [terraform-bundle](https://github.com/hashicorp/terraform/tree/master/tools/terraform-bundle) - generate a bundle with terraform and the providers for offline usage.
