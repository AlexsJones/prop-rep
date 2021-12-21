# Prop-rep

A golang command line tool to show you which organisations are contributing to a repository.

```
prop-rep scan -r cncf/tag-app-delivery -v
44 contributors found
This will take a while...
44/44 Users scanned
57 unique organisations found
Contributors          Organisation  Users
kubernetes            5             scottrigby, resouer, hongchaodeng, feloy, onlydole
gitops-working-group  4             scottrigby, roberthstrand, todaywasawesome, stefanprodan
jenkinsci             3             caniszczyk, lloydchang, torstenwalter
oam-dev               2             resouer, hongchaodeng
open-gitops           2             scottrigby, todaywasawesome
prometheus-community  2             scottrigby, torstenwalter
weaveworks            2             scottrigby, stefanprodan
todogroup             1             caniszczyk
crayon                1             roberthstrand
teamserverless        1             stefanprodan
pantsbuild            1             caniszczyk
```

## Dependencies

A `GITHUB_TOKEN` environment variable is required to access the GitHub API.