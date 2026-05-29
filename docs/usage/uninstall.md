# Uninstalling Verge CLI

Since Verge is compiled and distributed as a single standalone binary with no external configuration directories, daemon files, or registry footprints, uninstalling the application is extremely clean and transparent.

---

Run the following command depending on your active shell:

### For POSIX Shells (Linux, macOS, Git-Bash, WSL)
```bash
# Removes Verge from standard global and local executable paths
sudo rm -f /usr/local/bin/verge ~/.local/bin/verge /usr/bin/verge.exe /bin/verge.exe
```

### For Windows PowerShell
```powershell
# Automatically locates and force-deletes verge.exe from your Windows PATH
Remove-Item (Get-Command verge.exe).Path -Force
```

---

## 2. Path-by-Path Breakdown

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
