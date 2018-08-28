# depgoparser
Parses go project folder to produce a model in json format, which can be read by [Dependinator](https://github.com/michael-reichenauer/Dependinator), which will visualize the model in a map like interface.

## Installing

```bash
> go get -u github.com/michael-reichenauer/depgoparser
```
## Using
```bash
> depgoparser [-o path/to/model.json] path/to/go-project/folder
```