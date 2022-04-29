# getsum

## Overview

In many CI/CD scripts, we often see the pattern of:

```sh
curl <url> | bash
```

But this pattern places a lot of trust in the remote URL and can be the source of security issues:

- <https://about.codecov.io/security-update/>

At minimum, we should check provenance of the source we download (ex sget, etc.)!
However, true signed provenance relies on artifact providers to integrate and provide this information, but we're not there yet.

That's where getsum comes in! This tool takes a trust-on-first-use approach to at least guarantee that the remote artifact has not been tampered with since you first saw it.

## Getting started

```bash
$ getsum [config]
```

If `[config]` is not specified, the tool will default to `getsum.yaml`

When ran, the tool will fetch each file and verify the checksum.
If the checksums do not match, the tool will exit immediately with an error.

When ran with the --init flag, getsum will populate the config file with checksums
for the files.

## Configuration

`getsum` is configured through a yaml config file. This file contains a list of Modules,
which describe the files + checksums.

| Field    | Description                            |
| -------- | -------------------------------------- |
| url      | URL to fetch                           |
| checksum | checksum of the object to verify       |
| path     | output path of the file after download |

### Example

```yaml
modules:
  - checksum: 02acfe3bc856805c25d65bb620e414b98aa504c6194ff8e953ce169edfcc03c6
    url: https://github.com/google/ko/releases/download/v0.11.2/ko_0.11.2_Darwin_arm64.tar.gz
```

### Inspirations

- go.mod
- bazel workspaces
