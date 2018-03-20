# Permissions Health Checker

_permissions-checker_ performs a set of health checks in your web application to detect any possible issues with the permissions configuration.

## Installation

```
$> go get github.com/bitnami-labs/healthcheck-tools/cmd/permissions-checker/...
```

## Building from source

```
$> git clone https://github.com/bitnami-labs/healthcheck-tools.git
$> cd cmd/permissions-checker
$> make
```

## Basic usage

The tool is executed as follows:

```
$> permissions-checker -dir <DIRECTORY> -dir_default <DIR DEFAULT PERM> -file_default <FILE DEFAULT PERM> -owner <DEFAULT OWNER> -group <DEFAULT GROUP> -exclude <REGEXP> -show_hidden -verbose
```

The tool requires a set of parameters to work properly:

- *dir*: Directory to check. Default vale: */opt/bitnami*
- *dir_default*: Default directory permissions. Default vale: *rwxrwxr-x*
- *file_default*: Default file permissions. Default vale: *rwrw-r--*
- *owner*: Default owner. Default vale: *bitnami*
- *group*: Default group. Default vale: *daemon*
- *exclude*: Files and/or directories to e excluded.
- *show_hidden*: Show hidden files and directories
- *verbose*: Print every file and directory
- *version*: Show current version

## List of health checks

The tool will perform the following health checks:

- Check if the permissions are correct according to the default permissions.
- Check if the user owner is correct according to the default owner.
- Check if the group owner is correct according to the default group.

## Useful links

- [Troubleshoot Permissions issues (Bitnami Documentation pages)](https://docs.bitnami.com/general/how-to/troubleshoot-permission-issues/).
