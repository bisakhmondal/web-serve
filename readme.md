# Web-Serve
![image](https://user-images.githubusercontent.com/41498427/115450569-b0530b00-a239-11eb-9bc8-6b7cca810e1a.png)

## Procedure
- Generate your production build through `npm run build` (for eg.) or any appropiate way.
- Place the build directory inside `html` directory replacing the existing one.
-  Update the configuration in `conf.yml` if required. The YAML is self explanatory.
-  Build the binary.
 
 ```shell
 # Required go version 1.16 or higher
 export GO111MODULE=on
 go build -o <output binary name> main.go
  ```
  
- Run the binary by `./<binary_name>`

### Build the Binary using Docker
Now no hassle of installing dependencies to build the binary. FYI, I can't provide ready to use built binary for diferrent platforms as the web build assets is going to be embedded into the binary during compile time.

But be assured, i have provided a very easy way to resolve the issue. Just run these commands,

```shell
# build the docker image
docker build -t web-serve . 

# Running the build
docker run -v $(pwd):/web-serve -e OS="linux" -e ARCH="amd64" -e BIN_NAME="your binary name" web-serve

# Change ownership of the binary directory from root to host user
sudo chown -R $(id -u) binary
```

As binaries are platform dependent, you have to mention the Operating system and processor instruction Architecture via environment variables `OS` and `ARCH` respectively. A binary will be created at `binary/` dorectory of current working directory. The environment variable `BIN_NAME` (Optional) is the name of the output binary.

**A list of supported OS and Architecture has been made available [here](./dist.json). Just replace it with the environment variable `OS` and `ARCH` during docker run.**

## Supported Commands
`./<binary_name> help`
<pre>
web-serve - A lightweight production build webapp server

Usage:
  web-serve [commands] [flags]
  web-serve [command]

Available Commands:
  help        Help about any command
  install     install server as a os service
  remove      remove/uninstall the service related to the server
  start       start the server
  status      server current status running/stopped
  stop        stop the server

Flags:
  -h, --help   help for web-serve

Use "web-serve [command] --help" for more information about a command.
</pre>

## Output

![image](https://user-images.githubusercontent.com/41498427/115451880-59e6cc00-a23b-11eb-97fd-d5aae93c5e55.png)

## Contribute & Bug Report
If you find a bug in web-serve or have any questions you can ask here. 

New Contributions are always welcome. Please discuss bugfix, performance improvements, new feature request before opening a PR.
