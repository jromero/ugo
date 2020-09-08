## Using bash

_This is an example tutorial that has the necessary **hidden** annotations for Ugo. To see the hidden annotations, view the [raw](https://raw.githubusercontent.com/jromero/ugo/main/docs/examples/bash-basics.md) file._

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

Now let's append a few lines:

<!-- test:exec -->
```bash
echo "line 2" >> file.txt
echo "line 3" >> file.txt
echo "line 4" >> file.txt
echo "line 5" >> file.txt
```

If we `cat` it again:

<!-- test:exec -->
```bash
cat file.txt
``` 

We should see all of the content, something like:

<!-- test:assert=contains;ignore-lines=... -->
```text
some-content
line 2
...
line 5
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