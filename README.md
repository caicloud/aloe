# Aloe 

[![GoDoc](http://godoc.org/github.com/caicloud/aloe?status.svg)](http://godoc.org/github.com/caicloud/aloe)
[![Go Report Card](https://goreportcard.com/badge/github.com/caicloud/aloe)](https://goreportcard.com/report/github.com/caicloud/aloe)

Aloe is a declarative API test framework based on [ginkgo](https://github.com/onsi/ginkgo) and [gomega](https://github.com/onsi/gomega).
It aims to write and read API test cases simply.

Now only json API is supported.

## Get Started

Aloe assumes that API test is consist of `context` and `case`. 

If a GET API should be tested, the resource will be created in `context` and getted in `case`.

For example, a dir should be created to put `context` and `case`.

```
mkdir test/testdata
```

First, define `context` in `test/testdata/_context.yaml`.

`context` is defined in `_context.yaml`. Each dir means a context and context can be nested just like dir.
All test cases in the same dir will always run in same context.

Normally context is used to init data in database so that all test cases will be tested in determined context.

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

Then, define `case` in `test/testdata/get.yaml`.

`case` defines a test case which is defined in a yaml file. 

Normally test case should only be some simplest http requests.

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

Finally, some go codes should be writed in `test` dir to run the test case

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

## Variable

Sometimes variables should be defined because of auto-generated value by server(e.g. id).

```yaml
flow:
- description: "Init a product"
  definitions:
  - name: "testProductId"
    selector:
    - "id"
```

A variable called `testProductId` will be defined. The value of `testProductId` is from the roundtrip response.

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

If a varibale is defined, it can be used in roundtrip with format %{name}.

## Body validator

Sometimes only format of an auto-generated value need be checked. 

Some special validators are predefined, e.g. `$regexp`

```yaml
flow:
- description: "Create a product"
  response:
    # validate that id format should match regexp.
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

Now only `$regexp` and `$exists` is supported (more sp validator will be added).

## Examples

we can get some examples in:

- [get-started](./example/get-started)

- [crud](./example/crud)

