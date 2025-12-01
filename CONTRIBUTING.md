# 📄 **CONTRIBUTING.md**

````markdown
# Contributing to This Fork of `dexidp/dex`

This repository is a **fork** of the public Dex project:

- **Upstream (public project):** `dexidp/dex`
- **Fork (your working repo):** `fperron-kt88/dex`

Because this repo exists inside a larger project (e.g. `fx-test-generic-oidc`),
and because development happens on a fork—not directly on the public repo—it is
important to follow a clear and safe workflow.

This document explains:

- How remotes are configured
- How to create branches
- Where to push and where NOT to push
- How to pull changes from your fork vs from the real Dex project
- How to keep your fork in sync with upstream

---

## 🔧 Remote Configuration

Your local clone uses **two remotes**:

| Remote name         | URL                                      | Purpose                                |
| ------------------- | ---------------------------------------- | -------------------------------------- |
| `origin`            | `https://github.com/fperron-kt88/dex.git` | Your fork (read/write)                 |
| `upstream-dexidp`   | `https://github.com/dexidp/dex.git`       | Public Dex project (read-only, fetch)  |

You **always push** to `origin`.  
You **never push** to `upstream-dexidp`.

Pre-push hooks are installed to prevent mistakes.

---

## 🌿 Branching Model

### **Main development rules**

- `main` is your fork’s default branch.
- **NEVER develop directly on `main`.**
- Create feature branches off `main`:

```bash
git checkout main
git pull
git checkout -b feature/some-change
````

### **Feature branch naming**

Use:

```
feature/<description>
fix/<description>
chore/<description>
```

Example:

```
feature/fx-idp-integration
```

---

## 🔄 Pulling Changes

There are **two types of pulls**:

### 1. Pulling from **your fork** (`origin`)

Use this when:

* You pushed changes from another machine
* GitHub automation updated your fork
* You are collaborating with yourself across systems

Commands:

```bash
git pull        # pulls from origin for the current branch
git pull origin main
```

---

### 2. Pulling from **upstream Dex** (`upstream-dexidp`)

Use this to get new features, bug fixes, or releases from the real Dex project.

Fetch:

```bash
git fetch upstream-dexidp
```

Merge:

```bash
git merge upstream-dexidp/master
```

Rebase (cleaner):

```bash
git rebase upstream-dexidp/master
```

This will update your fork with the latest public Dex changes.

---

## 🚀 Pushing Code

Push with your safe alias:

```bash
git pushup
```

This automatically:

* Pushes the current branch to **origin**
* Sets the correct tracking branch (`origin/<branch>`)
* Avoids accidental pushes to upstream

---

## 🔒 Safety Tools

The repo includes:

### **1. Pre-push hook**

Blocks accidental pushes to the public Dex repo:

```
.git/hooks/pre-push
```

### **2. Pre-pull hook**

Warns when pulling from the wrong place on important branches.

### **3. Safety Status Command**

```
git safety-status
```

Example:

```
origin:            https://github.com/fperron-kt88/dex.git
upstream-dexidp:   https://github.com/dexidp/dex.git
current branch:    feature/fx-idp-integration
tracking branch:   origin/feature/fx-idp-integration
STATUS: ✔ SAFE - origin is writable
```

---

## 📦 Submodule Integration (optional)

If this repo is used as a submodule in a parent project:

* Always point the submodule to a **feature branch** or **specific tag**
* Never point submodules to `main` unless intentionally syncing upstream
* Use:

```bash
git submodule update --remote
```

Only when you want to import new upstream changes.

---

## 📝 Summary

* **Push** to `origin`, never upstream.
* **Pull** from `origin` for your own work,
  **pull from upstream-dexidp** to sync with Dex.
* **Always** work on feature branches.
* Use the **git safety tools** to avoid dangerous operations.

---

## 👍 Need Help?

If you encounter:

* Conflicts syncing upstream
* Branch divergence
* Submodule alignment issues
* Remote confusion

Use:

```
git safety-status
```

Or ask in the development notes.

```
