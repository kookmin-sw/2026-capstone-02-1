# Imp

Imp acts as the input language of traceinspector. It's a simple imperative language that supports the following:

- integer, boolean, and dynamic array values
- if-else statements
- while loops
- function calls (Imp is pass-by-value, but arrays are passed by reference)

Imp also supports the following builtin functions:
- `scanf(fmt_string string, ...vars) -> None`: Writes values read from stdin specified by `fmt_string` into locations `vars`. `fmt_string` is `"%t"` or `"%d"`, and `vars` should be assignable expressions (variable names or array indexes).
- `print(...vals) -> None`: Prints to stdout variadic arguments `vals`. Imp does not have strings, but only for `print`, string literals may be passed. Note that newline is not automatically added.
- `make_array(size int, default val) -> array[val_ty]`: Returns an array of length `size` with values set as `default`
- `len(arr array[var]) -> int`: Returns the length of the array
- 
