# Installing and Uninstalling Verge CLI

Verge is distributed as a single standalone compiled binary for Linux, macOS, and Windows. This document details how to install the CLI on your system, specify versions, install manually, and perform a complete uninstallation.

---

## 1. Install Using the Automated Script

The fastest way to install the latest version of Verge is using our official shell installer script.

### Installation
Run the following pipe-chainable command in your terminal:

```bash
curl -sSL https://raw.githubusercontent.com/armckinney/verge/main/install.sh | bash
```

The script automatically:
* Detects your operating system (Linux, macOS, or Windows via Git-Bash/MSYS/WSL) and CPU architecture (`amd64` or `arm64`).
* Queries GitHub for the latest release.
* Downloads the appropriate archive, extracts the binary, and places it in an executable directory (e.g., `/usr/local/bin` or `~/.local/bin`).
* Elevates permissions via `sudo` only if required to write to the destination directory.

---

## 2. Install a Specific Version Using the Script

If you need to pin to a specific version of Verge (for example, in CI/CD pipelines or matching a team environment), you can pass the `VERGE_VERSION` environment variable before executing the script.

```bash
curl -sSL https://raw.githubusercontent.com/armckinney/verge/main/install.sh | VERGE_VERSION=0.1.5 bash
```

> [!NOTE]
> The `VERGE_VERSION` variable supports tags with or without the `v` prefix (e.g., `0.1.5` or `v0.1.5`).

---

## 3. Manual Installation (curl, tar, mv)

If you prefer to download and install the binary manually without running an automated bash script, follow these steps:

### Step 1: Formulate the Download URL
Select the version you want (e.g., `0.1.5`) and match your OS (`linux`, `darwin`, `windows`) and CPU architecture (`amd64`, `arm64`).

```
https://github.com/armckinney/verge/releases/download/v[VERSION]/verge_[VERSION]_[OS]_[ARCH].tar.gz
```

*Example for macOS Apple Silicon (Darwin ARM64) v0.1.5:*
`https://github.com/armckinney/verge/releases/download/v0.1.5/verge_0.1.5_darwin_arm64.tar.gz`

### Step 2: Download the Archive
```bash
curl -sSL -O "https://github.com/armckinney/verge/releases/download/v0.1.5/verge_0.1.5_darwin_arm64.tar.gz"
```

### Step 3: Extract the Binary
```bash
tar -xzf verge_0.1.5_darwin_arm64.tar.gz
```

### Step 4: Move to Execution Path
Move the extracted `verge` binary to a directory on your system's `PATH` and ensure it has execution permissions:

* **For Global Unix/macOS installation:**
  ```bash
  sudo mv verge /usr/local/bin/
  sudo chmod +x /usr/local/bin/verge
  ```
* **For Local User installation (no sudo required):**
  ```bash
  mkdir -p ~/.local/bin
  mv verge ~/.local/bin/
  chmod +x ~/.local/bin/verge
  ```

---

## 4. Uninstalling Verge CLI

Since Verge is compiled and distributed as a single standalone binary with no external configuration directories, daemon files, or registry footprints, uninstalling the application is extremely clean and transparent.

### Automatically Uninstall
Run the following command depending on your active shell:

#### For POSIX Shells (Linux, macOS, Git-Bash, WSL)
```bash
# Removes Verge from standard global and local executable paths
sudo rm -f /usr/local/bin/verge ~/.local/bin/verge /usr/bin/verge.exe /bin/verge.exe
```

#### For Windows PowerShell
```powershell
# Automatically locates and force-deletes verge.exe from your Windows PATH
Remove-Item (Get-Command verge.exe).Path -Force
```

### Path-by-Path Breakdown
If you installed Verge to a custom path or want to verify exactly where it resides:

1. **Verify Binary Location:**
   Query your active shell to find the exact file path of the Verge executable:
   * **In Bash/Zsh:**
     ```bash
     which verge
     # Typical Output: /usr/local/bin/verge
     ```
   * **In PowerShell:**
     ```powershell
     (Get-Command verge.exe).Path
     # Typical Output: C:\Users\username\bin\verge.exe
     ```

2. **Standard Unix/macOS Global Paths:**
   * Global directory: `/usr/local/bin/verge` (requires `sudo` privileges to write/delete).
   * User-specific local directory: `~/.local/bin/verge` (requires no `sudo` privileges).

3. **Standard Windows Git-Bash / MSYS Paths:**
   * Global MSYS bin path: `/usr/bin/verge.exe` or `/bin/verge.exe`

Once you have identified the target path, delete the binary file using `rm` (on Unix/bash) or standard Windows Explorer. Your system is now 100% clean of all Verge remnants.
