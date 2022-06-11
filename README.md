# Ugo

A **testing framework** for your tutorials.

<img src="https://i.imgflip.com/4ddj1y.jpg" width="400" />

Are your tutorials constantly broken or outdated? Are you tired of manually testing your tutorials to make sure the steps actually work?

Fear no more! **Ugo** is here to help automate that!

### Install

```bash
go get -u github.com/jromero/ugo/cmd/ugo
```

### Usage

Integrating **Ugo** into your tutorials is done by adding a few hidden HTML comments. You can then use the `ugo` CLI to test your tutorials. 

It's that easy!

Here's a quick look at what that may look like.

1. Create the tutorial:
~~~markdown

# My Tutorial

<!-- test:suite=my-tutorial -->

<!-- test:teardown:exec -->
<!--
```
rm some-file.txt
```
-->

First, let's create a file `some-file.txt` with the following content:

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

2. Run `ugo`:

```bash
ugo run
```

3. Output:

```text
[info ][*] Suite 'my-tutorial' executing...
[info ][*][my-tutorial] Running task #1
[info ][*][my-tutorial] Running task #2
[info ][*][my-tutorial] Running task #3
[info ][*][my-tutorial] Running task #4
[info ][*] Nothing broken. Good job!
```

For more, check out these [examples](docs/examples).
