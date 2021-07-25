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

## Code Examples

### Basic HTML format usage  

```bash
➜ pretgo phtml wrap=80 <old.html >new.html
```

### Basic Json format usage

```bash
➜ cat mes.json | pretgo pjson
```

### Basic XML format usage

```bash
➜ cat mes.xml | pretgo pxml
```

## Status

Project is: _in progress_ .

## Contact

Created by [@muratmirgun](https://twitter.com/muratmirgun) - feel free to contact me!
