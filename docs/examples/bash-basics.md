## Using bash

> **NOTE:** To see the hidden data view the raw file.

<!-- test:suite=bash-basics -->

### Printing to terminal

This is how you print to the terminal using `echo`:

<!-- test:exec -->
```bash
echo "hello world!"
```

You should see:

<!-- test:assert=contains -->
```text
hello world!
```

### Writing and reading files

First let's create a file, we'll use echo and direct the output to a new file.

<!-- test:exec -->
```bash
echo "some-content" > file.txt
```

To read it we'll use `cat`:

<!-- test:exec -->
```bash
cat file.txt
```

You should see:

<!-- test:assert=contains -->
```text
some-content
```

### Executing a script

We'll create a file `my-script.sh` with contents:

<!-- test:file=my-script.sh -->
```shell script
echo "hello from script!"
```

Next, let's make it executable by running:

<!-- test:exec -->
```bash
chmod +x my-script.sh
```

Lastly, let's execute it:

<!-- test:exec -->
```bash
./my-script.sh
```

You should see:

<!-- test:assert=contains -->
```text
hello from script!
```