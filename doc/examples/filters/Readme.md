# Example `filters`

An example of using nxs-data-anonymizer in a simple configuration with one filters block.

All that needs to be specified in the configuration file is a list of tables to be processed with a list of replacement rules for the columns that need to be changed.

An example of a configuration:
```yaml
filters:
  table1: # Table name to process.
   columns:
      table1_col1: # Column name to process.
        value: "{{ - value_generation_function - }}" # Value for substitution.
        unique: true # The unique values ​​flag, default "false", is used only with the value generation function.
      table1_col2:
        value: "<fixed_value_or_generation_function_of_value>"
      table1_col3:
        value: "<fixed_value_or_generation_function_of_value>"
  table2:
    columns:
      table2_col1:
        value: "{{ - value_generation_function - }}"
        unique: true
      table2_col2:
        value: "<fixed_value_or_generation_function_of_value>"
      table2_col3:
        value: "<fixed_value_or_generation_function_of_value>"
```

Working examples of configurations in the `./MySQL` and `./PostgreSQL` directories.