# bundle-install

## `gcr.io/paketo-buildpacks/bundle-install`

A Cloud Native Buildpack to install gems from a Gemfile

This will be providing `gems`.

## Packaging

To package this buildpack for consumption:

```bash
./scripts/package.sh --version 0.8.18
```

This will build the buildpack for all target architectures specified in `buildpack.toml` (amd64 and arm64 by default) and create architecture-specific archives in the `build/` directory.

## Publishing

To publish this buildpack to a registry:

```bash
./scripts/publish.sh \
  --image-ref 348674388966.dkr.ecr.us-east-1.amazonaws.com/neeto-deploy/paketo/buildpack/bundle-install:0.8.18
```

The script will automatically:
- Read target architectures from `buildpack.toml`
- Detect architecture-specific buildpack archives
- Publish each architecture separately
- Create and push a multi-arch manifest list
