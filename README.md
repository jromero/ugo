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
[my-tutorial] Suite 'my-tutorial' executing...
[my-tutorial] Working directory: /var/folders/nx/x67fz2nj5hv_w43gn5h019hh0000gn/T/suite-my-tutorial-635763397
[my-tutorial][#1-default:file] --> Running task #1
[my-tutorial][#1-default:file] Writing file (some-file.txt) with contents:
some content
[my-tutorial][#2-default:exec] --> Running task #2
[my-tutorial][#2-default:exec] Executing the following:
cat some-file.txt
[my-tutorial][#2-default:exec] Output:
some content
[my-tutorial][#3-default:assert:contains] --> Running task #3
[my-tutorial][#3-default:assert:contains] Checking that output contained:
some content
[my-tutorial][#4-teardown:exec] --> Running task #4
[my-tutorial][#4-teardown:exec] Executing the following:
rm some-file.txt
[my-tutorial][#4-teardown:exec] Output:
Nothing broken. Good job!

```

For more, check out these [examples](docs/examples).
