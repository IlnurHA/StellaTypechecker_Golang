# Stella Typechecker on Go

```
Typecheck error:
Type Error: ERROR_UNEXPECTED_TYPE_FOR_PARAMETER
Unexpected type for parameter: a
Expected type:
        Nat
Actual type:
        Bool
in expression: fn (a : Bool) {
    return a;
  }
in function: main
```

## Build

- Install [go](https://go.dev/doc/install)

- Install [antlr4](https://www.antlr.org/)

- Clone repository:

    ```bash
    git clone https://github.com/IlnurHA/StellaTypechecker_Golang.git
    cd StellaTypechecker_Golang
    ```

- Generate parser:

    ```bash
    ./generate_parser.sh
    ```

- Build:

    ```bash
    go build cmd/typechecker/main.go
    ```

## Execution

### Command

```bash
./typecheck.sh [-dirPath=<path-to-directory>] [-filePath=<path-to-file>]
```

### Usage

>    `-dirPath` **string** Path to get tests from
>
>    `-filePath` **string** Path to source code on stella
>
>    `-noExitOnError` \[**bool**\] Skip exitting with status code 1 on type error

or execute it without flags to read from stdin

## Tests

If you are looking for tests you can find them here: https://github.com/ThatDraenGuy/stella_test_suite

Thanks @ThatDraenGuy for the repository
