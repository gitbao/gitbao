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

## Mac: Quickstart

Make sure that your PATH includes the PostGres command line tools. For the PostGres command line tools to work, you'll need to have PGDATA assigned. Then start PostGres.

```
export PATH=$PATH:'/Applications/Postgres.app/Contents/Versions/9.4/bin'
export PGDATA='/Users/Andrew/Library/Application Support/Postgres/var-9.4'
pg_ctl start
```

Mac can't run Docker, so boot2docker runs a linux VM in the background that docker connects to.

```
boot2docker init
boot2docker start
```

Assign GOPATH to something reasonable, then download gitbao and it's dependencies.

```
# Optional
export GOPATH="/usr/local/lib"
mkdir -p $GOPATH

# Not Optional
go get github.com/gitbao/gitbao
cd $GOPATH/src/github.com/gitbao/gitbao
go get ./...
```

In one terminal window start Kitchen, which is the main webapp.

```
# Optionally use tmux because it's awesome.
tmux new-session -s kitchen

# Load the boot2docker env and start kitchen.
$(boot2docker shellinit)
cd $GOPATH/src/github.com/gitbao/gitbao/cmd
go run kitchen/main.go
```

In another terminal window start Xialong, which creates baos.

```
# Optionally use tmux because it's awesome.
tmux new-session -s xiaolong

# Load the boot2docker env and start xiaolong.
$(boot2docker shellinit)
cd $GOPATH/src/github.com/gitbao/gitbao/cmd
go run xiaolong/main.go
```

### Troubleshooting

Make sure that no other docker instances exist when starting gitbao.

```
docker rm -f $(docker ps -a -q)
```