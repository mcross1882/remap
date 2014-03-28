## remap

`remap` is a command line utility that allows you to reorder and filter CSV files. It works by using a 
newline delimited map file. In the map file each row contains a field that represents a header in the input
file you want to extract. The row position represents column index in the output file. So ordering is done vertically
in the map file. For more information check out the `samples/` folder.

### Running

```
Syntax

$ go run remap.go [input file] [output file] [map file]

Example

$ go run remap.go sample/input.csv sample/output.csv sample/custom_headers.map
```
