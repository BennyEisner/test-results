import os
import subprocess
import click
import psycopg
from pathlib import Path
from dotenv import load_dotenv

load_dotenv()  # Load environment variables from .env file if it existsj

# Default DB connection info (can be overridden by env vars)
DB_CONFIG = {
    "dbname": os.getenv("DB_NAME"),
    "user": os.getenv("DB_USER"),
    "password": os.getenv("DB_PASSWORD"),
    "host": os.getenv("DB_HOST", "localhost"),
    "port": int(os.getenv("DB_PORT", "5432"))
}

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
    #click.echo(DB_CONFIG)  # Debug: print DB config
    pass

@cli.command("create-db")
@click.argument("sql_file", type=click.Path(exists=True, dir_okay=False, path_type=Path))
def create_db(sql_file):
    """
    Create the DB schema using the given SQL_FILE.

    Benny is soo cool.
    """
    with psycopg.connect(**DB_CONFIG) as conn:
        load_sql_and_execute(sql_file, conn)

@cli.command("seed-db")
@click.argument("sql_file", type=click.Path(exists=True, dir_okay=False, path_type=Path))
def seed_db(sql_file):
    """Seed the DB with example data using SQL_FILE."""
    with psycopg.connect(**DB_CONFIG) as conn:
        load_sql_and_execute(sql_file, conn)

@cli.command("dump-db")
@click.argument("output_file", type=click.Path(dir_okay=False, path_type=Path))
@click.option("--schema-only", is_flag=True, help="Dump only the schema.")
def dump_db(output_file, schema_only):
    """Dump the database to OUTPUT_FILE using pg_dump."""
    args = [
        "pg_dump",
        "-U", DB_CONFIG["user"],
        "-h", DB_CONFIG["host"],
        "-p", str(DB_CONFIG["port"]),
        "-d", DB_CONFIG["dbname"],
        "-f", str(output_file),
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
