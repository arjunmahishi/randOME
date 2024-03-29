# RandOME

A CLI tool to generate a (pseudo)random stream of open metrics data for debugging prometheus based TSDBs and tools that depend on it. This is very similar to [Avalanche](https://github.com/prometheus-community/avalanche). While avalanche focuses more on the "load testing" aspect of things, randOME focuses more on the actual data generated. 

By controlling the shape of the data, you can test/debug/benchmark specific scenarios related to alerting, visualization, storage, and performance. With randOME, you can create more targeted and meaningful tests that enable you to optimize your systems more effectively.

---

## Configuration

This tool can be configured with a simple YAML file. Here's a sample config:

```yaml
metrics:
  - name: metric_1
    value_min: 0
    value_max: 100
    labels:
      instance: [localhost:8080]
      cluster: [dev, prod, staging, test, qa]

  - name: metric_2
    max_cardinality: 10
    labels:
      method: [GET, POST, PUT, DELETE]
      status: [200, 400, 404, 500]
```

**Attributes**

| Field         | Type                 | Description |
|---------------|----------------------|-------------|
| `Name`        | `string`             | The name of the metric. |
| `ValueMin`    | `number`                | The minimum value that the metric can have. |
| `ValueMax`    | `number`                | The maximum value that the metric can have. |
| `Value`       | `number`           | The static value of the metric. This field is optional. |
| `Labels`      | `map[string][]string`| The labels associated with the metric, as a map of string to string slices (refer to the above examples for the exact format). |
| `MaxCardinality`| `number`              | The maximum number of unique label value combinations that can be associated with the metric. This field is optional. |

## Easiest way to run

**Emit default metrics**

```bash
$ docker run -d -p 9090:9090 arjunmahishi/randome
```

**Emit custom metrics**

```bash
$ cat > sample.yaml <<EOF
metrics:
  - name: metric_1
    value_min: 0
    value_max: 100
    labels:
      instance: [localhost:8080]
      cluster: [dev, prod, staging, test, qa]

  - name: metric_2
    max_cardinality: 10
    labels:
      method: [GET, POST, PUT, DELETE]
      status: [200, 400, 404, 500]
EOF

$ docker run -d -p 9090:9090 -v $(pwd)/sample.yaml:/app/sample.yaml  arjunmahishi/randome
```

**Remote write default metrics**

```bash
$ docker run -d -e REMOTE_WRITE_ADDR='http://<prometheus-compatible-host>/api/v1/write' arjunmahishi/randome
```

**Remote write custom metrics**

```bash
$ cat > sample.yaml <<EOF
metrics:
  - name: metric_1
    value_min: 0
    value_max: 100
    labels:
      instance: [localhost:8080]
      cluster: [dev, prod, staging, test, qa]

  - name: metric_2
    max_cardinality: 10
    labels:
      method: [GET, POST, PUT, DELETE]
      status: [200, 400, 404, 500]
EOF

$ docker run -d -e REMOTE_WRITE_ADDR='http://<prometheus-compatible-host>/api/v1/write' -v $(pwd)/sample.yaml:/app/sample.yaml arjunmahishi/randome
```

## Building locally

```
# for your OS
make build

# for a specific OS (check the makefile to see if the command is available for the OS you're looking for)
make build_<platform>

# both commands will create an executable binary in ./bin
```

## Running the binary directly

```
$ randOME --help
usage: randOME [<flags>] <command> [<args> ...]

Flags:
      --help           Show context-sensitive help (also try --help-long and --help-man).
  -f, --frequency=4s   Frequency of data points
  -c, --config=CONFIG  path to the config file
      --version        Show application version.

Commands:
  help [<command>...]
    Show help.

  print
    print the data to stdout

  remote-write --addr=ADDR
    Remote write to a Prometheus-compatible TSDB

  emit [<flags>]
    Emit metrics over HTTP on a given port (localhost:<port>/metrics)
```

**Print metrics to STDOUT**

```
randOME print -f 5s
```

**Emit metrics over HTTP**

```
randOME emit -f 5s
```

**Remote write to a prometheus compatible storage**

```
randOME remote-write -f 5s --addr http://localhost:8428/api/v1/write
```

**Defining custom metrics**

```
$ cat sample.yaml
metrics:
  - name: cpu_usage
    value_min: 0
    value_max: 100
    labels:
      instance: [localhost:8080]
      cluster: [dev, prod, staging, test, qa]

  - name: requests_total
    max_cardinality: 10
    labels:
      method: [GET, POST, PUT, DELETE]
      status: [200, 400, 404, 500]

$ randOME print -c sample.yaml -f 5s
requests_total{status='200',method='GET'} 2
requests_total{status='400',method='GET'} 7
requests_total{status='404',method='GET'} 6
requests_total{status='500',method='GET'} 0
requests_total{status='200',method='POST'} 5
requests_total{status='400',method='POST'} 3
requests_total{method='POST',status='404'} 4
requests_total{method='POST',status='500'} 3
requests_total{status='200',method='PUT'} 9
requests_total{status='400',method='PUT'} 6
cpu_usage{instance='localhost:8080',cluster='dev'} 21
cpu_usage{instance='localhost:8080',cluster='prod'} 52
cpu_usage{instance='localhost:8080',cluster='staging'} 36
cpu_usage{instance='localhost:8080',cluster='test'} 5
cpu_usage{instance='localhost:8080',cluster='qa'} 62
cpu_usage{cluster='dev',instance='localhost:8080'} 94
cpu_usage{instance='localhost:8080',cluster='prod'} 93
cpu_usage{cluster='staging',instance='localhost:8080'} 80
cpu_usage{instance='localhost:8080',cluster='test'} 67
cpu_usage{cluster='qa',instance='localhost:8080'} 59
```

