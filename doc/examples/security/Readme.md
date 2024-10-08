# Example `security`

An example of using nxs-data-anonymizer in a simple configuration with filters and security block.

An example of a configuration:
```yaml
security:
  policy: # Policy for handling undeclared tables and columns.
    tables: skip
    columns: randomize
  exceptions: # Excludes policy processing for the specified tables and columns.
    tables: # These tables will not be skipped.
    - table2
    columns: # These columns will not be randomized.
    - excluded_col1
    - excluded_col2
  defaults: # Overrides the default randomization value for a column with the specified value. If you use a function to generate a value, that value will be the same for all substitutions.
    columns:
      default_col1: # Column name to override the default value.
        value: "<fixed_value_or_generation_function_of_value>"

filters:
  table1:
    columns:
      table1_col1:
        value: "<fixed_value_or_generation_function_of_value>"
      table1_col2:
        value: "<fixed_value_or_generation_function_of_value>"
```

Working examples of configurations in the `./MySQL` and `./PostgreSQL` directories.