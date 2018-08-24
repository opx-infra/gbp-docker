# Frequently Asked Questions

**Sections**

1. `git-buildpackage` Questions

## `git-buildpackage` Questions

1. Why did I receive the error `gbp:error: upstream/XXX is not a valid treeish`?

`buildpackage` expects `3.0 (quilt)` packages to have a tag in the form of `vX.Y.Z` when the changelog reads `X.Y.Z-A`. If the upstream tag format is different, modify it with a flag. For example, an upstream format of `upstream/X.Y.Z` would use this command.

```bash
gbp buildpackage --git-upstream-tag='upstream/%(version)s'

# with dbp, it looks like this
dbp build repo -- --git-upstream-tag='upstream/%(version)s'
```
