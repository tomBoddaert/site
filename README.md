# Site

A simple static site builder

[Example site](https://gist.github.com/tomBoddaert/29b17843f88be23ea8e36bebf331834f)  
This is the recommended place to start if you just want to get going.

## Install

Make sure you have [`go`](https://go.dev/) installed.

```sh
go install github.com/tomboddaert/site@latest
```

## Use

The basic commands are [here](#using-the-tool).

The file structure should look like this:
```
.
├── out/                 - the built output
│   └── ...
├── pages/               - the templated pages
│   ├── <template name>/
│   │   └── ...
│   └── ...
├── raw/                 - the non-templated pages
│   └── ...
├── templates/           - the templates
│   ├── <single file
│   │    template name>
│   ├── <multi file
│   │    template name>/
│   │   └── ...
│   └── ...
├── data.json            - page variables
|    (optional)
├── siteConfig.json      - configuration
|    (optional)
└── tsconfig.json        - TypeScript configuration
     (optional)
```

### `siteConfig.json`

This file sets the config for site. The default config (do not use comments):
```js
{
  "TemplateDir":      "templates", // Sets the template directory
  "TemplatedSrcDir":  "pages",     // Sets the templated page directory
  "RawSrcDir":        "raw",       // Sets the raw directory
  "DstDir":           "out",       // Sets the output directory
  "DstMode":          "0755",      // Sets the file permissions for the output (octal)
  "DataFile":         "data.json", // Sets the file data is read from
  "FmtTemplatedHtml": false,       // Formats templated html files (may be slow)
  "FmtRawHtml":       false,       // Formats raw html files (may be slow)
  "TranspileTS":      true,        // Runs tsc to transpile TS
  "TSArgs":           [],          // Arguments to tsc (npx tsc [args])
  "IncludeTS":        false,       // Include TS files in the output

  "NotFoundPath":     "404"        // 404 page (serve)
}
```

### `out`

This is the default output directory, where the built pages will be.

### `templates`

This is the directory for templates; templates can either be a single file or a directory of files, which will be concatenated into one template. The name of the file or directory is the template name.

The templates use Go templating with [sprig functions](http://masterminds.github.io/sprig/). [Official documentation](https://pkg.go.dev/text/template); [Nomad Tutorial](https://developer.hashicorp.com/nomad/tutorials/templates/go-template-syntax).

### `pages`

This is the directory for pages that need templating. The pages should be in a subdirectory with the same name as the template.

### `rawPages`

This is the directory for files that do not need templating.

### `data.json`

This file is for variables used in the templates. The root should be an object, where the keys are the paths of the page in [`pages`](#pages), and the values are the data to be used in the page.  
There can also be a `default` key, whose corresponding values are used when not overridden.  
Defaults for a template can be set with a key corresponding to the name of the template, these take priority over the standard `default`.

### `tsconfig.json`

If TypeScript is used, this file is needed. The settings below need to be set. This will need to be adjusted if directory names are changed in `siteConfig.json`.

```json
{
  "compilerOptions": {
    "rootDir": "./raw",
    "outDir": "./out",
  },
  "include": [ "./raw" ]
}
```

The file can generated with:
```sh
npx tsc --init
```

`outFile` could be used to transpile a separate directory of TypeScript to one file of JavaScript. See the [tsconfig docs](https://aka.ms/tsconfig) for more information.
`rootDir` and `include` could be changed to separate TypeScript from the pages.

## Using the tool

### Help:
```sh
site help
```

### Building:
```sh
site build
```

### Serving:
```sh
site serve
```

#### Serving on a different address:
Linux:
```sh
SITE_ADDRESS=127.4.0.1:80 site serve
```

Windows cmd:
```cmd
set SITE_ADDRESS=127.4.0.1:80
site serve
```

PowerShell:
```powershell
$Env:SITE_ADDRESS = "127.4.0.1:80"
site serve
```

### Build and serve:
```sh
site build serve
```

## License

You maintain all rights to any work you create using the build functionality of this tool, and attribution is not required in the created work.

Licensed under the [MIT License](/LICENSE-MIT) ([source](https://opensource.org/licenses/MIT)) or [Apache License, Version 2.0](/LICENSE-APACHE) ([source](https://www.apache.org/licenses/LICENSE-2.0)) at your option.
