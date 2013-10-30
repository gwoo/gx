gx
=======

Remote execution via http. Scripts are base64 encoded and passed as the first element of the path. Obviously, be very careful. With great flexibility, you have to be Spiderman.

The main reason was for git deploys via webhook. Go made it possible to create a small (3.8MB) binary without any dependencies. Now, point a github repo at the host where you installed gx with a custom base64 encoded script

### Install Options
 - Source: `go get github.com/gwoo/gx`
 - Ubuntu: `bash < <(curl gist.github.com/gwoo/57c93572988116e47a3a/raw/install.sh)`
 - Download: https://github.com/gwoo/gx/releases

### Help
	./gx -help
	Usage of ./gx:
	  -host="localhost": Host for the server.
	  -password="test": Password for basic auth.
	  -port=8080: Port for the server.
	  -username="demo": Username for basic auth.

### Start the Server
	./gx

### Adding SSL
Files are served through SSL, if a `cert.pem` and `key.pem` exist in the current working directory where `gx` is started.

To generate, use openssl

	openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem

### A simple echo using bash
http://localhost:8080/IyEvdXNyL2Jpbi9lbnYgYmFzaAoKZWNobyAiaGVsbG8i

### A simple echo using ruby
http://localhost:8080/IyEvdXNyL2Jpbi9lbnYgcnVieQoKcHJpbnQgImhlbGxvIHdvcmxkIgo=

#### License
The BSD License http://opensource.org/licenses/bsd-license.php.

#### Todo
 - Write some godocs.
 - Add SSL support.
 - Add more options for hook handling.
 - Maybe add a config file.
