![photo_2023-07-28_15-28-52](https://github.com/nixys/nxs-data-anonymizer/assets/27485608/165a90a0-929f-460b-8dbd-2903c0d91f36)

# nxs-data-anonymizer

[![Telegram News][tg-news-badge]][tg-news-url]
[![Telegram Chat][tg-chat-badge]][tg-chat-url]

## Introduction

nxs-data-anonymizer is a tool for anonymizing **PostgreSQL** and **MySQL/MariaDB/Percona** databases' dump.

### Features

- Supported databases and versions:
  - PostgreSQL (9/10/11/12/13/14/15/all versions)
  - MySQL/MariaDB/Percona (5.7/8.0/8.1/all versions)
- Flexible data faking based on:
  - Go templates and [Sprig template’s library](https://masterminds.github.io/sprig/) like [Helm](https://helm.sh/docs/chart_template_guide/functions_and_pipelines/). You may also use values of other columns for same row to build more flexible rules
  - External commands you may execute to create table field values
  - Security enforcement rules
  - Link cells across the database to generate the same values
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
| `variables`     | Map of [Variables](#variables-settings) (key: variable name) | No | - | Global variables to be used in a filters. Variables are set at the init of application and remain unchanged during the runtime  |
| `link`     | Slice of [Link](#link-settings) | No | - | Rules to link specified columns across the database  |
| `filters`          | Map of [Filters](#filters-settings) (key: table name) | No      | -             | Filters set for specified tables (key as a table name). Note: for PgSQL you also need to specify a scheme (e.g. `public.tablename`) |
| `security`     | [Security](#security-settings) | No       | -      | Security enforcement for anonymizer |


##### Progress settings

| Option         | Type   | Required | Default value | Description                                                      |
|---             | :---:  | :---:    | :---:         |---                                                               |
| `rhythm`      | String | No       | `0s`    | Frequency write into the log a read bytes count. Progress will be written to the log only when this option is specified and has none-zero value. You may use a human-readable values (e.g. `30s`, `5m`, etc) |
| `humanize`      | Bool | No       | `false`    | Set this option to `true` if you need to write into the log a read bytes count in a human-readable format. On `false` raw bytes count will be written to the log |

##### Variables settings

| Option        | Type   | Required | Default value | Description                                                      |
|---            | :---:  | :---:    | :---:         |---                                                               |
| `type`        | String | No       | `template`    | Type of field `value`: `template` and `command` are available  |
| `value`       | String | Yes      | -             | The value to be used as global variable value within the filters. In accordance with the `type` this value may be either `Go template` or `command`. See below for details|

##### Link settings

Link is used to create the same data with specified rules for different cells across the database.

Each link element has following properties:
- Able to contain multiple tables and columns for each table
- All specified cells with the same data before anonymization will have same data after 
- One common rule to generate new values

| Option        | Type   | Required | Default value | Description                                                      |
|---            | :---:  | :---:    | :---:         |---                                                               |
| `type`        | String | No       | `template`    | Type of field `value`: `template` and `command` are available  |
| `value`       | String | Yes      | -             | The value to be used to replace at every cell in specified column. In accordance with the `type` this value may be either `Go template` or `command`. See below for details|
| `unique`      | Bool   | No       | `false`       | If true checks the generated value for cell is unique whole an all columns specified for `link` element |

##### Filters settings

Filters description for specified table.

| Option         | Type   | Required | Default value | Description                                                      |
|---             | :---:  | :---:    | :---:         |---                                                               |
| `columns`      | Map of [Columns](#columns-settings) (key: column name) | No       | -      | Filter rules for specified columns of table (key as a column name) |

###### Columns settings

| Option        | Type   | Required | Default value | Description                                                      |
|---            | :---:  | :---:    | :---:         |---                                                               |
| `type`        | String | No       | `template`    | Type of field `value`: `template` and `command` are available  |
| `value`       | String | Yes      | -             | The value to be used to replace at every cell in specified column. In accordance with the `type` this value may be either `Go template` or `command`. See below for details|
| `unique`      | Bool   | No       | `false`       | If true checks the generated value for cell is unique whole the column |

**Go template**

To anonymize a database fields you may use a Go template with the [Sprig template library's](https://masterminds.github.io/sprig/) functions. 

Additional filter functions:
- `null`: set a field value to `NULL`
- `isNull`: compare a field value with `NULL`
- `drop`: drop whole row. If table has filters for several columns and at least one of them returns drop value, whole row will be skipped during the anonymization process

You may also use the following data in a templates:
- Current table name. Statement: `{{ .TableName }}`
- Current column name. Statement: `{{ .CurColumnName }}`
- Values of other columns in the rules for same row (with values before substitutions). Statement: `{{ .Values.COLUMN_NAME }}` (e.g.: `{{ .Values.username }}`)
- Global variables. Statement: `{{ .Variables.VARIABLE_NAME }}` (e.g.: `{{ .Variables.password }}`)  
- Raw column data type. Statement: `{{ .ColumnTypeRaw }}`
- Regex's capturing groups for the column data type. This variable has array type so you need to use `range` or `index` to access specific element. Statement: `{{ index .ColumnTypeGroups 0 0 }}`. See [Types](#types-settings) for details

**Command**

To anonymize a database fields you may use a commands (scripts or binaries) with any logic you need. The command's concept has following properties:
- The command's `stdout` will be used as a new value for the anonymized field
- Command must return zero exit code, otherwise nxs-data-anonymizer will falls with error (in this case `stderr` will be used as an error text)
- Environment variables with the row data are available within the command:
  - `ENVVARTABLE`: contains a name of the current table
  - `ENVVARCURCOLUMN`: contains the current column name
  - `ENVVARCOLUMN_{COLUMN_NAME}`: contains values (before substitutions) for all columns for the current row
  - `ENVVARGLOBAL_{VARIABLE_NAME}`: contains value for specified global variable
  - `ENVVARCOLUMNTYPERAW`: contains raw column data type
  - `ENVVARCOLUMNTYPEGROUP_{GROUP_NUM}_{SUBGROUPNUM}`: contains regex's capturing groups for the column data type. See [Types](#types-settings) for details

##### Security settings

| Option         | Type   | Required | Default value | Description                                                      |
|---             | :---:  | :---:    | :---:         |---                                                               |
| `policy`      | [Policy](#policy-settings) | No       | -      | Security policy for entities |
| `exceptions`      | [Exceptions](exceptions-settings) | No       | -      | Exceptions for entities |
| `defaults`      | [Defaults](defaults-settings) | No       | -      | Default filters for entities |

###### Policy settings

| Option         | Type   | Required | Default value | Description                                                      |
|---             | :---:  | :---:    | :---:         |---                                                               |
| `tables`      | String | No       | `pass`      | Security policy for tables. If value `skip` is used all undescribed tables in config will be skipped while anonymization |
| `columns`      | String | No       | `pass`      | Security policy for columns. If value `randomize` is used all undescribed columns in config will be randomized (with default rules in accordance to types) while anonymization |

_Values to masquerade a columns in accordance with the types see below._

**PgSQL:**

| Type | Value to masquerade |
|---|:---:|
| `smallint`    | `0` |
| `integer`     | `0` |
| `bigint`      | `0` |
| `smallserial` | `0` |
| `serial`      | `0` |
| `bigserial`   | `0` |
| `decimal`     | `0.0` |
| `numeric`     | `0.0` |
| `real`        | `0.0` |
| `double`      | `0.0` |
| `character`   | `randomized character data"` |
| `bpchar`      | `randomized bpchar data` |
| `text`        | `randomized text data` |

**MySQL:**

| Type | Value to masquerade |
|---|:---:|
| `bit` |              `0` |
| `bool` |             `0` |
| `boolean` |          `0` |
| `tinyint` |          `0` |
| `smallint` |         `0` |
| `mediumint` |        `0` |
| `int` |              `0` |
| `integer` |          `0` |
| `bigint` |           `0` |
| `float` |            `0.0` |
| `double` |           `0.0` |
| `double precision` | `0.0` |
| `decimal` |          `0.0` |
| `dec` |              `0.0` |
| `char` |       `randomized char` (String will be truncated to "COLUMN_SIZE" length.)| 
| `varchar` |    `randomized varchar` (String will be truncated to "COLUMN_SIZE" length.) | 
| `tinytext` |   `randomized tinytext` | 
| `text` |       `randomized text` | 
| `mediumtext` | `randomized mediumtext` | 
| `longtext` |   `randomized longtext` | 
| `enum` |       Last value from `enum` | 
| `set` |        Last value from `set` | 
| `date` |       `2024-01-01` | 
| `datetime` |   `2024-01-01 00:00:00` | 
| `timestamp` |  `2024-01-01 00:00:00` | 
| `time` |       `00:00:00` | 
| `year` |       `2024` | 
| `json` |       `{"randomized": "json_data"}` |
| `binary` |     `cmFuZG9taXplZCBiaW5hcnkgZGF0YQo=` |
| `varbinary` |  `cmFuZG9taXplZCBiaW5hcnkgZGF0YQo=` |
| `tinyblob` |   `cmFuZG9taXplZCBiaW5hcnkgZGF0YQo=` |
| `blob` |       `cmFuZG9taXplZCBiaW5hcnkgZGF0YQo=` |
| `mediumblob` | `cmFuZG9taXplZCBiaW5hcnkgZGF0YQo=` |
| `longblob` |   `cmFuZG9taXplZCBiaW5hcnkgZGF0YQo=` |

###### Exceptions settings

| Option         | Type   | Required | Default value | Description                                                      |
|---             | :---:  | :---:    | :---:         |---                                                               |
| `tables`      | Slice of strings | No       | -      | Table names without filters which are not be skipped while anonymization if option `security.policy.tables` set to `skip` |
| `columns`      | Slice of strings | No       | -      | Column names (in any table) without filters which are not be randomized while anonymization if option `security.policy.columns` set to `randomize` |

###### Defaults settings

| Option         | Type   | Required | Default value | Description                                                      |
|---             | :---:  | :---:    | :---:         |---                                                               |
| `columns`      | Map of Filters | No       | -      | Default filter for columns (in any table). That filters will be applied for columns with this names without described filters |
| `types`      | Slice of [Types](#types-settings) | No       | -      | Custom filters for types (in any table). With this filter rules you may override default filters for types |

###### Types settings

| Option         | Type   | Required | Default value | Description                                                      |
|---             | :---:  | :---:    | :---:         |---                                                               |
| `regex`      | String | Yes       | -      | Regular expression. Will be checked for match for column data type (in `CREATE TABLE` section). Able to use capturing groups within the regex that available as an additional variable data in the filters (see [Columns](#columns-settings) for details). This ability helps to create more flexible rules to generate the cells value in accordance with  data type  features |
| `rule`      | [Columns](#columns-settings) | Yes       | -      | Rule will be applied columns with data types matched for specified regular expression |

#### Example

Imagine you have a simple database with two tables `users` and `posts` in your production PgSQL like this:

| id | username | password | api_key |
| :---: | :---: | :---: | :---: |
| 1 | `admin` | `ZjCX6wUxtXIMtip` | `epezyj0cj5rqrdtxklnzxr3f333uibtz6avek7926141t1c918` |
| 2 | `alice` | `tuhjLkgwwetiwf8` | `2od4vfsx2irj98hgjaoi6n7wjr02dg79cvqnmet4kyuhol877z` |
| 3 | `bob`   | `AjRzvRp3DWo6VbA` | `owp7hob5s3o083d5hmursxgcv9wc4foyl20cbxbrr73egj6jkx` |

| id | poster_id | title | content |
| :---: | :---: | :---: | :---: |
| 1 | 1 | `example_post_1` | `epezyj0cj5rqrdtxklnzxr3f333uibtz6avek7926141t1c918` |
| 2 | 2 | `example_post_2` | `2od4vfsx2irj98hgjaoi6n7wjr02dg79cvqnmet4kyuhol877z` |
| 3 | 3 | `example_post_3` | `owp7hob5s3o083d5hmursxgcv9wc4foyl20cbxbrr73egj6jkx` |
| 4 | 1 | `example_post_4` | `epezyj0cj5rqrdtxklnzxr3f333uibtz6avek7926141t1c918` |
| 5 | 2 | `example_post_5` | `2od4vfsx2irj98hgjaoi6n7wjr02dg79cvqnmet4kyuhol877z` |
| 6 | 3 | `example_post_6` | `owp7hob5s3o083d5hmursxgcv9wc4foyl20cbxbrr73egj6jkx` |

You need to get a dump with fake values:
- For `admin`: preset fixed value for a password and API key to avoid the need to change an app settings in your dev/test/stage or local environment after downloading the dump.
- For others: usernames in format `user_N` (where `N` it is a user ID) and unique random  passwords and API keys.
- Need to preserve data mapping between `users` and `posts` tables in `id` and `poster_id` columns
- Need to randomize contents of `content` column.
In accordance with these conditions, the nxs-data-anonymizer config may look like this:

```yaml
variables:
#Global variables.
  adminPassword:
    type: template
    value: "preset_admin_password"
  adminAPIKey:
    value: "preset_admin_api_key"

#Block defining rules of behavior with fields and tables for which filters are not specified.
security:
# Specifies the required actions for tables and columns that are not specified in the configuration.
  policy:
    tables: skip
    columns: randomize
# Excludes policy actions for the specified tables and columns.
  exceptions:
    tables: 
    - public.posts
    columns:
    - title
# Overrides the default policy actions for the columns specified in this block. The value is generated once and substituted into all instances of the field.
  defaults:
    columns:
      content:
        value: "{{- randAlphaNum 20 -}}"

#Here you define the rules that allow you to preserve the mapping of values ​​between tables.
link:
- rule:
#Value generation rule.
    value: "{{ randInt 1 15	}}"
    unique: true
  with:
#Tables and columns to which the rule is applied.
    public.users:
    - id
    public.posts:
    - poster_id

#Block describing replacement rules for fields.
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
        value: "{{ if eq .Values.username \"admin\" }}{{ .Variables.adminAPIKey }}{{ else }}{{- randAlphaNum 50 | nospace | lower -}}{{ end }}"
        unique: true
```

The `/path/to/script.sh` script content is following:

```bash
#!/bin/bash

# Print preset password if current user is admin
if [ "$ENVVARCOLUMN_username" == "admin" ];
then
    echo -n "$ENVVARGLOBAL_adminPassword"
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
| 5 | `admin` | `preset_admin_password` | `preset_admin_api_key` |
| 4 | `user_2` | `Pp4HY` | `dhx4mccxyd8ux5uf1khpbqsws8qqeqs4efex1vhfltzhtjcwcu` |
| 7 | `user_3`   | `vu5TW` | `lgkkq3csskuyew8fr52vfjjenjzudokmiidg3cohl2bertc93x` |

| id | poster_id | title | content |
| :---: | :---: | :---: | :---: |
| 1 | 5 | `example_post_1` | `EDlT6bGXJ2LOS7CE2E4b` |
| 2 | 4 | `example_post_2` | `EDlT6bGXJ2LOS7CE2E4b` |
| 3 | 7 | `example_post_3` | `EDlT6bGXJ2LOS7CE2E4b` |
| 4 | 5 | `example_post_4` | `EDlT6bGXJ2LOS7CE2E4b` |
| 5 | 4 | `example_post_5` | `EDlT6bGXJ2LOS7CE2E4b` |
| 6 | 7 | `example_post_6` | `EDlT6bGXJ2LOS7CE2E4b` |

It's easy. You can find more examples in doc/examples.

## Roadmap

Following features are already in backlog for our development team and will be released soon:
- [x] Global variables with the templated values you may use through the filters for all tables and columns
- [x] Ability to delete tables and rows from faked dump 
- [ ] Ability to output into log a custom messages. It’s quite useful it order to obtain some generated data like admin passwords, etc
- [ ] Support of a big variety of databases

## Feedback

For support and feedback please contact me:
- telegram: [@borisershov](https://t.me/borisershov)
- e-mail: b.ershov@nixys.io

For news and discussions subscribe the channels:
- Telegram community (news): [@nxs_data_anonymizer](https://t.me/nxs_data_anonymizer)
- Telegram community (chat): [@nxs_data_anonymizer_chat](https://t.me/nxs_data_anonymizer_chat)

## License

nxs-data-anonymizer is released under the [Apache License 2.0](LICENSE).

[tg-news-badge]: https://img.shields.io/endpoint?url=https%3A%2F%2Ftg.sumanjay.workers.dev%2Fnxs_data_anonymizer
[tg-chat-badge]: https://img.shields.io/endpoint?url=https%3A%2F%2Ftg.sumanjay.workers.dev%2Fnxs_data_anonymizer_chat
[tg-news-url]: https://t.me/nxs_data_anonymizer
[tg-chat-url]: https://t.me/nxs_data_anonymizer_chat
[aica-badge]: https://img.shields.io/badge/AI-Code%20Assist-EB9FDA
[aica-url]: https://app.commanddash.io/agent?github=https://github.com/nixys/nxs-data-anonymizer
