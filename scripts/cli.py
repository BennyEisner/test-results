import os
import subprocess
from pathlib import Path

import click
import psycopg
from dotenv import load_dotenv
from psycopg.errors import OperationalError

load_dotenv()  # Load environment variables from .env file if it exists

# Default DB connection info (can be overridden by env vars)
DB_CONFIG = {
    "dbname": os.getenv("DB_NAME"),
    "user": os.getenv("DB_USER"),
    "password": os.getenv("DB_PASSWORD"),
    "host": os.getenv("DB_HOST", "localhost"),
    "port": int(os.getenv("DB_PORT", "5432")),
}

PG_CONFIG = {**DB_CONFIG}
PG_CONFIG["dbname"] = "postgres"


def load_sql_and_execute(path: Path, conn: psycopg.Connection):
    if not path.exists():
        raise FileNotFoundError(f"SQL file not found: {path}")
    with path.open("r") as f:
        sql = f.read()
        with conn.cursor() as cur:
            cur.execute(sql)
            conn.commit()
    click.echo(f"✅ Executed SQL from {path}")


@click.group()
def cli():
    """Database utility CLI for test-results."""
    # click.echo(DB_CONFIG)  # Debug: print DB config
    pass


@cli.command("create-db")
def create_db():
    """
    Create a new database.
    """
    dbname = DB_CONFIG["dbname"]
    try:
        with psycopg.connect(**PG_CONFIG) as connection:
            connection.autocommit = True
            with connection.cursor() as cursor:
                cursor.execute(
                    "SELECT 1 FROM pg_database WHERE datname = %s", (dbname,)
                )
                if cursor.fetchone():
                    click.echo(f"Database '{dbname}' exists already")
                    return
                cursor.execute(f'CREATE DATABASE "{dbname}"')
                click.echo(f"Database '{dbname}' successfully created")
    except OperationalError as e:
        click.echo(f"Database creation failed: '{e}'")


@cli.command("init-schema")
@click.argument(
    "schema_file", type=click.Path(exists=True, dir_okay=False, path_type=Path)
)
def init_schema(schema_file):
    """
    Initialize the database schema using the given SCHEMA_FILE.
    """
    try:
        with psycopg.connect(**DB_CONFIG) as connection:
            load_sql_and_execute(schema_file, connection)
            click.echo("Database schema initialized")
    except OperationalError as e:
        click.echo(f"Schema init failed: {e}")


@cli.command("seed-db")
@click.argument(
    "sql_file", type=click.Path(exists=True, dir_okay=False, path_type=Path)
)
def seed_db(sql_file):
    """Seed the DB with example data using SQL_FILE."""
    try:
        with psycopg.connect(**DB_CONFIG) as conn:
            load_sql_and_execute(sql_file, conn)
            click.echo("Database seeded successfully")
    except OperationalError as e:
        click.echo(f"Database seeding failed: {e}")


@cli.command("dump-db")
@click.argument("output_file", type=click.Path(dir_okay=False, path_type=Path))
@click.option("--schema-only", is_flag=True, help="Dump only the schema.")
def dump_db(output_file, schema_only):
    """Dump the database to OUTPUT_FILE using pg_dump."""
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

    env = {**subprocess.os.environ, "PGPASSWORD": DB_CONFIG["password"]}

    try:
        subprocess.run(args, check=True, env=env)
        click.echo(f"✅ Dumped database to {output_file}")
    except subprocess.CalledProcessError as e:
        click.echo(f"❌ pg_dump failed: {e}")


if __name__ == "__main__":
    cli()
