# Go Import Cycle
[![Go Report Card](https://goreportcard.com/badge/github.com/samlitowitz/godepvis)](https://goreportcard.com/report/github.com/samlitowitz/godepvis)
[![Pipes](https://github.com/samlitowitz/godepvis/actions/workflows/pipes.yaml/badge.svg?branch=master)](https://github.com/samlitowitz/godepvis/actions/workflows/pipes.yaml)

`godepvis` is a tool to visualize Go imports resolved to the package or file level.

# Installation
`go install github.com/samlitowitz/godepvis/cmd/godepvis@v1.0.6`

# Usage
```shell
godepvis --path examples/simple/ --dot imports.dot
dot -Tpng -o assets/example.png imports.dot
```

![Example import graph resolved to the file level](assets/examples/direct-circular-dependency/file.png?raw=true "Example import graph resolved to the file level")

Red lines indicate files causing import cycles between packages. Packages involved in a cycle have their backgrounds colored red.

```shell
godepvis --path examples/simple/ --dot imports.dot --resolution package
dot -Tpng -o assets/example.png imports.dot
```
![Example import graph resolved to the package level](assets/examples/direct-circular-dependency/package.png?raw=true "Example import graph resolved to the package level")

Red lines indicate import cycles between packages.

## Configuration
The palette file follows the JSON Schema outlined in [assets/palette-schema](assets/palette-schema).

The [simple-palette example](examples/simple-palette) uses the following schema...

```yaml
base:
  packageName: "rgb(174, 209, 230)"
  packageBackground: "rgb(207, 232, 239)"
  fileName: "rgb(160, 196, 226)"
  fileBackground: "rgb(198, 219, 240)"
  importArrow: "rgb(133, 199, 222)"
cycle:
  packageName: "#FFB3C6"
  packageBackground: "#FFE5EC"
  fileName: "#FF8FAB"
  fileBackground: "#FFC2D1"
  importArrow: "#FB6F92"
```

...to produce the following outputs...

![Example import graph resolved to the file level](assets/examples/simple-palette/file.png?raw=true "Example import graph resolved to the file level")

![Example import graph resolved to the package level](assets/examples/simple-palette/package.png?raw=true "Example import graph resolved to the package level")
