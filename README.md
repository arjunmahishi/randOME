# RandOME

A CLI tool to generate a random stream of open metrics data for debugging prometheus based time series databases

---

### Easiest way to run

**Emit default metrics**

```bash
$ docker run -p 9090:9090 arjunmahishi/randome
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

$ docker run -p 9090:9090 -v $(pwd)/sample.yaml:/app/sample.yaml  arjunmahishi/randome
```

### Building locally

```
# for your OS
make my_build

# for multiple platforms
make build
```

### How to run the binary

```
 ➜  randOME --help
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

#### Cheat sheet

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
 ➜  cat sample.yaml
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

 ➜  randOME print -c sample.yaml -f 5s
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
