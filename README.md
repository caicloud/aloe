# Aloe

[![Build Status](https://travis-ci.org/caicloud/aloe.svg?branch=master)](https://travis-ci.org/caicloud/aloe)
[![Coverage Status](https://coveralls.io/repos/github/caicloud/aloe/badge.svg?branch=master)](https://coveralls.io/github/caicloud/aloe?branch=master)
[![GoDoc](http://godoc.org/github.com/caicloud/aloe?status.svg)](http://godoc.org/github.com/caicloud/aloe)
[![Go Report Card](https://goreportcard.com/badge/github.com/caicloud/aloe)](https://goreportcard.com/report/github.com/caicloud/aloe)

Aloe is a declarative API test framework based on [ginkgo](https://github.com/onsi/ginkgo), [gomega](https://github.com/onsi/gomega)
and [yaml](http://yaml.org/). It aims to help write and read API test cases simpler.

DISCLAIMER:
- only json API is supported now
- avoid using Aloe for extremely complex test, use ginkgo directly instead

## Terminology

There are two important concepts in aloe:

- `case`: case means a test case, for example, get an endpoint and verify its output.
- `context`: to simply put, a context is a group of cases. Usually, context is used to init data in database, so that all test cases will be tested in a determined context.

## Getting Started

Following is a getting started example for using Aloe. First, create a directory to put `context` and `case` yaml files.

```
mkdir -p test/testdata
```

Then, define your `context` in `_context.yaml`.

```yaml
# test/testdata/_context.yaml
summary: "CRUD Test example"
flow:
- description: "Init a product"
  request:
    api: POST /products
    headers:
      "Content-Type": "application/json"
    body: |
      {
        "id": "1",
        "title": "test"
      }
  response:
    statusCode: 201
  definitions:
  - name: "testProductId"
    selector:
    - "id"
```

As mentioned above, a context is used to run a group of test cases in a determined environment. In the above example,
we define a context which simply sends a POST request to `/products` with product name `test`; therefore, all test
cases in this context will expect product `test` exists.

Now with context setup, we can start defining `case`. Here, we define a test case in `test/testdata/get.yaml` to get
and verify product `test`.

```yaml
# test/testdata/get.yaml
description: "Try to GET a product"
flow:
- description: "Get the product with title test"
  request:
    api: GET /products/%{testProductId}
    headers:
      "Content-Type": "application/json"
  response:
    statusCode: 200
```

Finally, some go codes should be written in `test` directory to run the test case:

```go
func TestAPI(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	f := aloe.NewFramework("localhost:8080", cleanUp,
		"testdata",
	)
	if err := f.Run(); err != nil {
		fmt.Printf("can't run framework: %v", err)
		os.Exit(1)
	}
	ginkgo.RunSpecs(t, "API Suite")
}

func cleanUp() {
	// function to clean up context
	// normally it is used to drop database
}

var _ = ginkgo.BeforeSuite(func() {
	s := server.NewServer()
	go http.ListenAndServe(":8080", s)
})
```

Pay attention to the `cleanUp` method, which is used to clean up a context.

## Usage

### Variable

Variables can be defined to hold auto-generated value by server (e.g. id). For example:

```yaml
flow:
- description: "Init a product"
  definitions:
  - name: "testProductId"
    selector:
    - "id"
```

A variable called `testProductId` will be defined. The value of `testProductId` is from the round trip response.

Variable can also be a json snippet.

```
response:
{
	"id": "111",
	"title": "test",
	"comments":[
		"aaa",
		"bbb"
	]
}

// id will select value of string
["id"] => 111(without quote)

// empty selector will select whole body
[] => {
	"id": "111",
	"title": "test",
	"comments":[
		"aaa",
		"bbb"
	]
}

// comments will select partial body
["comments"] => [
	"aaa",
	"bbb"
]
```

If a variables is defined, it can be used in round trip with format `%{name}`.

### Body validator

Body validator is used to validate response fields. Some special validators are predefined, e.g. `$regexp`

```yaml
flow:
- description: "Create a product"
  response:
    # validate that id format matches regexp
    # validate that password field is not returned
    body: |
      {
        "id": {
          "$regexp": "[a-zA-Z][a-zA-Z0-9-]{11}"
        },
        "password": {
          "$exists": false,
        },
      }
```

Now only `$regexp` and `$exists` is supported (more special validator will be added in the future).

### Nested context

Context can be nested just like directory. Child context will see all setup in parent context.

```
tests
└── testdata
    ├── _context.yaml
    ├── basic
    │   ├── _context.yaml
    │   ├── create.yaml
    │   └── update.yaml
    ├── failure
    │   ├── _context.yaml
    │   └── create.yaml
    └── list
        ├── _context.yaml
        └── list_all.yaml
```

## Examples

For more examples, see:

- [crud](./examples/crud)
