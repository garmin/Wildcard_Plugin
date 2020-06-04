# Wildcard Plugin
This CF CLI Plugin allows users to search through and delete their applications using wildcards. It is useful for users who have multiple apps to manage in their spaces.

#Requirements
To prevent your shell from expanding the wildcard before the plugin sees it, wildcards should be escaped using a preceding '\\'.
```
cf wc-d app\*
cf wc-a app\?
```
# Installation

#### Install from github
##### OSX
`cf install-plugin http://github.com/garmin/Wildcard_Plugin/raw/1.0.4/bin/osx/wildcard_plugin`
##### WIN64
`cf install-plugin http://github.com/garmin/Wildcard_Plugin/raw/1.0.4/bin/win64/wildcard_plugin.exe`
##### LINUX
`cf install-plugin https://github.com/garmin/Wildcard_Plugin/raw/1.0.4/bin/linux64/wildcard_plugin`

#### Install from Source
```
go get code.cloudfoundry.org/cli
go get github.com/garmin/Wildcard_Plugin
cd $GOPATH/src/github.com/garmin/Wildcard_Plugin
go build *.go
cf install-plugin wildcard_plugin
```

## Usage

```
cf wildcard-apps APP_NAME_WITH_WILDCARD
```
```
cf wildcard-delete APP_NAME_WITH_WILDCARD [-f -r]
```

## Uninstall

```
cf uninstall-plugin wildcard
```
## Commands for wildcard-apps, wc-a

| Command/Option | Usage | Description|
| :--------------- |:---------------| :------------|
|`wildcard-apps, wc-a`| `cf wc-a APP_NAME_WITH_WILDCARD` |List all apps in the target space matching the wildcard pattern|

## Commands for wildcard-delete, wc-d

| Command/Option | Usage | Description|
| :--------------- |:---------------| :------------|
|`wildcard-delete, wc-d`| `cf wc-d APP_NAME_WITH_WILDCARD` |Displays list of matched apps and prompts the user for interactive deletion or force deletion of all matched apps|
|`-r`|`cf wc-d APP_NAME_WITH_WILDCARD -r`|Displays list of matched apps and prompts the user for interactive deletion or force deletion of all matched apps and their routes|
|`-f`|`cf wc-d APP_NAME_WITH_WILDCARD -f`|Force deletion of all apps in the target space matching the wildcard pattern without confirmation|
|`-f -r`|`cf wc-d APP_NAME_WITH_WILDCARD -f -r`|Force deletion of all apps and their mapped routes in the target space matching the wildcard pattern without confirmation|
