This tool is extracted from tanka.
It is equivalent to `helmvendor`, with all of the jsonnet specific functionality removed.
Original source: https://github.com/grafana/tanka
Documentation also extracted from https://tanka.dev/helm

# Helmvendor

Helm does not make vendoring incredibly easy by itself. helm pull provides the required plumbing, but it does not record its actions in a reproducible manner.

# Create a chartfile.yaml in the current directory, e.g. in lib/myLibrary
$ helmvendor init

$ # Install the MySQL chart at version 1.6.7 from the stable repository
$ helmvendor add stable/mysql@1.6.7

Adding charts: To add a chart, use the following:

$ helmvendor add <repo>/<name>@<version>

This will also call helmvendor vendor, so that the charts/ directory is updated.

Adding Repositories: By default, the stable repository is automatically set up for you. If you wish to add another repository, you can use the add-repo command:

# Add the official Grafana repository
$ helmvendor charts add-repo grafana https://grafana.github.io/helm-charts

Another way is to modify chartfile.yaml directly:

version: 1
repositories:
  - name: stable
    url: https://charts.helm.sh/stable
+ - name: grafana
+   url: https://grafana.github.io/helm-charts

Installing multiple versions of the same chart: If you wish to install multiple versions of the same chart, you can write them to a specific directory.
You can do so with a :<directory> suffix in the add command, or by modifying the chartfile manually.

helmvendor add stable/mysql@1.6.7:1.6.7
helmvendor add stable/mysql@1.6.8:1.6.8

The resulting chartfile will look like this:

version: 1
directory: charts
repositories:
- name: stable
  url: https://charts.helm.sh/stable
requires:
- chart: stable/mysql
  directory: 1.6.7
  version: 1.6.7
- chart: stable/mysql
  directory: 1.6.8
  version: 1.6.8

Install charts from chartfile: To install charts from an existing chartfile, use the following:

$ helmvendor vendor

Optionally, you can also pass the --prune flag to remove vendored charts that are no longer in the chartfile.
OCI Registry Support

Tanka supports pulling charts from OCI registries. To use one, the chart name must be split into two parts: the registry and the chart name.

As example, if you wanted to pull the oci://public.ecr.aws/karpenter/karpenter:v0.27.3 image, your chartfile would look like this:

version: 1
directory: charts
repositories:
- name: karpenter
  url: oci://public.ecr.aws/karpenter
requires:
- chart: karpenter/karpenter
  directory: v0.27.3
  version: v0.27.3

Registry login is not supported yet.
Troubleshooting
Helm executable missing

Helm support in Tanka requires the helm binary installed on your system and available on the $PATH. If Helm is not installed, you will see this error message:

evaluating jsonnet: RUNTIME ERROR: Expanding Helm Chart: exec: "helm": executable file not found in $PATH

To solve this, you need to install Helm. If you cannot install it system-wide, you can point Tanka at your executable using HELM_PATH
opts.calledFrom unset

This occurs, when Tanka was not told where it helm.template() was invoked from. This most likely means you didn't call new(std.thisFile) when importing tanka-util:

local tanka = import "github.com/grafana/jsonnet-libs/tanka-util/main.libsonnet";
local helm = tanka.helm.new(std.thisFile);
                       ↑ This is important

Failed to find Chart

helmTemplate: Failed to find a Chart at 'stable/grafana': No such file or directory.
helmTemplate: Failed to find a Chart at '/home/user/stuff/tanka/environments/default/grafana': No such file or directory.

Tanka failed to locate your Helm chart on the filesystem. It looked at the relative path you provided in helm.template(), starting from the directory of the file you called helm.template() from.

Please check there is actually a valid Helm chart at this place. Referring to charts as <repo>/<name> is disallowed by design.
Two resources share the same name

To make customization easier, helm.template() returns the resources not as the list it receives from Helm, but instead converts this into an object.

For the indexing key it uses kind_name by default. In some rare cases, this might not be enough to distinguish between two resources, namely when the same resource exists in two namespaces.

To handle this, pass a custom name format, e.g. to also include the namespace:

custom: helm.template('foo', './charts/foo', {
  nameFormat: '{{ print .namespace "_" .kind "_" .metadata.name | snakecase }}'
})

The literal default format used is {{ print .kind "_" .metadata.name | snakecase }}