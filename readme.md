# Gitbao

Wrapping up your github gists into delicious cloud-shaped bao buns. 

## Update

Proof of concept is done. Arbitrary go applications can be run from gists. Next steps:

 - Better docker management. Disk space, cpu cycles, caching of dependencies and pre-building docker containers. Also better cleanup/removal of dockerfiles. 
 - Auto-update of applications from github api listener.
 - Parse user configuration file on boa page, also detect this if it's a gist file
 - Support for other application types, python, node, static sites. 
 - System for supporting other applications, some kind of flexible format so that new application types can be easily added/modified. 

## Server Types

**Router:**
Sits on the wildcard subdomain `*.gitbao.com`. Routes all apps to their correct location.

**Kitchen**
Sits on `gitbao.com` & `gist.gitbao.com`. Handles creation of Bao's and serving of site pages. 

**XiaoLong**
Contain's Bao's. Can be triggered to create a new Bao, or provide logs/information for existing Bao's. 

## Spec

 - Github gist is https://gist.github.com/maxmcd/1bd61d7c9c68afb7fdda
 - Change url to https://gist.gitbao.com/maxmcd/1bd61d7c9c68afb7fdda
 - Url is parsed by server, ID is grabbed from gist url
 - Api request is made for gist metadata
     + `http://api.github.com/gists/1bd61d7c9c68afb7fdda`
     + Grab `files`, `git_pull_url`
     + Fallback is `https://gist.github.com/maxmcd/1bd61d7c9c68afb7fdda.json`
 - Redirect to config URL for application. `https://gitbao.com/bao/id`
     + This page allows for setting of ENV variables and other configuration. 
     + "deploy" button
     + Expires in 24hrs
     + Deplying triggers build on build server, logs streamed in with errors
     + Option to update and redeploy from same gist
 - Build server in docker, build logs stream right to the browser with long polling. 
 - Server at subdomain
 - Server is up for 24 hrs, more if you want it. Can log in to set up auto-updates, env variables, and other options. 
 - https://developer.github.com/v3/activity/events/types/#gistevent
 - https://developer.github.com/v3/auth/#via-oauth-tokens