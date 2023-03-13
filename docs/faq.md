# FAQ

## Using sunbeam as a filter

There is two ways to wire sunbeam to your programs:

- Wrap your command in sunbeam: `sunbeam run <my-command>`
- Pipe data to sunbeam: `<my-command> | sunbeam read -`

Both methods have advantages and drawbacks:

- Using the wrap method, sunbeam controls the whole process. However, it requires your user to invoke sunbeam instead of your command, so you will loose completions.
- The pipe method is more flexible and easier to integrate with existing programs, but you loose the ability to reload your script since sunbeam is not responsible for generating the data.

## Detecting that a script is running in sunbeam

Sunbeam set the `SUNBEAM_RUNNER` environment variable to `true` when it's running a script. You can use it to adapt the output of your script depending on the context.

## Configuring sunbeam appearance

You can configure the appearance of sunbeam by setting the following environment variables:

- `SUNBEAM_HEIGHT`: The maximum height of the sunbeam window, in lines. Defaults to `0` (fullscreen).
- `SUNBEAM_PADDING`: The padding around the sunbeam window. Defaults to `0`.

You can also set these options using the `--height` and `--padding` flags.

```bash
sunbeam --height 20 --padding 2 ./github.sh
```

## Validating the output of a script

To validate the output of a script, you can use the validate command:

```bash
# Validate a static page
sunbeam validate sunbeam.json

# validate a dynamic page
./github.sh | sunbeam validate

# You can even chain it with other sunbeam commands !
sunbeam run ./github.sh | sunbeam validate
sunbeam read sunbeam.json | sunbeam validate
```

The validate command will exit with a non-zero exit code if the output is invalid.