# GreyELK

Welcome to GreyELK. GreyELK imports Ethereum blocks into Elasticsearch.

Ensure that this folder is at the following location:
`${GOPATH}/src/github.com/manifoldfinance/greyelk`

## Getting Started with GreyELK

### Requirements

- [Golang](https://golang.org/dl/) 1.7

### Init Project

To get running with GreyELK and also install the
dependencies, run the following command:

```
make setup
```

It will create a clean git history for each major step. Note that you can always rewrite the history if you wish before pushing your changes.

To push GreyELK in the git repository, run the following commands:

```
git remote set-url origin https://github.com/manifoldfinance/greyelk
git push origin master
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

### Build

To build the binary for GreyELK run the command below. This will generate a binary
in the same directory with the name greyelk.

```
make
```

#### Using Docker

To build greyelk Docker image

```
docker build -t <docker-repo>/greyelk:latest .
```

### Run

#### Config

Before starting, there are configurations user need to set

- eth_rpc_addr The RPC endpoint of the blockchain.
  - It can be RPC port on an Ethereum fullnode, for example http://1.2.3.4:8545
  - If you use Infra, it can also be your Infra project, for example https://<network\>.infura.io/v3/YOUR-PROJECT-ID
- start_block The starting block to be imported. It must be below the current blockchain height.
  If the value is set to negative number, the import will start from the current blockchain height.

To run GreyELK with debugging output enabled, run:

```
./greyelk -c greyelk.yml -e -d "*"
```

#### Run with Docker image

```
 docker run -e ELASTIC_HOST="<Elasticsearch_host>:9200" -e ETH_RPC_ADDR="<Ethereum RPC endpoint>" --rm -d <docker-repo>/greyelk
```

### Test

To test GreyELK, run the following command:

```
make testsuite
```

alternatively:

```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `fields.yml` by running the following command.

```
make update
```

### Cleanup

To clean GreyELK source code, run the following command:

```
make fmt
```

To clean up the build directory and generated artifacts, run:

```
make clean
```

### Clone

To clone GreyELK from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/src/github.com/manifoldfinance/greyelk
git clone https://github.com/manifoldfinance/greyelk ${GOPATH}/src/github.com/manifoldfinance/greyelk
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
make release
```

This will fetch and create all images required for the build process. The whole process to finish can take several minutes.
