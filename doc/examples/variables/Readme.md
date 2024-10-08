# Example `variables`

An example of using nxs-data-anonymizer in a simple configuration with filters and variables blocks.

In order to use variables in the "filters" block, you need to describe the "variables" block:
```yaml
variables:
  some_variable: # Variable name
    value: "some_value" # Variable value
```

If you use a function to generate a value, then in this block it will be executed once and the resulting value will be the same for all its occurrences in the dump.

An example of a configuration:
```yaml
variables:
  var1:
    value: "<fixed_value_or_generation_function_of_value>"
  var2:
    value: "<fixed_value_or_generation_function_of_value>"
  var3:
    value: "<fixed_value_or_generation_function_of_value>"

filters:
  table1:
    columns:
      table1_col1:
        value: "<fixed_value_or_generation_function_of_value>"
      table1_col2:
        value: "{{ .Variables.var1 }}"
      table1_col3:
        value: "{{ .Variables.var2 }}"
  table2:
    columns:
      table2_col1:
        value: "<fixed_value_or_generation_function_of_value>"
      table2_col2:
        value: "{{ .Variables.var1 }}"
      table2_col3:
        value: "{{ .Variables.var3 }}"
```

Working examples of configurations in the `./MySQL` and `./PostgreSQL` directories.