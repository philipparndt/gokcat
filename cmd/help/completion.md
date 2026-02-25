# gokcat completion — generate shell completion scripts

Generate shell completion scripts for your shell. Supports bash, zsh, fish, and PowerShell.

## Usage

```
gokcat completion [bash|zsh|fish|powershell]
```

## Bash

```bash
source <(gokcat completion bash)
```

To load completions for each session, execute once:

```bash
# Linux:
gokcat completion bash > /etc/bash_completion.d/gokcat

# macOS:
gokcat completion bash > $(brew --prefix)/etc/bash_completion.d/gokcat
```

## Zsh

If shell completion is not already enabled in your environment, enable it first:

```zsh
echo "autoload -U compinit; compinit" >> ~/.zshrc
```

To load completions for each session, execute once:

```zsh
gokcat completion zsh > "${fpath[1]}/_gokcat"
```

Start a new shell for the setup to take effect.

## fish

```fish
gokcat completion fish | source
```

To load completions for each session, execute once:

```fish
gokcat completion fish > ~/.config/fish/completions/gokcat.fish
```

## PowerShell

```powershell
gokcat completion powershell | Out-String | Invoke-Expression
```

To load completions for every new session, run:

```powershell
gokcat completion powershell > gokcat.ps1
```

Then source this file from your PowerShell profile.
