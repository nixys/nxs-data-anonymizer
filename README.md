![photo_2023-07-28_15-28-52](https://github.com/nixys/nxs-data-anonymizer/assets/27485608/165a90a0-929f-460b-8dbd-2903c0d91f36)

# nxs-data-anonymizer

## Introduction

nxs-data-anonymizer is a tool for anonymizing **PostgreSQL** and **MySQL/MariaDB/Percona** databases' dump.

### Features

- Supported databases and versions:
  - PostgreSQL (9/10/11/12/13/14/15/all versions)
  - MySQL/MariaDB/Percona (5.7/8.0/8.1/all versions)
- Flexible data faking based on Go templates and [Sprig template’s library](https://masterminds.github.io/sprig/) like [Helm](https://helm.sh/docs/chart_template_guide/functions_and_pipelines/). You may also use values of other columns for same row to build more flexible rules
- Stream data processing. It means that you can a use the tool through a pipe in command line and redirect dump from source DB directly to the destination DB with required transformations
- Easy to integrate into your CI/CD

### Who can use the tool

Development and project teams which are dealing with production and test/dev/stage or dynamic namespaces with databases and need to ensure security and prevent data leaks.

## Quickstart

Inspect your database structure and [set up](#settings) the nxs-data-anonymizer config in accordance with the sensitive data you need to anonymize. 

You are able to use this tool in any way you want. Three most common ways are described below.

#### Console

In order to operate with your database anonymization via console you need to go through the following steps:
- Download and untar the nxs-data-anonymizer [binary](https://github.com/nixys/nxs-data-anonymizer/releases)
- Run the nxs-data-anonymizer through the command line with [arguments](#command-line-arguments) you want to use

For example, use the following command if you need to anonymize your PostgreSQL database from production to dev on fly (PostgreSQL Client need to be installed):
```console
export PGPASSWORD=password; pg_dump -U postgres prod | /path/to/nxs-data-anonymizer -t pgsql -c /path/to/nxs-data-anonymizer.conf | psql -U postgres dev
```

#### GitLab CI

This section describes how to integrate nxs-data-anonymizer into your GitLab CI. You may add jobs presented below into your `.gitlab-ci.yml` and adjust it for yourself.

##### Job: anonymize prod

Job described in this section is able to perform the following tasks:
- Run when special tag for `main` branch is set
- Create a `production` database dump, anonymize and upload it into s3 bucket

Job sample:
```yaml
anonymize:
  stage: anonymize
  image: nixyslab/nxs-data-anonymizer:latest
  variables:
    GIT_STRATEGY: none
    PG_HOST: ${PG_HOST_PROD}
    PG_USER: ${PG_USER_PROD}
    PGPASSWORD: ${PG_PASS_PROD}
  before_script: 
  - echo "${S3CMD_CFG}" > ~/.s3cmd
  - echo "${NXS_DA_CFG}" > /nxs-data-anonymizer.conf
  script:
  - pg_dump -h ${PG_HOST} -U ${PG_USER} --schema=${PG_SCHEMA} ${PG_DATABASE} | /nxs-data-anonymizer -t pgsql -c /nxs-data-anonymizer.conf | gzip | s3cmd put - s3://bucket/anondump.sql.gz
  only:
  - /^v.*$/
  except:
  - branches
  - merge_requests
```

##### Job: update stage

Job described in this section deals with the following:
- Manual job for `stage` branch
- Download the anonymized dump from s3 bucket and load into `stage` database

Job sample:
```yaml
restore-stage:
  stage: restore
  image: nixyslab/nxs-data-anonymizer:latest
  variables:
    GIT_STRATEGY: none
    PG_HOST: ${PG_HOST_STAGE}
    PG_USER: ${PG_USER_STAGE}
    PGPASSWORD: ${PG_PASS_STAGE}
  before_script: 
  - echo "${S3CMD_CFG}" > ~/.s3cmd
  script:
  - s3cmd --no-progress --quiet get s3://bucket/anondump.sql.gz - | gunzip | psql -h ${PG_HOST} -U ${PG_USER} --schema=${PG_SCHEMA} ${PG_DATABASE}
  only:
  - stage
  when: manual
```

##### CI/CD variables

This section contains a description of CI/CD variables used in GitLab CI job samples above.

###### General

| Variable | Description |
| :---: | :---: |
|`S3CMD_CFG` | S3 storage config |
|`PG_SCHEMA`| PgSQL schema |
|`PG_DATABASE`|PgSQL database name|

###### Production

| Variable | Description |
| :---: | :---: |
|`NXS_DA_CFG`|nxs-data-anonymizer config|
|`PG_HOST_PROD` |PgSQL host|
|`PG_USER_PROD`|PgSQL user|
|`PG_PASS_PROD`|PgSQL password|

###### Stage

| Variable | Description |
| :---: | :---: |
|`PG_HOST_STAGE`|PgSQL host|
|`PG_USER_STAGE`|PgSQL user|
|`PG_PASS_STAGE`|PgSQL password|

#### GitHub Actions

This section describes how to integrate nxs-data-anonymizer into your GitHub Actions. You may add jobs presented below into your `.github` workflows and adjust it for yourself.

##### Job: anonymize prod

Job described in this section is able to perform the following tasks:
- Run when special tag is set
- Create a `production` database dump, anonymize and upload it into s3 bucket

```yaml
on:
  push:
    tags:
    - v*.*

jobs:
  anonymize:
    runs-on: ubuntu-latest
    container:
      image: nixyslab/nxs-data-anonymizer:latest
      env:
        PG_HOST: ${{ secrets.PG_HOST_PROD }}
        PG_USER: ${{ secrets.PG_USER_PROD }}
        PGPASSWORD: ${{ secrets.PG_PASS_PROD }}
        PG_SCHEMA: ${{ secrets.PG_SCHEMA }}
        PG_DATABASE: ${{ secrets.PG_DATABASE }}
    steps:
    - name: Create services configs
      run: |
        echo "${{ secrets.S3CMD_CFG }}" > ~/.s3cmd
        echo "${{ secrets.NXS_DA_CFG }}" > /nxs-data-anonymizer.conf
    - name: Anonymize
      run: |
        pg_dump -h ${PG_HOST} -U ${PG_USER} --schema=${PG_SCHEMA} ${PG_DATABASE} | /nxs-data-anonymizer -t pgsql -c /nxs-data-anonymizer.conf | gzip | s3cmd put - s3://bucket/anondump.sql.gz
```

##### Job: update stage

Job described in this section deals with the following:
- Manual job
- Download the anonymized dump from s3 bucket and load into `stage` database

```yaml
on: workflow_dispatch

jobs:
  restore-stage:
    runs-on: ubuntu-latest
    container:
      image: nixyslab/nxs-data-anonymizer:latest
      env:
        PG_HOST: ${{ secrets.PG_HOST_STAGE }}
        PG_USER: ${{ secrets.PG_USER_STAGE }}
        PGPASSWORD: ${{ secrets.PG_PASS_STAGE }}
        PG_SCHEMA: ${{ secrets.PG_SCHEMA }}
        PG_DATABASE: ${{ secrets.PG_DATABASE }}
    steps:
    - name: Create services configs
      run: |
        echo "${{ secrets.S3CMD_CFG }}" > ~/.s3cmd
    - name: Restore
      run: |
        s3cmd --no-progress --quiet get s3://bucket/anondump.sql.gz - | gunzip | psql -h ${PG_HOST} -U ${PG_USER} --schema=${PG_SCHEMA} ${PG_DATABASE}
```

##### GitHub Actions secrets

This section contains a description of secrets used in GitHub Actions job samples above.

###### General

| Variable | Description |
| :---: | :---: |
|`S3CMD_CFG` | S3 storage config |
|`PG_SCHEMA`| PgSQL schema |
|`PG_DATABASE`|PgSQL database name|

###### Production

| Variable | Description |
| :---: | :---: |
|`NXS_DA_CFG`|nxs-data-anonymizer config|
|`PG_HOST_PROD` |PgSQL host|
|`PG_USER_PROD`|PgSQL user|
|`PG_PASS_PROD`|PgSQL password|

###### Stage

| Variable | Description |
| :---: | :---: |
|`PG_HOST_STAGE`|PgSQL host|
|`PG_USER_STAGE`|PgSQL user|
|`PG_PASS_STAGE`|PgSQL password|

### Settings

Default configuration file path: `/nxs-data-anonymizer.conf`. The file is represented in yaml.

#### Command line arguments

| Argument         | Short   | Required | Having value | Default value | Description                                                      |
| :---:            | :---:   | :---:    | :---:     | :---:         |---                                                               |
| `--help`      | `-h` | No   | No    |  -      | Show program help message |
| `--version`      | `-v` | No   | No    |  -      | Show program version |
| `--conf`      | `-c` | No   | Yes    |  `/nxs-data-anonymizer.conf`      | Configuration file path |
| `--input`      | `-i` | No   | Yes    |  -     | File to read data from. If not specified `stdin` will be used |
| `--log-format` | `-l` | No   | Yes    |  `json`     | Log file format. You are available to use either `json` or `plain` value |
| `--output`      | `-o` | No   | Yes    |  -     | File to write data to. If not specified `stdout` will be used |
| `--type`      | `-t` | Yes   | Yes    |  -     | Database dump file type. Available values: `pgsql`, `mysql` |

#### General settings

| Option         | Type   | Required | Default value | Description                                                      |
|---             | :---:  | :---:    | :---:         |---                                                               |
| `logfile`      | String | No       | `stderr`      | Log file path. You may also use `stdout` and `stderr` |
| `loglevel`     | String | No       | `info`        | Log level. Available values: `debug`, `warn`, `error` and `info` |
| `progress`     | [Progress](#progress-settings) | No | - | Anonymization progress logging |
| `filters`          | Map of [Filters](#filters-settings) | No      | -             | Filters set for specified tables (key as a table name). Note: for PgSQL you also need to specify a scheme (e.g. `public.tablename`) |

##### Progress settings

| Option         | Type   | Required | Default value | Description                                                      |
|---             | :---:  | :---:    | :---:         |---                                                               |
| `rhythm`      | String | No       | `0s`    | Frequency write into the log a read bytes count. Progress will be written to the log only when this option is specified and has none-zero value. You may use a human-readable values (e.g. `30s`, `5m`, etc) |
| `humanize`      | Bool | No       | `false`    | Set this option to `true` if you need to write into the log a read bytes count in a human-readable format. On `false` raw bytes count will be written to the log |

##### Filters settings

Filters description for specified table.

| Option         | Type   | Required | Default value | Description                                                      |
|---             | :---:  | :---:    | :---:         |---                                                               |
| `columns`      | Map of [Columns](#columns-settings) | No       | -      | Filter rules for specified columns of table (key as a column name) |

###### Columns settings

| Option        | Type   | Required | Default value | Description                                                      |
|---            | :---:  | :---:    | :---:         |---                                                               |
| `type`        | String | No       | `template`    | Type of field `value`: `template` and `command` are available  |
| `value`       | String | Yes      | -             | The value to be used to replace at every cell in specified column. In accordance with the `type` this value may be either `Go template` or `command`. See below for details|
| `unique`      | Bool   | No       | `false`       | If true checks the generated value for cell is unique whole the column |

**Go template**

To anonymize a database fields you may use a Go template with the [Sprig template library's](https://masterminds.github.io/sprig/) functions. You may also use values of other columns in the rules for same row (with values before substitutions).

Additional filter functions:
- `null`: set a field value to `NULL`
- `isNull`: compare a field value with `NULL`

**Command**

To anonymize a database fields you may use a commands (scripts or binaries) with any logic you need. The command's concept has following properties:
- The command's `stdout` will be used as a new value for the anonymized field
- Command must return zero exit code, otherwise nxs-data-anonymizer will falls with error (in this case `stderr` will be used as an error text)
- Environment variables with the row data are available within the command:
  - `ENVVARTABLE`: contains a name of the filtered table
  - `ENVVARCURCOLUMN`: contains the current column name 
  - `ENVVARCOLUMN_{COLUMN_NAME}`: contains values (before substitutions) for all columns for the current row

#### Example

Imagine you have a simple table `users` in your production PgSQL like this:

| id | username | password | api_key |
| :---: | :---: | :---: | :---: |
| 1 | `admin` | `ZjCX6wUxtXIMtip` | `epezyj0cj5rqrdtxklnzxr3f333uibtz6avek7926141t1c918` |
| 2 | `alice` | `tuhjLkgwwetiwf8` | `2od4vfsx2irj98hgjaoi6n7wjr02dg79cvqnmet4kyuhol877z` |
| 3 | `bob`   | `AjRzvRp3DWo6VbA` | `owp7hob5s3o083d5hmursxgcv9wc4foyl20cbxbrr73egj6jkx` |

You need to get a dump with fake values:
- For `admin`: preset fixed value for a password and API key to avoid the need to change an app settings in your dev/test/stage or local environment after downloading the dump
- For others: usernames in format `user_N` (where `N` it is a user ID) and unique random  passwords and API keys

In accordance with these conditions, the nxs-data-anonymizer config may look like this:

```yaml
filters:
  public.users:
    columns:
      username:
        value: "{{ if eq .Values.username \"admin\" }}{{ .Values.username }}{{ else }}user_{{ .Values.id }}{{ end }}"
      password:
        type: command
        value: /path/to/script.sh
        unique: true
      api_key:
        value: "{{ if eq .Values.username \"admin\" }}preset_admin_api_key{{ else }}{{- randAlphaNum 50 | nospace | lower -}}{{ end }}"
        unique: true
```

The `/path/to/script.sh` script content is following:

```bash
#!/bin/bash

# Print preset password if current user is admin
if [ "$ENVVARCOLUMN_username" == "admin" ];
then
    echo -n "preset_admin_password"
    exit 0
fi

# Generate password for other users
p=$(pwgen -s 5 1 2>&1) 
if [ ! $? -eq 0 ];
then

    # On error print message to stderr and exit with non zero code

    echo -n "$p" >&2
    exit 1
fi

# Print generated password
echo $p | tr -d '\n'

exit 0
```

Now you may execute the following command in order to load anonymized data into your dev DB:

```
pg_dump ... | ./nxs-data-anonymizer -c filters.conf | psql -h localhost -U user example
```

As a result:
| id | username | password | api_key |
| :---: | :---: | :---: | :---: |
| 1 | `admin` | `preset_admin_password` | `preset_admin_api_key` |
| 2 | `user_2` | `Pp4HY` | `dhx4mccxyd8ux5uf1khpbqsws8qqeqs4efex1vhfltzhtjcwcu` |
| 3 | `user_3`   | `vu5TW` | `lgkkq3csskuyew8fr52vfjjenjzudokmiidg3cohl2bertc93x` |

It's easy.

## Roadmap

Following features are already in backlog for our development team and will be released soon:
- Global variables with the templated values you may use through the filters for all tables and columns
- Ability to delete tables and rows from faked dump 
- Ability to output into log a custom messages. It’s quite useful it order to obtain some generated data like admin passwords, etc
- Support of a big variety of databases

## Feedback

For support and feedback please contact me:
- telegram: [@borisershov](https://t.me/borisershov)
- e-mail: b.ershov@nixys.io

For news and discussions subscribe the channels:
- Telegram community (news): [@nxs_data_anonymizer](https://t.me/nxs_data_anonymizer)
- Telegram community (chat): [@nxs_data_anonymizer_chat](https://t.me/nxs_data_anonymizer_chat)

## License

nxs-data-anonymizer is released under the [Apache License 2.0](LICENSE).
