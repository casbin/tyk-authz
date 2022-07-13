# tyk-authz

Casbin authorization plugin for Tyk.


## Installation

Compile it:
``` bash
docker run --rm -v `pwd`:/plugin-source tykio/tyk-plugin-compiler:v3.2.3 plugin.so
```


Copy the example Tyk Gateway api definition `examples\casbin-authz-api-example.json` (as well as the model&policy configfile) and the compiled `plugin.so` to the Tyk plugins directory. 

Modify the path in the `casbin-authz-api-example.json` to point to the compiled `plugin.so` file & config files. 

Start the Tyk Gateway and test the plugin.

```bash
curl http://localhost:8080/test/get -i  # it will receive a 401
curl -H "username:alice" http://localhost:8080/test/get -i -v  # it will receive a 200
```
