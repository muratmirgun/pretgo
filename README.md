# PretGO

> So basic cli for format json,html and xml!

## Table of contents

* [Screenshots](#screenshots)
* [Setup](#setup)
* [Status](#status)
* [Contact](#contact)

## Screenshots

![Example screenshot](./static/pretgo.gif)

## Setup

### First clone project

```bash
git clone https://github.com/muratmirgun/pretgo
```

### And install it

```bash
go install
```
### Or use Dockerfile to generate image and container

```bash
# inside project
docker build -t pretgo-local .
```

### Then use it with json
```bash
➜ cat mes.json | docker run -i --rm pretgo-local -format="json"
```

### Or html

```bash
➜ cat mes.html | docker run -i --rm pretgo-local -format="html"
```

### Or xml format

```bash
➜ cat mes.xml | docker run -i --rm pretgo-local -format="xml"
```

## Code Examples

### Basic HTML format usage  

```bash
➜ cat mes.html | pretgo -format="html"
```

### Basic Json format usage

```bash
➜ cat mes.json | pretgo -format="json"
```

### Basic XML format usage

```bash
➜ cat mes.xml | pretgo -format="xml"
```

## Status

Project is: _in progress_ .

## Contact

Created by [@muratmirgun](https://twitter.com/muratmirgun) - feel free to contact me!
