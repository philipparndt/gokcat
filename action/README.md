# setup-gokcat

This GitHub Action downloads and installs the latest version of the gokcat binary, with support for caching to speed up subsequent runs.

## Usage

```yaml
- name: Setup gokcat
  uses: ./action
  with:
    install-dir: /usr/local/bin # optional
```

## Inputs
- `install-dir`: Directory to install gokcat binary (default: `/usr/local/bin`).

## Outputs
- `gokcat-path`: Path to the installed gokcat binary.

## Development

- Install dependencies: `npm install`
- Build: `npm run build`
- Package for GitHub Actions: `npm run prepare`

## Notes
- The action uses caching to avoid repeated downloads.
- Only Linux binaries are provided by gokcat.
