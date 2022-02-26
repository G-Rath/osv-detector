# Open Source Vulnerability Detector

An auditing tool for detecting vulnerabilities using the Open Source
Vulnerability specification

## Usage

The detector maintains a cache of the OSV database it's using locally, which is
updated with any changes at the start of every run.

You can have the detector work solely off this cache with the `--offline` flag:

```shell
osv-detector --offline
```

This requires the detector to have successfully cached an offline copy of the
OSV database at least once.
