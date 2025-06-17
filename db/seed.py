import psycopg2
import os
import random
from faker import Faker
from typing import Optional, Tuple, List

# Initialize Faker
fake = Faker()

class DatabaseSeeder:
    def __init__(self, connection_url: str):
        self.connection = psycopg2.connect(connection_url)
        self.cursor = self.connection.cursor()
        print("Successfully connected to the database.")

    def _execute_query(self, query: str, params: Optional[tuple] = None) -> Optional[List[Tuple]]:
        self.cursor.execute(query, params)
        try:
            return self.cursor.fetchall()
        except psycopg2.ProgrammingError: # No results to fetch
            return None

    def _execute_returning_id(self, query: str, params: Optional[tuple] = None) -> Optional[int]:
        self.cursor.execute(query, params)
        result = self.cursor.fetchone()
        return result[0] if result else None

    def create_project(self, name: str) -> Optional[int]:
        print(f"Creating project: {name}")
        return self._execute_returning_id(
            "INSERT INTO projects (name) VALUES (%s) ON CONFLICT (name) DO NOTHING RETURNING id",
            (name,)
        )

    def create_test_suite(self, project_id: int, name: str, time: float) -> Optional[int]:
        print(f"  Creating test suite: {name} for project_id: {project_id}")
        return self._execute_returning_id(
            "INSERT INTO test_suites (project_id, name, time) VALUES (%s, %s, %s) RETURNING id",
            (project_id, name, time)
        )

    def create_build(self, test_suite_id: int, build_number: str, ci_provider: str, ci_url: Optional[str]) -> Optional[int]:
        print(f"    Creating build: {build_number} for test_suite_id: {test_suite_id}")
        return self._execute_returning_id(
            "INSERT INTO builds (test_suite_id, build_number, ci_provider, ci_url) VALUES (%s, %s, %s, %s) RETURNING id",
            (test_suite_id, build_number, ci_provider, ci_url)
        )

    def seed_data(self, num_projects: int, num_suites_per_project: int, num_builds_per_suite: int):
        print("Starting database seeding...")
        project_count = 0
        suite_count = 0
        build_count = 0

        for i in range(num_projects):
            project_name = fake.company() + " Project " + str(i+1)
            project_id = self.create_project(project_name)
            if project_id is None: # Project might already exist if names collide by chance
                self.cursor.execute("SELECT id FROM projects WHERE name = %s", (project_name,))
                project_id = self.cursor.fetchone()[0]
            
            if project_id:
                project_count +=1
                for j in range(num_suites_per_project):
                    suite_name = fake.bs().capitalize() + " Test Suite " + str(j+1)
                    suite_time = round(random.uniform(1.0, 100.0), 2)
                    test_suite_id = self.create_test_suite(project_id, suite_name, suite_time)
                    
                    if test_suite_id:
                        suite_count += 1
                        for k in range(num_builds_per_suite):
                            build_number_str = str(fake.random_int(min=1000, max=9999)) + "-" + fake.lexify(text="??????", letters="ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
                            ci_providers = ["GitHub Actions", "Jenkins", "Travis CI"]
                            ci_provider_str = random.choice(ci_providers)
                            ci_url_str = fake.url() if random.choice([True, False]) else None
                            
                            build_id = self.create_build(test_suite_id, build_number_str, ci_provider_str, ci_url_str)
                            if build_id:
                                build_count += 1
        
        self.connection.commit()
        print(f"\nSeeding complete!")
        print(f"  Total projects created: {project_count}")
        print(f"  Total test suites created: {suite_count}")
        print(f"  Total builds created: {build_count}")


    def close_connection(self):
        if self.connection:
            self.cursor.close()
            self.connection.close()
            print("Database connection closed.")

def main():
    # Database connection details from environment variables
    db_host = os.getenv("DB_HOST", "localhost")
    db_port = os.getenv("DB_PORT", "5433")
    db_user = os.getenv("DB_USER", "postgres")
    db_password = os.getenv("DB_PASSWORD", "postgrespassword")
    db_name = os.getenv("DB_NAME", "test_results")

    connection_url = f"postgres://{db_user}:{db_password}@{db_host}:{db_port}/{db_name}"

    seeder = None
    try:
        seeder = DatabaseSeeder(connection_url)

        # Define how much data to generate
        num_projects_to_create = random.randint(10, 20)
        num_suites_per_project_to_create = random.randint(10, 20)
        num_builds_per_suite_to_create = random.randint(2, 5)

        print(f"Attempting to create:")
        print(f"  - {num_projects_to_create} projects")
        print(f"  - {num_suites_per_project_to_create} test suites per project")
        print(f"  - {num_builds_per_suite_to_create} builds per test suite")
        
        seeder.seed_data(
            num_projects=num_projects_to_create,
            num_suites_per_project=num_suites_per_project_to_create,
            num_builds_per_suite=num_builds_per_suite_to_create
        )

    except psycopg2.OperationalError as e:
        print(f"Database connection error: {e}")
        print("Please ensure the database is running and accessible, and environment variables are set correctly:")
        print(f"  DB_HOST={db_host}, DB_PORT={db_port}, DB_USER={db_user}, DB_NAME={db_name}")
        print(f"  (DB_PASSWORD is read from env but not printed for security)")
    except Exception as error:
        print(f"An unexpected error occurred: {error}")
    finally:
        if seeder:
            seeder.close_connection()

if __name__ == "__main__":
    main()
