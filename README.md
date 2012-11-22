# decisions

API & web app for deciding randomly among choices.

## Try It

You can try *decisions* at [decisionmaker.herokuapp.com](https://decisionmaker.herokuapp.com).


## Develop

A working [Go](http://golang.org) installation is required to compile *decisions*.  

```bash
$ cd $GOPATH
$ go get github.com/jmcvetta/decisions
```


## Deploy to Heroku

``Procfile`` and ``.godir`` files are included, making deployment to Heroku a breeze:

```bash
$ cd $GOPATH/src/github.com/jmcvetta/decisions
$ heroku create --buildpack git://github.com/kr/heroku-buildpack-go.git
$ git push heroku master
```


## License

*decisions* is Free Software, licensed under the [GPL
v3](http://www.gnu.org/copyleft/gpl.html).
