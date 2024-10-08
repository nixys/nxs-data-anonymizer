# Example `variables`

An example of using nxs-data-anonymizer in a simple configuration with filters and link blocks.

Configuration block description:
```yaml
link:
- rule: # Description of the value generation rule.
    value: "some_value_generator" # Some function for generating values ​​with given parameters.
    unique: true # Flag of uniqueness of the value.
  with: # Description of tables with fields for which value coherence must be maintained.
    some_table1:
    - id
    some_table2:
    - id_from_table1
```

An example of a configuration:
```yaml
link:
- rule:
    value: "<generation_function_of_value>"
    unique: true
  with:
    table1:
    - col1_linked
    table2:
    - col2_linked

filters:
  table1:
    columns:
      col2:
        value: "<fixed_value_or_generation_function_of_value>"
      col3:
        value: "<fixed_value_or_generation_function_of_value>"
  table2:
    columns:
      col1:
        value: "<fixed_value_or_generation_function_of_value>"
      col3:
        value: "<fixed_value_or_generation_function_of_value>"
```

Working examples of configurations in the `./MySQL` and `./PostgreSQL` directories.