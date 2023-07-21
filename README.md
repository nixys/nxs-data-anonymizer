# nxs-data-anonymizer

## Introduction

nxs-data-anonymizer is a tool to anonymize a **PostgreSQL** and **MySQL** databases dump.

### Features

- PgSQL (_all versions_) and MySQL/MariaDB/Percona (_all versions_) are supported
- Flexible data faking based on Go templates and [Sprig template library](https://masterminds.github.io/sprig/) like the [Helm](https://helm.sh/docs/chart_template_guide/functions_and_pipelines/). Also you may use a values of other columns for same row to build more flexible rules
- Stream data processing. That's mean you can a use the tool through a pipe in command line and redirect dump from source DB directly to the destination DB with required transformations
- Easy to integrate into your CI/CD

### Who use the tool

Development teams and projects who has production and test/dev/stage or dynamic namespaces with databases and needs to be ensure for security and prevent data leaks.

## Quickstart

Inspect your database structure and [set up](#Settings) the nxs-data-anonymizer config in accordance with the sensitive data you need to anonymize. 

You are able to use this tool in any way you want. Three most common ways are described below.

#### Console

To operate with your database anonymization via console you need following:
- Download and untar the nxs-data-anonymizer [binary](https://github.com/nixys/nxs-data-anonymizer/releases)
- Run the nxs-data-anonymizer through the command line with an [arguments](#command-line-arguments) you want to use

For example, use the following command if you need to anonymize your PostgreSQL database from production to dev on fly (PostgreSQL Client need to be installed):
```console
export PGPASSWORD=password; pg_dump -U postgres prod | /path/to/nxs-data-anonymizer -t pgsql -c /path/to/nxs-data-anonymizer.conf | psql -U postgres dev
```

#### GitLab CI

This section describes how to integrate nxs-data-anonymizer into your GitLab CI. You may add jobs below into your `.gitlab-ci.yml` and adjust it for yourself.

##### Job: anonymize prod

Job described in this section do the following:
- Run when special tag for `main` branch is set
- Create a `production` database dump, anonymize it and updaload into s3 bucket

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

Job described in this section do the following:
- Manual job for `stage`branch
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

This section contains a description for CI/CD variables used in GitLab CI job samples above.

###### General

| Variable | Description |
| :---: | :---: |
| `S3CMD_CFG` | S3 storage config |
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

_Comming soon ..._

### Settings

Default configuration file path: `/nxs-data-anonymizer.conf`. File represented in yaml.

#### Command line arguments

| Argument         | Short   | Required | Has value | Default value | Description                                                      |
| :---:            | :---:   | :---:    | :---:     | :---:         |---                                                               |
| `--help`      | `-h` | No   | No    |  -      | Show program help message |
| `--version`      | `-v` | No   | No    |  -      | Show program version |
| `--conf`      | `-c` | No   | Yes    |  `/nxs-data-anonymizer.conf`      | Configuration file path |
| `--input`      | `-i` | No   | Yes    |  -     | File to read data from. If not specified `stdin` will be used |
| `--output`      | `-o` | No   | Yes    |  -     | File to write data to. If not specified `stdout` will be used |
| `--type`      | `-t` | Yes   | Yes    |  -     | Data base dump file type. Available values: `pgsql`, `mysql` |

#### General settings

| Option         | Type   | Required | Default value | Description                                                      |
|---             | :---:  | :---:    | :---:         |---                                                               |
| `logfile`      | String | No       | `stderr`      | Log file path. Also you may use `stdout` and `stderr` |
| `loglevel`     | String | No       | `info`        | Log level. Available values: `debug`, `warn`, `error` and `info` |
| `filters`          | Map of [Filters](#filters-settings) | No      | -             | Filters set for specified tables (key as a table name). Note: for PgSQL you also need to specify a scheme (e.g. `public.tablename`) |

##### Filters settings

Filters description for specified table.

| Option         | Type   | Required | Default value | Description                                                      |
|---             | :---:  | :---:    | :---:         |---                                                               |
| `columns`      | Map of [Columns](#columns-settings) | No       | -      | Filter rules for specified columns of table (key as a column name) |

###### Columns settings

| Option        | Type   | Required | Default value | Description                                                      |
|---            | :---:  | :---:    | :---:         |---                                                               |
| `value`       | String | No       | -             | A value to be used for replace at every cell in specified column. This value may be either fixed value or Go template with the [Sprig template library](https://masterminds.github.io/sprig/) functions. You also may use in the rules a values of other columns for same row (before substitutions) |
| `unique`      | Bool   | No       | `false`       | If true checks the generated value for cell ischemes unique whole the column |

#### Example

Imagine you have a simple table `users` in your production PgSQL like this:

| id | username | api_key |
| :---: | :---: | :---: |
| 1 | `admin` | `epezyj0cj5rqrdtxklnzxr3f333uibtz6avek7926141t1c918` |
| 2 | `alice` | `2od4vfsx2irj98hgjaoi6n7wjr02dg79cvqnmet4kyuhol877z` |
| 3 | `bob`   | `owp7hob5s3o083d5hmursxgcv9wc4foyl20cbxbrr73egj6jkx` |

You need to get a dump with fake values:
- For `admin`: preset fixed value for an API key to avoid a need to change an app settings in your dev/test/stage or local environment after downloading the dump
- For others: usernames in format `user_N` (where `N` it is a user ID) and unique random API keys

In accordance with this conditions, the nxs-data-anonymizer config may looks like this:

```yaml
filters:
  public.users:
    columns:
      username:
        value: "{{ if eq .Values.username \"admin\" }}{{ .Values.username }}{{ else }}user_{{ .Values.id }}{{ end }}"
      api_key:
        value: "{{ if eq .Values.username \"admin\" }}preset_admin_api_key{{ else }}{{- randAlphaNum 50 | nospace | lower -}}{{ end }}"
        unique: true
```

Now you may execute following command to load anonymized data into your dev DB:

```
pgdump ... | ./nxs-data-anonymizer -c filters.conf | psql -h localhost -U user example
```

As a result:
| id | username | api_key |
| :---: | :---: | :---: |
| 1 | `admin` | `preset_admin_api_key` |
| 2 | `user_2` | `dhx4mccxyd8ux5uf1khpbqsws8qqeqs4efex1vhfltzhtjcwcu` |
| 3 | `user_3`   | `lgkkq3csskuyew8fr52vfjjenjzudokmiidg3cohl2bertc93x` |

It's easy.

## Roadmap

Following features are already in backlog for our development team and will be released soon:
- Global variables with the templated values you may use through the filters for all tables and columns
- Ability to delete tables and rows from faked dump 
- Ability to output into log a custom messages. Useful to obtain some generated data like a admin passwords, etc
- Support for more databases

## Feedback

For support and feedbaсk please contact me:
- telegram: [@borisershov](https://t.me/borisershov)
- e-mail: b.ershov@nixys.ru

## License

nxs-data-anonymizer is released under the [Apache License 2.0](LICENSE).
