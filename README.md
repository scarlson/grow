# grow #

Golang reddit OAuth Wrapper

[reddit api]: http://www.reddit.com/dev/api

### install ###

```
go get github.com/scarlson/grow
```

### getting started ###

Visit https://ssl.reddit.com/prefs/apps/ to setup a new app.
Set redirect url to yourhost/login.

Rename example/example_config.json to example/config.json.
Edit config.json to fill in app id, app secret, and user agent.


```
cd example
go run main
```

Will get your server up and running.

### disclaimer ###

This library is still an early work in progress.  Expect many things to change/break in the near future.
