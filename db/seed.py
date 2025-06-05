import psycopg2
from typing import Optional, Tuple

class ProjectRepository:
    """Repository for managing projects in the database."""
    
    def __init__(self, connection):
        """Initialize with database connection."""

        self.connection = connection
        self.cursor = connection.cursor()

    def exists_by_name(self, name: str) -> bool:
        """Check if a project with the given name exists."""

        self.cursor.execute("SELECT EXISTS(SELECT 1 FROM projects WHERE name = %s)", (name,))
        return self.cursor.fetchone()[0]

    def create(self, name: str) -> Optional[Tuple[int, str]]:
        """Create a new project if it doesn't already exist."""

        if not self.exists_by_name(name):
            self.cursor.execute("INSERT INTO projects (name) VALUES (%s) RETURNING id, name",(name,))
            return self.cursor.fetchone()
        return None


    def delte_by_id(self, project_id: int) -> bool:
        self.cursor.execute("DELETE FROM projects WHERE id = %s RETURNING id", (project_id,))
        deleted = self.cursor.fetchone()
        return deleted is not None
    
    def save_changes(self):
        """Commit changes to the database."""
        self.connection.commit()


def main():
    connection = None
    
    try:
        url = "postgres://admin:secret@localhost:5432/testdb"
        connection = psycopg2.connect(url)
        print("Connected to database")
        
        # Initialize repository
        repo = ProjectRepository(connection)
        
        # Create several projects
        projects = ["abcd", "efgh", "ijkl", "Duplicate checker example project addition"]

        for name in projects:
            result = repo.create(name)
            if result:
                print(f"Created project: {result}")
            else:
                print(f"Project '{name}' already exists")
        
        # Save all changes
        repo.save_changes()
        
        # Display PostgreSQL version
        repo.cursor.execute("SELECT version()")
        version = repo.cursor.fetchone()
        print(f"PostgreSQL version: {version[0]}")
        
        # Example of searching for a specific project
        repo.cursor.execute("SELECT * FROM projects WHERE name = %s", ("abcd",))
        result = repo.cursor.fetchone()
        print(f"Search for 'abcd': {result if result else 'Not found'}")
        
        # Delte a project by an ID
        delete_id = 2 # Hard coded but will add flexability later 
        print(f"Trying to deleting project with ID={delete_id}")

        success = repo.delte_by_id(delete_id)
        if success:
            print(f"Successfully deleted project with ID={delete_id}")
        else:
            print(f"No project found with ID={delete_id}")
        repo.save_changes()

    except Exception as error:
        print(f"Error: {error}")
        
    finally:
        # Clean up connection
        if connection:
            connection.close()
            print("Connection closed")


if __name__ == "__main__":
    main()