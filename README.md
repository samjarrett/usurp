# usurp

Usurp is a command-line shim to assume an AWS IAM Role, and pass the credentials to the subsequent command.

## Installation

`usurp` is provided as a statically compiled binary for Darwin (macOS, in amd64/x86_64 for traditional 
Intel macs, and arm64 for M1 macs), Linux (in amd64/x86_64 for Intel/AMD, and arm64 for ARM-based 
architectures) and Windows (amd64/x86_64 only).

The latest binaries are attached to [Github releases](https://github.com/samjarrett/usurp/releases).


## Usage

```bash
usurp <role arn> <command> [<args ...>]
```

For example:

```bash
$ usurp arn:aws:iam::123456789012:role/my-role-name aws sts get-caller-identity
```
which will assume `my-role-name` in account `123456789012` and execute `aws sts get-caller-identity`.
