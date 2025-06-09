# scripts/cli.py
import os
import shutil
import subprocess
from pathlib import Path

import click
import psycopg
from dotenv import load_dotenv
from psycopg import errors as psycopg_errors

# Load environment variables from .env file if it exists
load_dotenv()

# Default DB connection info (can be overridden by env vars)
DB_CONFIG = {
    "dbname": os.getenv("DB_NAME"),
    "user": os.getenv("DB_USER"),
    "password": os.getenv("DB_PASSWORD"),
    "host": os.getenv("DB_HOST", "localhost"),
    "port": int(os.getenv("DB_PORT", "5432")),
}

# Configuration for connecting to the postgres database
PG_CONFIG = {**DB_CONFIG}
PG_CONFIG["dbname"] = "postgres"


def validate_config():
    """Validate that all required database configuration is present."""
    missing = []
    for key in ["dbname", "user", "password"]:
        if not DB_CONFIG.get(key):
            missing.append(key)

    if missing:
        raise ValueError(
            f"Missing required database configuration: {', '.join(missing)}"
        )


def load_sql_and_execute(path: Path, conn: psycopg.Connection):
    """
    Load SQL from a file and execute it against the database.

    Args:
        path: Path to the SQL file
        conn: Database connection

    Raises:
        FileNotFoundError: If the SQL file doesn't exist
        ValueError: If the SQL file is too large
    """
    if not path.exists():
        raise FileNotFoundError(f"SQL file not found: {path}")

    # Basic safety check - don't execute extremely large files
    file_size = path.stat().st_size
    if file_size > 10 * 1024 * 1024:  # 10MB limit
        raise ValueError(f"SQL file too large ({file_size} bytes). Limit is 10MB.")

    with path.open("r") as f:
        sql = f.read()
        try:
            with conn.cursor() as cur:
                cur.execute(sql)
                conn.commit()
            click.echo(f"✅ Executed SQL from {path}")
        except Exception as e:
            conn.rollback()
            click.echo(f"❌ SQL execution failed: {e}")
            raise


@click.group()
def cli():
    """Database utility CLI for test-results."""
    try:
        validate_config()
    except ValueError as e:
        click.echo(f"Configuration error: {e}")
        exit(1)


@cli.command("create-db")
@click.option("--dbname", help="Override the database name from environment")
def create_db(dbname):
    """
    Create the database if it doesn't exist.

    Args:
        dbname: Override the database name from environment variables
    """
    if dbname:
        db_name = dbname
    else:
        db_name = DB_CONFIG["dbname"]

    try:
        with psycopg.connect(**PG_CONFIG) as connection:
            connection.autocommit = True
            with connection.cursor() as cursor:
                cursor.execute(
                    "SELECT 1 FROM pg_database WHERE datname = %s", (db_name,)
                )
                if cursor.fetchone():
                    click.echo(f"Database '{db_name}' exists already")
                    return
                cursor.execute(f'CREATE DATABASE "{db_name}"')
                click.echo(f"Database '{db_name}' successfully created")
    except psycopg_errors.OperationalError as e:
        click.echo(f"Database creation failed: '{e}'")
    except Exception as e:
        click.echo(f"Unexpected error during database creation: {e}")


@cli.command("init-schema")
@click.argument(
    "schema_file", type=click.Path(exists=True, dir_okay=False, path_type=Path)
)
def init_schema(schema_file):
    """
    Initialize the database schema using the given SCHEMA_FILE.

    Args:
        schema_file: Path to the SQL file containing the schema definition
    """
    try:
        with psycopg.connect(**DB_CONFIG) as connection:
            load_sql_and_execute(schema_file, connection)
            click.echo("Database schema initialized")
    except psycopg_errors.OperationalError as e:
        click.echo(f"Schema initialization failed: {e}")
    except Exception as e:
        click.echo(f"Unexpected error during schema initialization: {e}")


@cli.command("seed-db")
@click.argument(
    "sql_file", type=click.Path(exists=True, dir_okay=False, path_type=Path)
)
def seed_db(sql_file):
    """
    Seed the database with example data using SQL_FILE.

    Args:
        sql_file: Path to the SQL file containing the seed data
    """
    try:
        with psycopg.connect(**DB_CONFIG) as conn:
            load_sql_and_execute(sql_file, conn)
    except Exception as e:
        click.echo(f"Database seeding failed: {e}")


@cli.command("dump-db")
@click.argument("output_file", type=click.Path(dir_okay=False, path_type=Path))
@click.option("--schema-only", is_flag=True, help="Dump only the schema.")
def dump_db(output_file, schema_only):
    """
    Dump the database to OUTPUT_FILE using pg_dump.

    Args:
        output_file: Path where the database dump will be saved
        schema_only: If True, dump only the schema without data
    """
    # Check if pg_dump is installed
    if not shutil.which("pg_dump"):
        click.echo(
            "❌ pg_dump not found. Make sure PostgreSQL client tools are installed."
        )
        return

    args = [
        "pg_dump",
        "-U",
        DB_CONFIG["user"],
        "-h",
        DB_CONFIG["host"],
        "-p",
        str(DB_CONFIG["port"]),
        "-d",
        DB_CONFIG["dbname"],
        "-f",
        str(output_file),
    ]
    if schema_only:
        args.append("--schema-only")

    env = {**os.environ, "PGPASSWORD": DB_CONFIG["password"]}

    try:
        subprocess.run(args, check=True, env=env)
        click.echo(f"✅ Dumped database to {output_file}")
    except subprocess.CalledProcessError as e:
        click.echo(f"❌ pg_dump failed: {e}")


@cli.command("delete-db")
@click.option("--force", is_flag=True, help="Force drop without confirmation")
@click.option("--dbname", help="Override the database name from environment")
def delete_db(force, dbname):
    """
    Drop (delete) the database completely.

    This command deletes the entire database specified in your environment
    variables or via the --dbname option. Use with extreme caution as this
    will permanently delete all data.

    Args:
        force: If True, skip confirmation prompt
        dbname: Override the database name from environment variables
    """
    db_name = dbname if dbname else DB_CONFIG["dbname"]

    # Safety confirmation unless --force is used
    if not force:
        confirmation = click.prompt(
            f"Are you sure you want to drop database '{db_name}'?"
            f"Type '{db_name}' to confirm",
            type=str,
        )
        if confirmation != db_name:
            click.echo("❌ Database drop cancelled.")
            return

    # Connect to postgres database (can't drop a database while connected to it)
    try:
        with psycopg.connect(**PG_CONFIG) as connection:
            # Set autocommit mode (required for dropping databases)
            connection.autocommit = True

            with connection.cursor() as cursor:
                # First terminate any existing connections to the database
                cursor.execute(
                    """
                    SELECT pg_terminate_backend(pg_stat_activity.pid)
                    FROM pg_stat_activity
                    WHERE pg_stat_activity.datname = %s
                    AND pid <> pg_backend_pid()
                    """,
                    (db_name,),
                )

                # Now drop the database
                cursor.execute(f'DROP DATABASE IF EXISTS "{db_name}"')
                click.echo(f"✅ Database '{db_name}' dropped successfully!")
    except psycopg_errors.OperationalError as e:
        click.echo(f"❌ Database drop failed: {e}")
    except Exception as e:
        click.echo(f"❌ Unexpected error during database drop: {e}")


if __name__ == "__main__":
    cli()
