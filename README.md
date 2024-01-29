This tool is extracted from tanka.

It is equivalent to `tk tool charts`, with all of the jsonnet specific functionality removed.

Original source: https://github.com/grafana/tanka

Documentation also extracted from https://tanka.dev/helm

# Helmvendor

Helm does not make vendoring incredibly easy by itself. helm pull provides the required plumbing, but it does not record its actions in a reproducible manner.

# Create a chartfile.yaml, install charts
Start by creating a chartfile using the init command:
```
$ helmvendor init
```

To add a chart, use the following:
```
$ helmvendor add <repo>/<name>@<version>
```

This will also call `helmvendor vendor``, so that the charts/ directory is updated.

For example, to install the MySQL chart at version 1.6.7 from the stable repository:
```
$ helmvendor add stable/mysql@1.6.7
```

# Add repository
Adding Repositories: By default, the stable repository is automatically set up for you. If you wish to add another repository, you can use the add-repo command:
```
$ helmvendor charts add-repo grafana https://grafana.github.io/helm-charts
```

Another way is to modify chartfile.yaml directly:
```
version: 1
repositories:
  - name: stable
    url: https://charts.helm.sh/stable
+ - name: grafana
+   url: https://grafana.github.io/helm-charts
```

Installing multiple versions of the same chart: If you wish to install multiple versions of the same chart, you can write them to a specific directory.
You can do so with a :<directory> suffix in the add command, or by modifying the chartfile manually.

```
helmvendor add stable/mysql@1.6.7:1.6.7
helmvendor add stable/mysql@1.6.8:1.6.8
```

The resulting chartfile will look like this:
```
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
```

# Install charts from chartfile
To install charts from an existing chartfile, use the following:
```
$ helmvendor vendor
```
Optionally, you can also pass the --prune flag to remove vendored charts that are no longer in the chartfile.
