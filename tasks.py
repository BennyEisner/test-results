from invoke.tasks import task

PYTHON = "python3"
FLIT = "flit"
SRC_DIR = "scripts"
TEST_DIR = "tests"

@task(name="install", aliases=["i"])
def install(c):
    """Install project using Flit (including optional dev dependencies)."""
    c.run(f"{FLIT} install --symlink --deps=all")

@task(name="lint", aliases=["l"])
def lint(c):
    """Run static analysis using flake8."""
    c.run(f"flake8 {SRC_DIR} {TEST_DIR}")

@task(name="typecheck", aliases=["t"])
def typecheck(c):
    """Run mypy for type checking."""
    c.run(f"mypy {SRC_DIR}")

@task(name="test", aliases=["T"])
def test(c):
    """Run tests using pytest."""
    c.run("pytest")

@task(name="coverage", aliases=["c"])
def coverage(c):
    """Run test coverage analysis."""
    c.run("coverage run -m pytest")
    c.run("coverage report")

@task(pre=[lint, typecheck, test], name="check", aliases=["C"])
def check(c):
    """Run all code quality checks and tests."""
    print("✅ All checks passed!")

@task(name="clean", aliases=["x"])
def clean(c):
    """Clean up temporary files and __pycache__."""
    c.run("find . -type d -name '__pycache__' -exec rm -r {} +")
    c.run("rm -rf .pytest_cache .mypy_cache .coverage coverage.xml")

@task(name="fmt", aliases=["f"])
def fmt(c):
    """Format code using black and isort."""
    c.run(f"isort {SRC_DIR}")
    c.run(f"black {SRC_DIR}")
