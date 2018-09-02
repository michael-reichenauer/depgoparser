# Dependinator Go Parser
Parses Go files within a folder to produce a model in json format, which can be read by [Dependinator](https://github.com/michael-reichenauer/Dependinator), which will visualize the model in a map-like interface.

## Installing

```bash
> go get -u github.com/michael-reichenauer/depgoparser
```
## Using
```bash
> depgoparser [-o path/to/model.json] path/to/go-project/folder
```
By default, the ouput model json file is written into the folder, but you can use the "-o" option to specify the output model path.
