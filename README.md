# Snowid
Snowid is a service which generates unique sortable time-based 64 bit long identifiers.

# Format

```
[ 00000000000000000000000000000000000000000    0    0000000000 000000000000 ]
  |--------------(1)----------------------| |-(2)-| |---(3)--| |---(4)----|
```
1) 41 bits for Unix Timestamp with milliseconds
2) Reserved bit
3) 10 bits for Machine ID
4) 12 bits for incrementing number

| field      | max value (dec) | max value (hex) |
|------------|-----------------|-----------------|
| Timestamp  | 2199023255551   | 0x1FFFFFFFFFF   |
| Machine ID | 1023            | 0x3FF           |
| Number     | 4095            | 0xFFF           |

The format is based on idea that each service in the cluster has its unique machine ID. 
It makes some points for running service. For example, if you run snowid in Kubernetes cluster, you should use StatefulSet instead of a regular Deployment. 

The service also supports custom epoch to use it in the timer. Be careful with changing epoch setting for the existing services as it may bring collisions. 
usually, you don't need to change epoch during the lifetime.   

# Installing

## With docker

1. Clone the repository
2. Run `docker build . -t snowid:latest`

## From sources

```
go get github.com/azzz/snowid/cmd/snowid
```

# Running

## Configuration
The service is configured by setting environment variables.

| Environment Variable | Optional | Description                                                | Default | Example        |
|----------------------|----------|------------------------------------------------------------|---------|----------------|
| LOG_LEVEL            | optional | Log level                                                  | info    | info           |
| LISTEN               | required | Host and port the service runs on                          |         | ":8080"        |
| EPOCH                | required | Epoch the timer start counting time from                   |         | 20210413001805 |
| MACHINE_ID           | required | Machine/service id. Must be unique in the scope of cluster |         | 42             |

## Example

```
LOG_LEVEL=info LISTEN=":8080" EPOCH="20210413001805" MACHINE_ID=333 snowid`
```

# Usage

## HTTP

The service can be called by HTTP

### GET /id64

**description:** Return an ID in both numeric and string formats. The string format might be useful if the language does not support 64 bit numbers. For example, JavaScript expects the numbers are 52 bits long. 

**content-type:** application/json

**Example:**

```
{"numeric":41859390008987648,"string":"41859390008987648"}
```

## Examples

```
> curl -i http://localhost:8080/id64
HTTP/1.1 200 OK
Date: Wed, 09 Jun 2021 18:25:13 GMT
Content-Length: 58
Content-Type: text/plain; charset=utf-8

{"numeric":41859390008987648,"string":"41859390008987648"}
```
