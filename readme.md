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

!!TODO: I will make available a Dockerfile with the steps. 

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

### Output

![image](https://user-images.githubusercontent.com/41498427/115451880-59e6cc00-a23b-11eb-97fd-d5aae93c5e55.png)

## Contribute & Bug Report
If you find a bug in web-serve or have any questions you can ask here. 

New Contributions are always welcome. Please discuss bugfix, performance improvements, new feature request before opening a PR.

