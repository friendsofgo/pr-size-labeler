# Pull Request Size Labeler

Sample GitHub Action written in Go to visualize and optionally limit the size of your pull requests.

Developed for fun and learn. 
Heavily inspired on [codelytv/pr-size-labeler](https://github.com/CodelyTV/pr-size-labeler).

**All credits for the original authors&trade;**

[![Status](https://github.com/friendsofgo/pr-size-labeler/workflows/labeler/badge.svg)](https://github.com/friendsofgo/pr-size-labeler)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/friendsofgo/pr-size-labeler)](https://golang.org/doc/go1.15)
[![Version](https://img.shields.io/github/release/friendsofgo/pr-size-labeler.svg?style=flat-square)](https://github.com/friendsofgo/pr-size-labeler/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/friendsofgo/pr-size-labeler)](https://goreportcard.com/report/github.com/friendsofgo/pr-size-labeler)
[![Total alerts](https://img.shields.io/lgtm/alerts/g/friendsofgo/pr-size-labeler.svg?logo=lgtm&logoWidth=18)](https://lgtm.com/projects/g/friendsofgo/pr-size-labeler/alerts/)
[![FriendsOfGo](https://img.shields.io/badge/powered%20by-Friends%20of%20Go-73D7E2.svg)](https://friendsofgo.tech)

## Usage

Create a file named `labeler.yml` inside the `.github/workflows` directory and paste:

```yml
name: labeler

on: [pull_request]

jobs:
  labeler:
    runs-on: ubuntu-latest
    name: Label the PR size
    steps:
      - uses: friendsofgo/pr-size-labeler@v1.0
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          xs_max_size: '10'
          s_max_size: '100'
          m_max_size: '500'
          l_max_size: '1000'
          fail_if_xl: 'false'
          message_if_xl: 'This PR is so big! Please, split it 😊'
          files_to_ignore: 'go.mod *.js'
```

> If you want, you can customize all `*_max_size` with the size that fits in your project.

> By setting `fail_if_xl` to `'true'` you'll make fail all pull requests bigger than `l_max_size`.