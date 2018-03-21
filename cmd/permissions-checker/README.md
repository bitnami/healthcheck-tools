# Permissions Health Checker

_permissions-checker_ performs a set of health checks in your server or local machine to detect any possible issues with the permissions configuration comparing the current permissions, owner and groups with those that should have.

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
- *exclude*: File or directory to be excluded (you can use RegExp).
- *show_hidden*: Include hidden files and directories in the check.
- *verbose*: Print every file and directory in a hierarchical way.
- *version*: Show current version.

## List of health checks

The tool will perform the following health checks:

- Check if the permissions are correct according to the default permissions.
- Check if the user owner is correct according to the default owner.
- Check if the group owner is correct according to the default group.

## Examples

#### Verbose

```
==================================================
	PERMISSIONS CHECKS
==================================================
Starting checks with these parameters:
	- Directory to check: /Applications/wordpress-4.9.4-6/apache2
	- Default file permissions: -rw-r--r--
	- Default dir permissions: rwxrwxr-x
	- Owner: crhernandez
	- Group: admin
	- Exclude: (?!.*)
	- Show hidden: false
	- Verbose: true
==================================================

-- Checking permissions --
 (d) htdocs -rwxr-xr-x (expected rwxrwxr-x) crhernandez admin
    (f) 503.html -rw-r--r-- crhernandez admin
    (f) bitnami.css -rw-r--r-- crhernandez admin
    (f) favicon.ico -rw-r--r-- crhernandez admin
    (d) img -rwxr-xr-x (expected rwxrwxr-x) crhernandez admin
      (f) background.png -rw-r--r-- crhernandez admin
      (f) bitnami.png -rw-r--r-- crhernandez admin
      (f) header_bg.png -rw-r--r-- crhernandez admin
      (f) lampstack.png -rw-r--r-- crhernandez admin
      (f) lappstack.png -rw-r--r-- crhernandez admin
      (f) menu_bg.png -rw-r--r-- crhernandez admin
      (f) module_table_bottom.png -rw-r--r-- crhernandez admin
      (f) module_table_top.png -rw-r--r-- crhernandez admin
      (f) plain-background.png -rw-r--r-- crhernandez admin
```

#### No verbose

```
==================================================
	PERMISSIONS CHECKS
==================================================
Starting checks with these parameters:
	- Directory to check: /Applications/wordpress-4.9.4-6/apache2
	- Default file permissions: -rw-r--r--
	- Default dir permissions: rwxrwxr-x
	- Owner: crhernandez
	- Group: admin
	- Exclude: (?!.*)
	- Show hidden: false
	- Verbose: false
==================================================

-- Checking permissions --
(d) /Applications/wordpress-4.9.4-6/apache2/bin -rwxr-xr-x (expected rwxrwxr-x) crhernandez admin
(f) /Applications/wordpress-4.9.4-6/apache2/bin/ab -rwxr-xr-x (expected -rw-r--r--) crhernandez admin
(f) /Applications/wordpress-4.9.4-6/apache2/bin/ab.bin -rwxr-xr-x (expected -rw-r--r--) crhernandez admin
(f) /Applications/wordpress-4.9.4-6/apache2/bin/apachectl -rwxr-xr-x (expected -rw-r--r--) crhernandez admin
(f) /Applications/wordpress-4.9.4-6/apache2/bin/apxs -rwxr-xr-x (expected -rw-r--r--) crhernandez admin
```

## Useful links

- [Troubleshoot Permissions issues (Bitnami Documentation pages)](https://docs.bitnami.com/general/how-to/troubleshoot-permission-issues/).
