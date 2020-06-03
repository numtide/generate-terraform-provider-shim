package main

import (
	"bytes"
	"encoding/hex"
	"text/template"

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

func generateShim(downloadURL, pluginName, version, binaryName string, hash []byte) (string, error) {
	t, err := template.New("shim").Parse(shimTemplate)
	if err != nil {
		return "", errors.Wrap(err, "while parsing shim template")
	}

	bb := new(bytes.Buffer)
	err = t.Execute(bb, templateData{
		DownloadURL: downloadURL,
		PluginName:  pluginName,
		Version:     version,
		BinaryName:  binaryName,
		SHA1:        hex.EncodeToString(hash[:]),
	})

	if err != nil {
		return "", errors.Wrap(err, "while rendering template")
	}

	return bb.String(), nil

}
