# Ugo

A **testing framework** for your tutorials.

Are you tired of manually testing your tutorials to make sure the steps actually work? Are they broken or outdated? 

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
[my-tutorial] Suite 'my-tutorial' executing...
[my-tutorial] Working directory: /var/folders/nx/x67fz2nj5hv_w43gn5h019hh0000gn/T/suite-my-tutorial-297835053
[my-tutorial][task#1] --> Running task #1
[my-tutorial][task#1] Writing file (some-file.txt) with contents:
some content
[my-tutorial][task#2] --> Running task #2
[my-tutorial][task#2] Executing the following:
cat some-file.txt
[my-tutorial][task#2] Output:
some content
[my-tutorial][task#3] --> Running task #3
[my-tutorial][task#3] Checking that output contained:
some content
Nothing broken. Good job!
```

For more, check out these [examples](docs/examples).
