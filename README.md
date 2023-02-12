# Site

A simple static site builder

[Example site](https://gist.github.com/tomBoddaert/eeb098d80db8f16accb2efda6b68182a)  
This is the recommended place to start if you just want to get going.

## Install

If you have [`go`](https://go.dev/) installed:
```sh
go install github.com/tomboddaert/site@latest
```

There are prebuilt binaries on the releases page.

[Building from source](#building-the-tool)

## Use

The basic commands are [here](#using-the-tool).

The file structure should look like this:
```
.
├── docs                 - the built output
├── rawPages             - the non-templated pages
│   └── ...
├── templatedPages       - the templated pages
│   ├── <template name>
│   │   └── ...
│   └── ...
├── templates            - the templates
│   ├── <single file
│   │    template name>
│   ├── <multi file
│   │    template name>
│   │   └── ...
│   └── ...
├── pageVariables.json   - page variables
|    (optional)
└── tsconfig.json        - TypeScript configuration
     (optional)
```

### `docs`

This is the output directory, where the built pages will be.

### `templates`

This is the directory for templates; templates can either be a single file or a directory of files, which will be concatenated into one template. The name of the file or directory is the template name.

The templates use Go templating. [Official documentation](https://pkg.go.dev/text/template); [Nomad Tutorial](https://developer.hashicorp.com/nomad/tutorials/templates/go-template-syntax).

### `templatedPages`

This is the directory for pages that need templating. The pages should be in a subdirectory with the same name as the template. Inside there should be the structure wanted in [`docs`](#docs).

### `rawPages`

This is the directory for files that do not need templating. The files should be in the structure wanted in [`docs`](#docs).

### `pageVariables.json`

This file is for variables used in the templates. The root should be an object, where the keys are the paths of the page in [`templatedPages`](#templatedpages), and the values are the data to be used in the page.  
There can also be a `default` key, whose corresponding value is used in any page that is not in the file.  
Defaults for a template can be set with a key corresponding to the name of the template.

### `tsconfig.json`

If TypeScript is used, this file is needed. The settings below need to be set.

```json
{
  "compilerOptions": {
    "rootDir": "./rawPages",  // recommended but can be changed
    "outDir": "./docs",
  },
  "include": [ "rawPages" ]
}
```

`outFile` could be used to transpile a separate directory of TypeScript to one file of JavaScript. See the [tsconfig docs](https://aka.ms/tsconfig) for more information.
`rootDir` and `include` could be changed to separate TypeScript from the pages.

## Building the tool

```sh
go get github.com/tomboddaert/site
go build
```

## Using the tool

(If the `site` binary is in the same directory rather than in installed globally, change `site` => `./site` in the commands below)

Help:
```sh
site help
```

Building:
```sh
site build
```

Serving:
```sh
site serve
```

Serving on a different address:  
Linux:
```sh
SITE_ADDRESS=127.4.0.1:80 site serve
```

Windows:
```cmd
set SITE_ADDRESS=127.4.0.1:80
site serve
```

Build and serve:
```sh
site build serve
```

## Adding to PATH

On Linux, I prefer to add a symlink to the binary from the local bin directory, which I have in PATH.
```sh
ln -s <path to>/site/site ~/.local/bin/site
```

## License

You maintain all rights to any work you create using the build functionality of this tool, and attribution is not required in the created work.
However, the code itself is licensed under <a href="http://creativecommons.org/licenses/by/4.0/?ref=chooser-v1" target="_blank" rel="license noopener noreferrer" style="display:inline-block;">CC BY 4.0<img style="height:22px!important;margin-left:3px;vertical-align:text-bottom;" src="https://mirrors.creativecommons.org/presskit/icons/cc.svg?ref=chooser-v1"><img style="height:22px!important;margin-left:3px;vertical-align:text-bottom;" src="https://mirrors.creativecommons.org/presskit/icons/by.svg?ref=chooser-v1"></a>, where attribution is required.

<p xmlns:cc="http://creativecommons.org/ns#" xmlns:dct="http://purl.org/dc/terms/"><a property="dct:title" rel="cc:attributionURL" href="https://github.com/tomboddaert/site">Site</a> by <a rel="cc:attributionURL dct:creator" property="cc:attributionName" href="https://tomboddaert.com/">Tom Boddaert</a> is licensed under <a href="http://creativecommons.org/licenses/by/4.0/?ref=chooser-v1" target="_blank" rel="license noopener noreferrer" style="display:inline-block;">CC BY 4.0<img style="height:22px!important;margin-left:3px;vertical-align:text-bottom;" src="https://mirrors.creativecommons.org/presskit/icons/cc.svg?ref=chooser-v1"><img style="height:22px!important;margin-left:3px;vertical-align:text-bottom;" src="https://mirrors.creativecommons.org/presskit/icons/by.svg?ref=chooser-v1"></a></p>
