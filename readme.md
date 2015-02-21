# Gitbao

Wrapping up your github gists into delicious cloud-shaped buns. 

## Spec

 - Github gist is https://gist.github.com/maxmcd/1bd61d7c9c68afb7fdda
 - Change url to https://gist.gitbao.com/maxmcd/1bd61d7c9c68afb7fdda
 - Url is parsed by server, ID is grabbed from gist url
     + Make sure there aren't other formats for gist urls
     + Expect all kinds of noise at the end of the url
     + https implementation is clearly not optional
 - Api request is made for gist metadata
     + `http://api.github.com/gists/1bd61d7c9c68afb7fdda`
     + Grab `files`, `git_pull_url`
     + Fallback is `https://gist.github.com/maxmcd/1bd61d7c9c68afb7fdda.json`
 - Redirect to config URL for application. `https://gitbao.com/subdomain`
     + This page allows for setting of ENV variables and other configuration. 
     + "deploy" button
     + Expires in 24hrs
     + Deplying triggers build on build server, logs streamed in with errors
     + Option to update and redeploy from same gist
 - Build server in docker, build logs stream right to the browser with long polling. 
     + Could get files from githubusercontent, although not sure about rate limits. 
     + Could download file from .git endpoint as well, likely the best way
     + Need to build docker file with required imports, installs and other options.
     + `FROM golang:onbuild` in a dockerfile will build the app and all dependencies in that directory.
     + Alternatively `go get ./...` will get and install all child dependencies.
- Server at subdomain
    + There should be one users database
    + Servers can be very modular, just need the entrypoint servers to talk to a DB. 
    + Routing is interesting. Need to keep track of various servers with specific endpoints. 
    + Could write new ips and ports to nginx file and route through one load balancer.
    + Maybe see if AWS has resources for stuff like this, writing dynamic routes to ELB would be cool
 - Personal git url? Nah, just update the gist. 
 - Server is up for 24 hrs, more if you want it. Can log in to set up auto-updates, env variables, and other options. 
 - https://developer.github.com/v3/activity/events/types/#gistevent
 - https://developer.github.com/v3/auth/#via-oauth-tokens