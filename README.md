# aloe

## get started

aloe use `ContextConfig` struct serialize your testdata and execute all RoundTrip in context.

init all test data in _context.yaml(create operate).
```yaml
summary: "CRUD Test example"
flow:
- description: "Create a product "
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
  - name: "testProduct"
  - name: "testProductId"
    selector:
    - "id"
```
define your variable in `definitions`. as above，we can use %{testProduct} define body and use %{testProductId} define product ID. then you can test `GET /products/%{testProductId}` api in your testcases.
```
description: "Try get a product"
flow:
- description: "Get product inited in _context.yaml"
  request:
    api: GET /products/%{testProductId}
    headers:
      "Content-Type": "application/json"
  response:
    statusCode: 200
    body: "%{testProduct}"
```
finally, run all testcases in aloe. aloe will do `cleanUp` function when every context execute finish. usual you need clean up your database to guarantee environment cleanly.
```go
func TestAPI(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	f := framework.NewFramework("localhost:8080", cleanUp,
		"testdata",
	)
	if err := f.Run(); err != nil {
		fmt.Printf("can't run framework: %v", err)
		os.Exit(1)
	}
	ginkgo.RunSpecs(t, "API Suite")
}

func cleanUp() {
	// clean up databases
}

var _ = ginkgo.BeforeSuite(func() {
	server.Product{}.Register()
	go http.ListenAndServe(":8080", nil)
})
```

## example

we can get some example in:

- [get started](./example/crud)

- [crud](./example/get_started)