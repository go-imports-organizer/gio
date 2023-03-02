<!---
#
# Copyright 2023 Go Imports Organizer Contributors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# 	https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
-->
# gio
A customizable imports organizer for the Go programming language

* [Summary](#summary)
* [Installation](#installation)
* [Configuration](#configuration)

# <a name='summary'></a>Summary
`gio` is a fully customizable Go imports organizer. The configuration
is project based and is stored in a `gio.yaml` file in the root of your
module's project folder alongside the `go.mod` file. For consistency
the `gio.yaml` file should be committed to your projects vcs.

# <a name='installation'></a>Installation

## Command Line Tool

```
  $ go install github.com/go-imports-organizer/gio@latest
```

## Go project configuration

### Example scripts/tools.go file
This file will ensure that the `github.com/go-imports-organizer/gio` repo is vendored into your project.
```
//go:build tools
// +build tools

package hack

// Add tools that scripts depend on here, to ensure they are vendored.
import (
	_ "github.com/go-imports-organizer/gio"
)

```

### Example scripts/verify-imports.sh script
This file will check if there are any go files that need to be formatted. If there are, it will print a list of them, and exit with status one (1), otherwise it will exit with status zero (0). Make sure that you make the file executable with `chmod +x scripts/verify-imports.sh`.
```
#!/bin/bash

bad_files=$(go run ./vendor/github.com/go-imports-organizer/gio -l)
if [[ -n "${bad_files}" ]]; then
        echo "!!! gio needs to be run on the following files:"
        echo "${bad_files}"
        echo "Try running 'make imports'"
        exit 1
fi
```

### Example Makefile sections
```
imports: ## Organize imports in go files using gio. Example: make imports
	go run ./vendor/github.com/go-imports-organizer/gio
.PHONY: imports

verify-imports: ## Run import verifications. Example: make verify-imports
	hack/verify-imports.sh
.PHONY: verify-imports
```

# <a name='Configuration'></a>Configuration
The `gio.yaml` configuration file is a well formatted yaml file.

### Excludes
An array of Exclude definitions.

Each Exclude definition is a rule that tells `gio` which files and folders to ignore while looking for Go files that it should organize the imports of. The default configuration file ignores the `.git` directory and the `vendor` directory. The more files and folders that you can ignore at this stage the faster `gio` will run.

#### RegExp
A string, valid values are any valid Go regular expression.

A well formatted Regular Expression that is used to match against. Be as specific as possible.

#### MatchType
A string, valid values are `[name, path]`.

Lets `gio` know whether to match agains the files `name` or the files `path` relative to the modules root directory _(where the go.mod file is located)_.

### Groups
An array of Group definitions

Each group definition is a rule that tells `gio` how you would like the `imports` in your Go files organized. Each definition represents a block of import statements, describing the order in which that block should be displayed and how to identify the items that it should contain.

#### Description
A string, valid values are any valid string value

A friendly name to identify the definition by instead of trying to decipher the regular expression each time to remember what it does.

#### RegExp
A string, valid values are any valid Go regular expression.

A well formatted Regular Expression that is used to match against. Be as specific as possible.

Note:

There is one keyword that is available for the RegExp value that is a special keyword, it is `%{module}%`. This keyword automatically creates a regular expression that matches the current module name as defined by the go.mod file. To ensure that it captures the correct imports you should always set the `MatchOrder` to `0` for this definition.

#### MatchOrder
An integer, valid values are -n...n

Tells `gio` which order the definitions should be matched against in. Lower numbers are first, higher numbers are last.

It is important to ensure the correct `matchorder` is used, expecially if any of your `regexp` have any kind of overlap, such as having a module name of `github.com/example/mymodule` and a group definition for `github.com/example`. You would want to make sure that your `module` definition was matched first or those imports would get rolled into the `github.com/example` one because it is less specific.


#### DisplayOrder
An integer, valid values are -n...n

The order in which the blocks of imports should be displayed in each file. Lower numbers are first, higher numbers are last
