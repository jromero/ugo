# Ugo

A **testing framework** for your tutorials.

Are you tired of manually testing your tutorials to make sure the steps actually work? Are they broken or outdated? 

Fear no more! **Ugo** is here to help automate that!

### Install

```bash
go get -u github.com/jromero/ugo/cmd/ugo
```

### Usage

Integrating **Ugo** into your tutorials is done by just adding a few hidden HTML comments. You can then use the `Ugo` CLI to test your tutorials. 

It's that easy!

Here's a quick look at what that may look like.

~~~markdown

# My Tutorial

<!-- test:suite=my-tutorial -->

First, let's create a file `some-file.txt`

<!-- test:file=some-file.txt -->
```text
some content
```

Then, we'll execute `cat` to read the file:

<!-- test:exec -->
```bash
cat some-file.txt
```

Finally, we're make sure the output contains what we expect:

<!-- test:assert=contains -->
```text
some content
```

~~~

For more, check out these [examples](docs/examples).