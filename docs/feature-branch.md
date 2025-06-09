# 🧑🏾‍💼 Git Workflow: Creating and Pushing a Feature Branch

This guide walks you through the basic steps to create a new feature branch, add your work to it, and push it to GitHub.

---

## ✅ 1. Create a New Branch

Use a descriptive name for the feature you're working on (e.g., `feature/login-page`):

```bash
git checkout -b feature/your-feature-name
```

This creates and switches to the new branch.

---

## 📂 2. Add New or Modified Files

After making your changes or creating new files, check the status of your working directory:

```bash
git status
```

Add specific files:

```bash
git add path/to/file1 path/to/file2
```

Or add all changes:

```bash
git add .
```

---

## 📝 3. Commit Your Changes

Write a clear commit message describing what you did:

```bash
git commit -m "Add login form and validation logic"
```

---

## 🚀 4. Push the Feature Branch to GitHub

```bash
git push -u origin feature/your-feature-name
```

The `-u` flag sets the upstream so future pushes can be done with just `git push`.

---

## 📀 5. Update Your Branch (as needed)

To keep your feature branch up to date with the `main` branch:

```bash
git checkout main
git pull origin main
git checkout feature/your-feature-name
git merge main
```

---

## 📌 Tips

* Use kebab-case or slash-separated names for branches (e.g., `feature/new-api-endpoint`).
* Commit frequently with meaningful messages.
* Push regularly to avoid losing your work.

