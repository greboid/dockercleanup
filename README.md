### Docker Cleanup

CLI utility to cleanup docker images and containers in the background.

## Docker usage

Available as a docker image, will access either CLI arguments or environmental variables for configuration.

```
version: '3.7'

services:
  dockercleanup:
    image: greboid/dockercleanup
    environment:
      DURATION: 24h
    privileged: true
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    restart: always
```

## Basic CLI Usage

This can also be installed and run directly:

```
go install github.com/greboid/dockercleanup
```
    
```
  dockercleanup --duration [repeat every X duration]
```

## Config Options

You can optionally pass the follow options


  ```--allimages=true``` to cleanup all images, rather than just dangling images
  
  
  ```--containers=true``` to cleanup non running containers

