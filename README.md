
# Worgenda

If you have some [org-mode](http://orgmode.org/) files shared on Dropbox or other
similar service, you must understand my frustration when you try to view or edit
your diary scheduling in differents devices (smartphone, tablet, office
computer, etc.) without Emacs. There are differents solutions well-known in the community as
org-mobile or orgzly but both are mobile oriented and its usability is a little
"tricky" 

Worgenda is a pet project just for improve my Go skills. It just
covers a little subset of the orgmode features but by now is enough for me:

* [x] Show a calendar with the events notes.
* [x] Show the orgfiles in raw mode
* [ ] Add new notes to the orgmode files in the remote service
* [ ] Show TODOs task list


## Dependencies

If you want to compile Worgenda you need install the
[Dropbox client for Go](https://github.com/stacktic/dropbox)

```
go get github.com/stacktic/dropbox
```

## Configuration

Before run Worgenda you have to create a configuration directory named `./var`
with some important files:

* `user.json`: This file is the database of users (one user in fact). For
  example:
  
```
	[
		{
			"Username":"sdemingo",
			"Fullname":"Sergio",
			"Password":"md5sum-password"
		}
	]
```

* `config.json`: This file stores the tokens to authentication in Dropbox
  service. To get them you must register your app in the
  [Dropbox Apps Panel](https://www.dropbox.com/developers/apps). The `Files`
  property has the relative paths from Dropbox root of your org-mode files.

```
{
"AppKey":"<appkey>",
"AppSecret":"<appsecret>",
"Token":"<apptoken>",
"Files":[
    "org/work.org",
    "org/home.org",
    "org/others.org"]
}
```

* `private_key` and `public_key`: To complete a TLS session with each client.


## Run

When your `./var` directory is ready just run worgenda with:

```
$ worgenda -d <domain-or-ip>
```

Flag `-d` is the domain name or the ip address of worgenda machine. By default
'localhost' is used.
