# User Dashboard Order

## Problems:

Order of implementation for Users and Dashboard Configuration
Method of implementing both systems

## Context:

Need to implement a system where users can customize their dashboard. Unsure what order to implement required components for this.
Also need to decide how both systems should be done (ie should users and user config be a seperate database from testresults)

## Options Considered:

### Prolem 1 (Order of implementation)

#### Option 1: Users First, Then Dashboard Configuration

Create user authentication system first
Once users can be authenticated, add dashboard configuration capabilities
**Pros:**
Clear user ownership of configurations from the start
Authentication security established early
Natural progression (users → their preferences)
**Cons:**
Delays dashboard functionality
May require retrofitting dashboard features to user system
Can't test dashboard features with real configurations early

#### Option 2: Dashboard Configuration First, Then Users

Create dashboard configuration system independent of users
Later integrate user authentication and link users to configurations
**Pros:**
Dashboard functionality available sooner
Can test and refine dashboard features before user integration
More flexible initial development
**Cons:**
May require rework to link configurations to users later
Temporary solution needed for storing/identifying configurations
Security considerations handled later in process

### Prolem 2 (Method of implementation)

#### Option 1: Unifieid Database

Stores users, configuartion, and test results all in a single postgres database

**Pros:**
Simpler architecture
Easier to query

**Cons:**
Could run into scaling issues
Changing database schema will likley require a substantial refactoring of current code

#### Option 2: Seperate database approaches

Keep user and user configuration settings seperated from test results database

**Pros:**
Components will scale independently
Errors more isolated
can optimize/structure each database for its specific needs
**Cons:**
More comples architecture
Requires database communication

## Decision outcomes

#### Prolem 1 (Order of implementation)

#### Prolem 2 (Method of implementation)
