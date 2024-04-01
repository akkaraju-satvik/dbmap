# dbmap

`dbmap` is a database migration management tool for PostgreSQL.

>Note: `dbmap` is not a SQL generator. It is a database migration management tool.
> `dbmap` will help you manage them by keeping track of which scripts have been run, and which have not.

## Installation

### From the releases page

Download the latest release from the [releases page](https://github.com/akkaraju-satvik/dbmap/releases).

### From source

- Clone the repository

```bash
git clone https://github.com/akkaraju-satvik/dbmap
```

- Install the dependencies

```bash
go mod download
```

- Install the binary

```bash
go install
```

## Usage

### Initialize a new project

```bash
dbmap init -c <connection-string>
```

### Create a new migration

```bash
dbmap create-migration
```

### Apply migrations

```bash
dbmap apply-migrations
```
