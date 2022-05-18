# frog-redmine

## About this plugin
This plugin uses x-ray to scan a maven project for vulnerabilities.
It also maintains a redmine issue tracking up to date, adding and closing known issues automatically.

## Usage

### Build
```
$ go build
```

### Install as jFrog CLi plugin
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
[docker-compose](RELEASE.md) file used to deploy a redmine issue tracking platform locally.

## Release Notes
The release notes are available [here](RELEASE.md).
