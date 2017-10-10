[![travis Status](https://travis-ci.org/mattstratton/bowie.svg?branch=master)](https://travis-ci.org/mattstratton/bowie) [![Build status](https://ci.appveyor.com/api/projects/status/u7pu7ins2csxngxu?svg=true)](https://ci.appveyor.com/project/mattstratton/bowie)
 [![Go Report Card](https://goreportcard.com/badge/github.com/mattstratton/bowie)](https://goreportcard.com/report/github.com/mattstratton/bowie) [![GoDoc](https://godoc.org/github.com/mattstratton/bowie?status.svg)](http://godoc.org/github.com/mattstratton/bowie) [![GitHub release](https://img.shields.io/github/release/mattstratton/bowie.svg)](https://github.com/mattstratton/bowie/releases) [![MIT License](http://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE) <br />
[![Ebert](https://ebertapp.io/github/mattstratton/bowie.svg)](https://ebertapp.io/github/mattstratton/bowie) [![Codacy Badge](https://api.codacy.com/project/badge/Grade/c3f1eeb0b0bb4c68ba94df95cffefe0c)](https://www.codacy.com/app/matt.stratton/bowie?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=mattstratton/bowie&amp;utm_campaign=Badge_Grade) [![BCH compliance](https://bettercodehub.com/edge/badge/mattstratton/bowie?branch=master)](https://bettercodehub.com/) [![Coverage Status](https://coveralls.io/repos/github/mattstratton/bowie/badge.svg?branch=master)](https://coveralls.io/github/mattstratton/bowie?branch=master) [![codebeat badge](https://codebeat.co/badges/cbd7bfdf-e8d6-44b4-a377-20662bb2dbac)](https://codebeat.co/projects/github-com-mattstratton-bowie-master)

![bowie](bowie-logo.png)
> Time may change me<br>
> But I can't trace time<br>
> *-- David Bowie*

A pretty changelog generator 

Built with :heart: by [@mattstratton](https://github.com/mattstratton) in Go.

Inspired by [skywinder/github-changelog-generator](https://github.com/skywinder/github-changelog-generator).

This project adheres to the Contributor Covenant [code of conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. We appreciate your contribution. Please refer to the [contributing guidelines](CONTRIBUTING.md) for details on how to help.

<!-- TOC depthFrom:2 -->

- [Installation](#installation)
    - [Go](#go)
    - [Homebrew](#homebrew)
    - [Yum](#yum)
    - [Apt](#apt)
    - [Chocolatey](#chocolatey)
- [Usage](#usage)
    - [GitHub token](#github-token)
- [Contributing](#contributing)
- [Versioning](#versioning)
- [Authors](#authors)
- [Acknowledgments](#acknowledgments)
- [License](#license)

<!-- /TOC -->

## Installation

You can install `bowie` in multiple ways.

### Go

`go get github.com/mattstratton/bowie`

### Homebrew

`brew install ...` *TODO: Not yet available*

### Yum

`yum install ...` *TODO: Not yet available*

### Apt

`apt-get ...` *TODO: Not yet available*

### Chocolatey

`choco ...` *TODO: Not yet available*

## Usage

`bowie -u [username] -p [projectname]`

Where `[username]` is your GitHub user or organization (i.e., "mattstratton" or "google"), and `[projectname]` is the name of the GitHub project (i.e., "bowie").

If you doing specify a user and org, it will try to infer it based upon the git remote set as `origin` *TODO: This doesn't work yet.*

There will be other flags later. Also, remember to set the environment variable `BOWIE_GITHUB_TOKEN` or use the `-t` flag to pass it at the command line *(TODO: add details on the token)*

### GitHub token

`bowie` requires a GitHub token to function. 

Follow these steps:

- [Generate a token here](https://github.com/settings/tokens/new?description=GitHub%20Changelog%20Generator%20token) - you only need "repo" scope for private repositories
- Either:
    - Run the script with `--token <your-40-digit-token>`; **OR**
    - Set the `BOWIE_GITHUB_TOKEN` environment variable to your 40 digit token

You can set an environment variable by running the following command at the prompt, or by adding it to your shell profile (e.g., `~/.bash_profile` or `~/.zshrc`):

    export BOWIE_GITHUB_TOKEN="«your-40-digit-github-token»"


## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/mattstratton/bowie/tags). 

## Authors

* **Matt Stratton** - *Initial work* - [mattstratton](https://github.com/mattstratton)

## Acknowledgments

* [skywinder/github-changelog-generator](https://github.com/skywinder/github-changelog-generator)
* [goreleaser/goreleaser](https://github.com/goreleaser/goreleaser)
* [spf13/cobra](https://github.com/spf13/cobra)
* Logo via [David Bowie by James Fenton](https://thenounproject.com/term/david-bowie/128345/) from [the Noun Project](https://thenounproject.com/)
## License

bowie - A pretty changelog generator 

|                      |                                          |
|:---------------------|:-----------------------------------------|
| **Author:**          | Matt Stratton (<matt.stratton@gmail.com>)
| **Copyright:**       | Copyright 2017, Matt Stratton
| **License:**         | The MIT License

```markdown
The MIT License (MIT)

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

```
