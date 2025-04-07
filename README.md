# Prepare Requirement

- golang
- MacOs OR Linux OR Windows

# Installation

```
git clone git@github.com:BounkhongDev/go-generator.git
```

```
cd go-module-generator
```

### MacOs OR Linux
Run this command in your terminal
```
$ ./install.sh
```

### Windows

To add a directory to your system's PATH:

1. Copy this floder to Local Disk (C:)
2. Right-click on 'This PC' or 'Computer' on your desktop or in File Explorer.
3. Click on 'Properties'.
4. Click on 'Advanced system settings'.
5. Click on the 'Environment Variables' button.
6. In the 'System variables' section, find the Path variable and select it.
7. Click 'Edit', then 'New', and add the directory path (C:\go-module-generator).
8. Click 'OK' to close all windows.

# Usege

### use for initial project:
```
mkdir <yourProjectName>
```
```
go-gen init
```
```
Enter Project Name: <yourProjectName>
```
```
go mod tidy
```

### use for generate module
> [!IMPORTANT]
> In your directory structure must has been "/src/"

```
go-gen <yourServiceName>
```
