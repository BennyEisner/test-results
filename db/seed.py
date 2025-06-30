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

    def create_build(self, test_suite_id: int, build_number: str, ci_provider: str, ci_url: Optional[str], test_case_count: int, duration: Optional[float] = None) -> Optional[int]:
        print(f"    Creating build: {build_number} for test_suite_id: {test_suite_id}")
        if duration is None:
            duration = round(random.uniform(60.0, 600.0), 2)
        return self._execute_returning_id(
            "INSERT INTO builds (test_suite_id, build_number, ci_provider, ci_url, test_case_count, duration) VALUES (%s, %s, %s, %s, %s, %s) RETURNING id",
            (test_suite_id, build_number, ci_provider, ci_url, test_case_count, duration)
        )

    def create_test_case(self, suite_id: int, name: str, classname: str) -> Optional[int]:
        print(f"      Creating test case definition: {name} for suite_id: {suite_id}")
        return self._execute_returning_id(
            "INSERT INTO test_cases(suite_id, name, classname) VALUES (%s, %s, %s) RETURNING id",
            (suite_id, name, classname)
        )

    def create_build_test_case_execution(self, build_id: int, test_case_id: int, status: str, execution_time: float) -> Optional[int]:
        print(f"        Creating execution for build_id: {build_id}, test_case_id: {test_case_id}")
        return self._execute_returning_id(
            "INSERT INTO build_test_case_executions (build_id, test_case_id, status, execution_time) VALUES (%s, %s, %s, %s) RETURNING id",
            (build_id, test_case_id, status, execution_time)
        )

    def create_failure(self, build_test_case_execution_id: int, message: str, type: str, details: str) -> Optional[int]:
        print(f"          Creating failure for execution_id: {build_test_case_execution_id}")
        return self._execute_returning_id(
            "INSERT INTO failures (build_test_case_execution_id, message, type, details) VALUES (%s, %s, %s, %s) RETURNING id",
            (build_test_case_execution_id, message, type, details)
        )
    def seed_data(self, num_projects: int, num_suites_per_project: int, num_builds_per_suite: int, num_test_case_definitions_per_suite: int):
        print("Starting database seeding...")
        project_count = 0
        suite_count = 0
        build_count = 0
        test_case_definitions_count = 0
        build_test_case_executions_count = 0
        failures_count = 0

        for i in range(num_projects):
            project_name = fake.company() + " Project " + str(i+1)
            project_id = self.create_project(project_name)
            if project_id is None: # Project might already exist if names collide by chance
                self.cursor.execute("SELECT id FROM projects WHERE name = %s", (project_name,))
                project_id = self.cursor.fetchone()[0]
            
            if project_id:
                project_count += 1
                for j in range(num_suites_per_project):
                    suite_name = fake.bs().capitalize() + " Test Suite " + str(j+1)
                    suite_time = round(random.uniform(1.0, 100.0), 2) 
                    test_suite_id = self.create_test_suite(project_id, suite_name, suite_time)
                    
                    if test_suite_id:
                        suite_count += 1
                        
                        # Create test case definitions for this suite
                        current_suite_test_case_ids = []
                        for tc_idx in range(num_test_case_definitions_per_suite):
                            tc_name = fake.sentence(nb_words=3).replace('.', '') + f" Def {tc_idx+1}"
                            tc_classname = fake.word().capitalize() + "." + fake.word().capitalize() + "Tests"
                            test_case_definition_id = self.create_test_case(test_suite_id, tc_name, tc_classname)
                            if test_case_definition_id:
                                test_case_definitions_count += 1
                                current_suite_test_case_ids.append(test_case_definition_id) # Tracks number of test cases per suite

                        # Create builds for this suite
                        for k in range(num_builds_per_suite):
                            build_number_str = str(fake.random_int(min=1000, max=9999)) + "-" + fake.lexify(text="??????", letters="ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
                            ci_providers = ["GitHub Actions", "Jenkins", "Travis CI"] 
                            ci_provider_str = random.choice(ci_providers)
                            ci_url_str = fake.url() if random.choice([True, False]) else None
                            
                            build_id = self.create_build(test_suite_id, build_number_str, ci_provider_str, ci_url_str, len(current_suite_test_case_ids))
                            if build_id:
                                build_count += 1
                                
                                # Create executions for each test case definition in this build
                                for test_case_def_id in current_suite_test_case_ids:
                                    exec_time = round(random.uniform(0.01, 15.0), 3)
                                    statuses = ["passed", "failed", "skipped", "error"]
                                    weights = [0.70, 0.15, 0.10, 0.05] 
                                    exec_status = random.choices(statuses, weights=weights, k=1)[0]
                                    
                                    execution_id = self.create_build_test_case_execution(
                                        build_id=build_id,
                                        test_case_id=test_case_def_id,
                                        status=exec_status,
                                        execution_time=exec_time
                                    )
                                    if execution_id:
                                        build_test_case_executions_count += 1
                                        if exec_status == "failed" or exec_status == "error":
                                            failure_message = fake.sentence(nb_words=10)
                                            failure_type = fake.word().capitalize() + "Error"
                                            failure_details = fake.text(max_nb_chars=500)
                                            failure_id = self.create_failure(execution_id, failure_message, failure_type, failure_details)
                                            if failure_id:
                                                failures_count +=1
        self.connection.commit()
        print(f"\nSeeding complete!")
        print(f"  Total projects created: {project_count}")
        print(f"  Total test suites created: {suite_count}")
        print(f"  Total test case definitions created: {test_case_definitions_count}")
        print(f"  Total builds created: {build_count}")
        print(f"  Total build test case executions created: {build_test_case_executions_count}")
        print(f"  Total failures created: {failures_count}")

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
        num_projects_to_create = random.randint(5, 10)
        num_suites_per_project_to_create = random.randint(3, 7) 
        num_builds_per_suite_to_create = random.randint(2, 4) # Reduced
        num_test_case_definitions_per_suite_to_create = random.randint(15, 30)

        print(f"Attempting to create:")
        print(f"  - {num_projects_to_create} projects")
        print(f"  - {num_suites_per_project_to_create} test suites per project")
        print(f"  - {num_test_case_definitions_per_suite_to_create} test case definitions per suite")
        print(f"  - {num_builds_per_suite_to_create} builds per test suite (each executing all test cases defined for the suite)")
        
        seeder.seed_data(
            num_projects=num_projects_to_create,
            num_suites_per_project=num_suites_per_project_to_create,
            num_builds_per_suite=num_builds_per_suite_to_create,
            num_test_case_definitions_per_suite=num_test_case_definitions_per_suite_to_create
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
