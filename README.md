# Gourlex

## Overview

Gourlex is a simple tool that can be used to extract URLs and paths from web pages. 
It can be helpful during web application assessments to uncover additional targets.

![gourlex](https://github.com/trap-bytes/gourlex/blob/main/static/gourlex.png)

## Features

- **URLs and Paths Extraction**
  - The tool can be used to extract only URLs, only paths, or both.
- **Silent mode for easy integration with other tools**
  - The tool provides a silent mode, making it easy to integrate its output into other tools during the reconnaissance and enumeration phases.

## Install

```
go install github.com/trap-bytes/gourlex@latest
```
## Usage:

```
gourlex -h
```

This will display help for the tool. Here are all the arguments it supports.

```
Usage:
  gourlex [arguments]

The arguments are:
  -t string    Specify the target URL (e.g., domain.com or https://domain.com)  
  -p string    Specify the proxy URL (e.g., 127.0.0.1:8080)
  -c string    Specify cookies (e.g., user_token=g3p21ip21h; 
  -r string    Specify headers (e.g., Myheader: test
  -s           Silent Mode, avoid printing banner and other messages
  -uO          Extract only full URLs
  -pO          Extract only URL paths
  -h           Display help

Example:
  gourlex -t domain.com
```
