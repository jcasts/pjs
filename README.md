# pjs
Pjs is an extremely simple command line tool to pretty print and search JSON data.

```
$ pjs ~/path/to/file.json
$ curl api.twitter.com/1.1/notifications.json | pjs
```

### What does it do?
* Reads JSON from a STDIN stream or a file.
* Supports continuous streaming of multiple JSON values.
* Formats and indents your JSON for readability. (Yes, indentation is configurable)
* Colorizes by data type. (Yes, you can turn it off)
* Sorts JSON map keys alphabetically, for readability.
* Can filter JSON data by specific keys or values.

### Installation
Pjs is written in Go, which means it's highly portable. You can get binaries here:

* [darwin/amd64](https://github.com/jcasts/pjs/blob/master/bin/darwin_amd64/pjs.zip?raw=true)
* [linux/amd64](https://github.com/jcasts/pjs/blob/master/bin/linux_amd64/pjs.zip?raw=true)

If you prefer compiling from source, there are no dependencies besides Go:

```
$ go get github.com/jcasts/pjs
```

### Filtering JSON
Filtering is done by passing glob-like paths to pjs after the double dash -- delimiter.

```
$ pjs test.json -- object/name
{
  "object": {
     "name": "Jim"
  }
}
```

Multiple paths can be specified, although the examples will stick to one for simplicity.

#### Key Matching
Keys may be complete or partial matches. A partial match is achieved with the wildcard * character:

```
$ pjs test.json -- object/n*
{
  "object": {
     "name": "Jim",
     "new_user": false
  }
}
```

An or operator may also be used to specify multiple exact keys, represented by the pipe character | :

```
$ pjs test.json -- "object/name|new_user"
{
  "object": {
     "name": "Jim",
     "new_user": false
  }
}
```

#### Arrays
If you have an array of objects, you can specify a wildcard * for the index in the path:

```
$ pjs test.json -- objects/*/name
{
  "objects": [
    {
       "name": "Jim"
    },
    {
      "name": "Amy"
    },
    {
      "name": "Lea"
    },
    {
      "name": "Alison"
    }
}
```

An index or index range may also be specified to select parts of an array:

```
$ pjs test.json -- objects/1..2/name
{
  "objects": [
    {
      "name": "Amy"
    },
    {
      "name": "Lea"
    }
}
```

#### Values
Similarly, JSON data may be filtered by value, by using the = character:

```
$ pjs test.json -- objects/*/name=Amy
{
  "objects": [
    {
      "name": "Amy"
    }
}
```

Wildcards and "or" selectors work the same way on values as they do on keys:

```
$ pjs test.json -- objects/*/name=A*
{
  "objects": [
    {
      "name": "Amy"
    },
    {
      "name": "Alison"
    }
}
```

```
$ pjs test.json -- "objects/*/name=Amy|Alison"
{
  "objects": [
    {
      "name": "Amy"
    },
    {
      "name": "Alison"
    }
}
```

#### Objects by Matching Attributes
Often you'll need to get the full object of a specific user. You can search for it by it's name or id, and then use the parent function, represented by double dots .., to select the parent.
```
$ pjs test.json -- objects/*/name=Amy/..
{
  "objects": [
    {
      "age": 43,
      "id": 84729478,
      "last_name": "Smith",
      "name": "Amy"
    }
}
```
 
#### Recursive Searches
For some data, you may not know how deep the information you need is buried. You can specify a recursive search to get the highest level matches in your data. Use the double wildcard ** to recursively search for matches:

```
$ pjs test.json -- **/name=Amy/..
{
  "objects": {
    "users": [
      {
        "age": 43,
        "id": 84729478,
        "last_name": "Smith",
        "name": "Amy"
      }
    ]
  }
}
```

These can be chained to search through multiple specific structures.

Escape any special characters with backslash \

### License
MIT
