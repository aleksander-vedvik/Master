# Master

### Preliminary

The entirety of Gorums is added to this repository. The main contribution to Gorums in this Master's thesis are outlined in the readme file inside the broadcast folder of Gorums: `./gorums/broadcast/README.md`

The root folder is defined as the folder this file is in.

## Benchmarks

Three implementations of Paxos and two implementations of PBFT have been benchmarked:

1. Paxos using BroadcastCall
2. Paxos using QuorumCall
3. Paxos using QuorumCall with BroadcastOption
4. PBFT using Gorums
5. PBFT without using Gorums

All benchmarks are configured to run in docker containers. A single client container will be created which is responsible of collecting results and saving them into csv files. These files will be saved in the folder: `./csv`

### Prerequisites

The experiments have been run on Ubuntu and in WSL2 on Windows. The prerequisites to run the benchmarks are as follows:

- [Make](https://www.gnu.org/software/make/manual/make.html)
- [Docker](https://www.docker.com/)

### Folder structure

The folder structure is as follows:

    ├── csv
    ├── example
    │   ├── reliablebroadcast
    ├── gorums
    ├── logs
    ├── src
    │   ├── benchmark
    │   ├── leaderelection
    │   ├── paxos.bc    # Paxos using BroadcastCall
    │   ├── paxosqc     # Paxos using QuorumCall
    │   ├── paxosqcb    # Paxos using QuorumCall with BroadcastOption
    │   ├── pbft.gorums # PBFT using Gorums
    │   ├── pbft.plain  # PBFT without using Gorums
    │   ├── util

The `example` folder contains the example given in the thesis. The `csv` and `logs` folder contains the result files from running benchmarks. The `gorums` folder contains all the Gorums code including the broadcast framework developed in this thesis. Each implementation of Paxos and PBFT is denoted by a comment in the folder structure above. The `src/benchmark` folder contains the code for running the benchmarks, while `src/util` is used to process the results from the benchmarks.

### Options

Eight docker containers will be started when running the benchmarks. One of these is the client docker container while the rest are running a single server instance. This enables running the benchmarks in different server configurations:

- Paxos can be run with either 3, 5, or 7 servers.
- PBFT is run with 4 servers.

Redundant servers (i.e. when using a server configuration with less than 7 servers) will exit immediately after creation to not interfere with the benchmark.

The benchmarks are controlled entirely by using the `.env` file in the root directory. Each option will be presented in the following sections.

#### Benchmark type

To specify which benchmark to run we can specify the `BENCH` option in the `.env` file. This accepts a number between 1-5 and the numbers correspond to:

1. PaxosBC
2. PaxosQC
3. PaxosQCB
4. PBFT.With.Gorums
5. PBFT.Without.Gorums

For example, to benchmark PBFT using Gorums we have to specify the following in the `.env` file:

```yml
BENCH=4
```

#### Metric type

Two types of metrics have been measured: Performance and Throughput vs Latency. The performance benchmark measures the performance of the implentations under common operations. Additionally, it will append "Performance" to the csv files when creating the results. The throughput vs latency benchmark sends client requests at the target throughput and measures the round trip latency of the requests.

This metric type option is compatible with all benchmark types, and it accepts a number:

0. Throughput vs Latency
1. Performance

For example, to run the Performance benchmark we have to specify the following in the `.env` file:

```yml
TYPE=1
```

#### Server configuration

There is a set of configuration files in the root folder: `conf.*.yml`. These files are used to configure the servers. The different server configurations are as follows:

- 3 servers (Paxos): `paxos.3`
- 5 servers (Paxos): `paxos.5`
- 7 servers (Paxos): `paxos.7`
- 4 servers (PBFT): `pbft`

For example, to run with 5 servers we have to specify the following in the `.env` file:

```yml
CONF=paxos.5
```

#### Throughput vs Latency configuration

It is possible to configure the Throughput vs Latency benchmark.

To specify the max throughput (client requests/second) the benchmark will run at, it is possible to specify the following field:

```yml
THROUGHPUT=5000
```

By using the `DUR` field, it is possible to define how long the benchmark should run. It accepts a number which will make the benchmark run at the target throughput for the given number of seconds. E.g. if `THROUGHPUT=5000` the benchmark can be configured to run at this throughput for 3 seconds by defining the following in the `.env` file:

```yml
DUR=3
```

The `STEPS` field divides the benchmark into the given number of steps. It will run the benchmark the given number of times with incremental steps until the max throughput. E.g. if `STEPS=2` and `THROUGHPUT=5000`, then it will first run the benchmark with `target throughput = 2500`. Then it will purge all state after each step in order to not make the benchmarks not interfere with each other. Afterwards, it will run the same benchmark with `target throughput = 5000`. To run a benchmark in 10 steps, you can specify the following in the `.env` file:

```yml
STEPS=10
```

Lastly, to run a benchmark several times it is possible to use the `RUNS` field. This will run the benchmark with identical setup for the given number of times. If `STEPS=5` and `RUNS=2`, then the benchmark will be run in total 10 times. To run a benchmark setup once, you can specify the following in the `.env` file:

```yml
RUNS=1
```

NOTE: These options will be ignored when running the Performance benchmark.

#### Defaults

There are two options left that can be modified. The first one is `PRODUCTION`. This should always be set to `1` when running in docker containers. This option can be set to `0` if running each server as a separate process. This was only used for test purposes and should therefore not be changed.

```yml
PRODUCTION=1
```

It is also possible to enable logging. This will incur overhead when running the benchmarks. Hence, logging is disabled by default but it can be enabled by setting:

```yml
LOG=1
```

Log files will be stored in the folder: `./logs`

### How to run

When running the benchmarks for the first time, we need to create the docker network and build the docker containers. This will also start the benchmark corresponding to what is specified in the `.env` file:

    make network
    make docker

Servers are not configured to stop automatically after a benchmark has been run. The client docker container will however stop after it is finished. Hence, it is necessary to stop the server containers after the client container has been stopped (docker will print `client-1 exited with code 0` when the client container has been stopped) by pressing:

    Ctrl + c

After the docker containers have been built the first time, it is sufficient to run:

    make eval

NOTE: There is no need to rebuild the docker containers after changing the `.env` file. Hence, to run a different benchmark you can change the `.env` file accordingly and then run `make eval`. This will load and start the correct benchmark based on the `.env` file. If the changes you make in the `.env` file is not reflected when running `make eval`, then run `make docker` to rebuild the containers.
