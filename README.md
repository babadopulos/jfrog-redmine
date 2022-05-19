# JFrog-redmine

## About this plugin
This plugin uses x-ray to scan a maven project for vulnerabilities.
It also maintains a redmine issue tracking up to date, adding and closing known issues automatically.

## Usage

### Build
```
$ go build
```

### Install as JFrog CLi plugin
```
$ mkdir -p ~/.jfrog/plugins/jfrog-redmine/bin/
$ cp jfrog-redmine ~/.jfrog/plugins/jfrog-redmine/bin/
```

### Commands
```
jf jfrog-redmine audit --help
```

* audit
    - Options:
      - --source     [Mandatory] Source code directory.
      - --project    [Mandatory] Redmine Project identifier
      - --dryrun     [Default: false] Show what would have been done
    - Example:
    ```
  $ jf jfrog-redmine audit --source=/path/to/maven/project --project=project-1 --dryrun

  ```

### Environment variables
* REDMINE_API_ENDPOINT - Base endpoint to redmine [example: http:/127.0.0.1:8080]**
* REDMINE_API_KEY - Redmine API key used for user authentication**

## Additional info

[docker-compose](docker-compose.yml) file used to deploy a redmine issue tracking platform locally.

[mysql dump](redmine.sql) file used to restore mysql to initial redmine state.


# How to deploy a working environment to test this plugin

Start docker containers
```bash
docker-compose up
```

Export Mysql container ID
```bash
REDMINE_DB_CONTAINER_ID=$(docker ps -aqf "name=jfrog-redmine_db_1")
```

Restore Mysql to Redmine initial state
```bash
cat redmine.sql | docker exec -i $REDMINE_DB_CONTAINER_ID /usr/bin/mysql -u root --password=mysqlpwd redmine
```

point your browser to: http://127.0.0.1:8080

Redmine UI credentials
```text
user: admin
pass: redmineadmin
```

Export env variables used by this plugin
```bash
export REDMINE_API_ENDPOINT=http://127.0.0.1:8080
export REDMINE_API_KEY=86e264af0a51949f8d0364058557bfe14b046e54
```

Build the plugin
```bash
go build
```

Install binary as jfrog plugin
```bash
mkdir -p ~/.jfrog/plugins/jfrog-redmine/bin/
cp jfrog-redmine ~/.jfrog/plugins/jfrog-redmine/bin/
```

Using jfrog-plugin
```bash
jf jfrog-redmine audit --help
```

<pre>
Name:
   jfrog-redmine audit - Audit mvn project

Usage:
    jfrog-redmine audit [command options]

Options:
  --source     [Mandatory] Source code directory.
  --project    [Mandatory] Redmine Project identifier
  --dryrun     [Default: false] Show what would have been done

Environment Variables:
  REDMINE_API_ENDPOINT
    Redmine API endpoint

  REDMINE_API_KEY
    Redmine API key
</pre>

Audit example
```
jf jfrog-redmine audit --source=./sample-maven-project/java-vulnerable-sample --project=my-mvn-project --dryrun
```


## Release Notes
The release notes are available [here](RELEASE.md).
