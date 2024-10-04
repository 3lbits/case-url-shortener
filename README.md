# Case assignment: Finish the URL shortener service

This is a case assignment for a Junior Platform Engineer position at ElBits AS. The task is to implement the last piece of a simple URL redirection service, and to deploy it.

**TODO**

- [ ] Implement [`LinkFile.Load`](./linkfile.go)
- [ ] Deploy the application to [fly.io](https://fly.io).

## Prerequisites

1. Ensure the `git` command works.
2. Clone this repository: `git clone https://github.com/3lbits/case-url-shortener.git`
3. [Install Go](https://go.dev/learn/).
4. Ensure the `go` and `gofmt` commands work.

## Running the tests

To check your implementation, run `go test`:

```shell
$ go test ./...
```

This should give output like this:

```shell
ok  	github.com/3lbits/case-url-shortener	0.277s
```

## Running the service

The service runs as a small Go server. You can start it like this:

```shell
$ go run . -serve -addr ":9090" -linkfile links.txt
```

You can now send HTTP requests to the server:

```shell
$ curl -i localhost:9090/elbits
```

You should get something that looks like this:

```text
HTTP/1.1 500 Internal Server Error
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 04 Oct 2024 09:29:55 GMT
Content-Length: 14

implement me!
```

## The linkfile

Links are stored in plain text files.
Each line may contain either a comment, whitespace only, or a link definition.

1. Lines containing only whitespace are ignored.
2. Lines that start with a `#` character are ignored.
3. Link definitions consist of any non-space characters, a space character, and a valid URL.

Here's an example file:

```text
# Comments are useful
e https://elbits.no/

# lunch orders for Lede
lede-food-orders https://forms.google.com/sefyuaebrgoubaerogubaoeurfg82991384019401734
```

## Your task

1. Implement [`LinkFile.Long`](./linkfile.go). It doesn't need to be fancy or optimized, it just needs to work.
2. Deploy the application to `fly.io`.

### General guidelines

A working HTTP server is implemented in `main.go`. You shouldn't need to change it.

You _can_ modify any file except the existing test files in [`testdata`](./testdata) and [`linkfile_pass_test.go`](./linkfile_pass_test.go).
If you do find any bugs in `main.go` or `linkfile_pass_test.go`, feel free to fix them and explain what you fixed in the commit message (this is not a test, we don't expect you to find anything specific).

We'll look at the bundled commit when evaluating the assignment.
You can commit directly to `main`.

### Launch the app

Launch the application on `fly.io`. [Create an account and follow the quickstart](https://fly.io/docs/getting-started/launch/).

When told to run the `flyctl launch` command, you can provide this:

```shell
flyctl launch \
  --org "[UPDATE THIS]"
  --vm-size "shared-cpu-1x" \
  --region arn \
  --name urlshortener \
  --internal-port "8080" \
  --env "PORT=8080"
```

Please send your app URL (`flyctl deploy` or `flyctl launch` should print this) together with the code bundle.

### Send us the code

Run `sh done.sh` to generate `finished-case.bundle` and send it to `case-url-shortener@elbits.no`.

Enjoy!
