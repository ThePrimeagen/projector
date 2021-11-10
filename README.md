## Projector
A simple key value store per path.

### Building From Source
```
git clone git@github.com:ThePrimeagen/projector.git
cd projector

# Install it where you want
go build -o ~/.local/bin/projector ./cmd
```

### Commands
#### Available Commands
<pre>
prints out the configuration for the current director and its parents.
If you provide no arguments to projector, this is the default behavior
<b>projector</b> [--pwd] <b>print</b> [<i>key</i>]
<b>projector</b> [--pwd] [<i>key</i>]
</pre>

<pre>
will link the cwd (or pwd) to another directories config.
This is useful for mirroring configs
<b>projector</b> [--pwd] <b>link</b> <i>absolute/to/other/directory</i>
</pre>

<pre>
will remove any previously established link
<b>projector</b> [--pwd] <b>unlink</b>
</pre>

<pre>
delete a key in the current directories config.  Will not delete key out of
parent
<b>projector</b> [--pwd] <b>del</b> <i>key_to_delete</i>
</pre>

<pre>
adds key and value to the current directories config.
<b>projector</b> [--pwd] <b>add</b> <i>key_to_delete</i> <i>arguments</i> ...
</pre>

<pre>
shows the current path to the config file
<b>projector</b> [--pwd] <b>add</b> <i>key_to_delete</i> <i>arguments</i> ...
</pre>
