# Guardian [![Build Status](https://travis-ci.org/monochromegane/guardian.svg?branch=master)](https://travis-ci.org/monochromegane/guardian)

Monitor file changes and execute custom commands for each event.
This tool supports cross-platform file system notifications by using [fsnotify](https://github.com/fsnotify/fsnotify).

## Usage

```sh
$ guardian -create "echo %e %p" /path/to/monitor [...]
```

`guardian` monitors file **recursively** and execute custom commands for each event.

### Placeholder

You can use placeholder in command arguments. `%e` will be replace by event name, `%p` will be replace by path name.

### Events

- CREATE
- WRITE
- CHMOD
- REMOVE
- RENAME

## License

[MIT](https://github.com/monochromegane/guardian/blob/master/LICENSE)

## Author

[monochromegane](https://github.com/monochromegane)
