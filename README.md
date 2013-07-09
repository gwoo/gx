GExec
=======

Remote execution via http. Scripts are base64 encoded and passed as the first element of the path. Obviously, be very careful. With great flexibility, you have to be Spiderman.

The main reason was for git deploys via webhook. Go made it possible to create a small (3.8MB) binary without any dependencies. Now, point a github repo at the host where you installed gexec with a custom base64 encoded script

	Usage of ./gexec:
	  -host="localhost": Host for the server.
	  -password="test": Password for basic auth.
	  -port=8080: Port for the server.
	  -username="demo": Username for basic auth.

### Start the Server
./gexec

### A simple echo using bash
http://localhost:8080/IyEvdXNyL2Jpbi9lbnYgYmFzaAoKZWNobyAiaGVsbG8i

### A simple echo using ruby
http://localhost:8080/IyEvdXNyL2Jpbi9lbnYgcnVieQoKcHJpbnQgImhlbGxvIHdvcmxkIgo=

#### License
The BSD License http://opensource.org/licenses/bsd-license.php.
