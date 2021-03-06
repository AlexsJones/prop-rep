# Prop-rep

[![Maintainability](https://api.codeclimate.com/v1/badges/0d8492035cb7a262ef8b/maintainability)](https://codeclimate.com/github/AlexsJones/prop-rep/maintainability)

A golang command line tool to show you which organisations are contributing to a repository.

Get the latest binary [here](https://github.com/AlexsJones/prop-rep/releases/)

![demo](./images/demo.gif)

```
prop-rep scan -r prometheus/prom2json -v
17 contributors found
This will take a while...
17/17 Users scanned
24 unique organisations found
Organisation             Contributors  Users
prometheus               6             beorn7, SuperQ, simonpasquier, juliusv, discordianfish, roidelapluie
FOSDEM                   1             SuperQ
open-cluster-management  1             chaitanyaenr
voxpupuli                1             roidelapluie
noisetor                 1             SuperQ
```

```
prop-rep scan -r prometheus/prom2json -g 
open graph.png
```

![graph](images/graph.png)

## Dependencies

Generate a Personal Access Token at [https://github.com/settings/tokens](https://github.com/settings/tokens) which needs the `read:org` scope. Use this token as `GITHUB_TOKEN` environment variable, which is required to access the GitHub API.

## Building

Linux users may find `stdio.h` is missing as an error message. 
For platforms such as Ubuntu you can circumvent this by installing `sudo apt install libc6-dev`

## Caveats

- Organisations a user is displayed as a member of will only be visible if they are a public member of that organisation ( https://github.com/orgs/<ORG>/people).
