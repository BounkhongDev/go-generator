# Go Generator

A simple and opinionated project generator for Go (Golang), designed to help you quickly scaffold projects and modules using a clean structure with repositories.

---

## ✅ Features

- Generate a new Go project structure
- Easily create new modules with built-in templates
- Supports `macOS`, `Linux`, and `Windows`
- Automatically adds the generator to your system path

---

## 📦 Requirements

- [Go (Golang)](https://golang.org/dl/) installed
- One of the following operating systems:
  - macOS
  - Linux
  - Windows

---

## 🚀 Installation

### Clone the Repository

```bash
git clone git@github.com:BounkhongDev/go-generator.git
cd go-generator
```

### For macOS or Linux

Run the installation script:

```bash
./install.sh
```

This will move the binary to `/usr/local/bin` and make `go-gen-r` available globally.

### For Windows

To add the directory to your system’s `PATH` manually:

1. Copy the `go-generator` folder to your Local Disk (`C:`)
2. Right-click on `This PC` or `Computer` on your desktop or in File Explorer
3. Select `Properties`
4. Click on `Advanced system settings`
5. Click the `Environment Variables` button
6. In the **System variables** section, find the `Path` variable and select it
7. Click `Edit`, then `New`, and add the path:  
   `C:\go-generator`
8. Click `OK` to save and close all windows

Now, you can use `go-gen-r` from any terminal window.

---

## 🛠️ Usage

### 🔧 Initialize a New Project

```bash
mkdir <yourProjectName>
cd <yourProjectName>
go-gen-r init
```

You’ll be prompted to enter your project name:

```
Enter Project Name: <yourProjectName>
```

Then, tidy up dependencies:

```bash
go mod tidy
```

---

### 🧱 Generate a New Module

> **IMPORTANT**  
> Your project must follow the directory structure and include a `/src/` folder.

Inside your project directory:

```bash
go-gen-r <yourServiceName>
```

This will generate a new module inside `/internal/<yourServiceName>` with boilerplate code for handler, service, repositories, and routes.

---
