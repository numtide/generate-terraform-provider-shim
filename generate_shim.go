package main

import (
	"strings"

	"github.com/numtide/generate-terraform-provider-shim/shim/targz"
	zipshim "github.com/numtide/generate-terraform-provider-shim/shim/zip"
	"github.com/pkg/errors"
)

var shimTemplate = `#!/usr/bin/env bash
set -e -o pipefail

plugin_url="{{.DownloadURL}}"
plugin_unpack_dir="${XDG_CACHE_HOME:-$HOME/.cache}/terraform-providers/{{.PluginName}}_v{{.Version}}"
plugin_binary_name="{{.BinaryName}}"
plugin_binary_path="${plugin_unpack_dir}/${plugin_binary_name}"
plugin_binary_sha1="{{.SHA1}}"

if [[ ! -d "${plugin_unpack_dir}" ]]; then
    mkdir -p "${plugin_unpack_dir}"
fi

if [[ -f "${plugin_binary_path}" ]]; then
    current_sha=$(git hash-object "${plugin_binary_path}")
    if [[ $current_sha != "${plugin_binary_sha1}" ]]; then
        rm "${plugin_binary_path}"
    fi
fi

if [[ ! -f "${plugin_binary_path}" ]]; then
    curl -L "${plugin_url}" | tar xzvfC - "${plugin_unpack_dir}"
    chmod 755 "${plugin_binary_path}"
fi

current_sha=$(git hash-object "${plugin_binary_path}")
if [[ $current_sha != "${plugin_binary_sha1}" ]]; then
    echo "plugin binary sha does not match ${current_sha} != ${plugin_binary_sha1}" >&2
    exit 1
fi

exec "${plugin_binary_path}" $@
`

type templateData struct {
	DownloadURL string
	PluginName  string
	Version     string
	BinaryName  string
	SHA1        string
}

func generateShim(downloadURL, pluginName, version, binaryName string) (string, error) {

	if strings.HasSuffix(downloadURL, ".tar.gz") {
		return targz.GenerateShim(downloadURL, pluginName, version, binaryName)
	}

	if strings.HasSuffix(downloadURL, ".zip") {
		return zipshim.GenerateShim(downloadURL, pluginName, version, binaryName)
	}

	return "", errors.New("cannot generate shim")

}
