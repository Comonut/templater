Templater 
==========

Templater is a CLI tool that makes it easier to make, use and share code and project templates.


```
Usage: templater <template folder> <output path> [flags]

Flags:
  -h, --help            help for templater
  -m, --mode string     Output writing mode - one of append, ignore, replace, merge (default "replace")
  -v, --values string   File containing parameter values
 ```

Rundown
---

<img src="https://user-images.githubusercontent.com/3967269/125961518-3b69f73e-cf32-4cd2-bcc6-06fcda53eed1.gif" width="600" height="600"/>


A templates constitutes a folder that contains a `templates` folder and a `params.yaml` file.

### Templates
Templates are defined using the [Go template format](https://pkg.go.dev/text/template) extended with [sprig functions library](https://github.com/Masterminds/sprig)
The template for the example shown above and [found in examples](./examples/helloworld) looks as follows
``` python
{{if (eq .timestamp "Yes")}}# Generated at {{now}}
{{end}}print("{{.greeting}} world from {{.name}}!")
```
The relative path of each folder is also a template in itself.
So for examples `/{{.filename}}` would render a file named with the filename variable, but if we named it and  `/{{uuidv4}}` the filename will be a different uuid on each execution. You can also use the templating on the folders that preceed the final, which can be useful for instance for generating packages.

### Parameters
To make it easier to use the template, you should create a `params.yaml` file which describes your input variables.
The `params.yaml` for the example above looks like this:
``` yaml
- param: file
  title: Generated file name
  type: textfield
  example: ex. script

- param: greeting
  title: Preffered greeting?
  type: choice
  options:
    - Hello
    - Greetings
    - What's up

- param: name
  title: Your name
  type: textfield
  example: Alex

- param: timestamp
  title: Add timestamp
  type: choice
  options:
    - "Yes"
    - "No"
```
This file is then consumed by the CLI to generate a simple terminal user interface that acts a setup wizard.

(TODO add docs for param types)

You can optionally omit this file and populate your parameters using a YAML file passed with the `-v` flag.

### Writing modes
By default, if there is already a file on the path of a new generated file, the new one will replace the old one.
However, you can change this beahviour using the `--mode` flag to one of `append, ignore, replace, merge`


| mode    | description                                                    | existing file content | new file content | result      |
|---------|----------------------------------------------------------------|-----------------------|------------------|-------------|
| append  | Appends the generated content to the end of the existing file  | a<br/> b<br/> c<br/>                | a<br/> X<br/> c<br/>          | a<br/> b<br/> c<br/>  a<br/> X<br/> c<br/>    |
| ignore  | Ignores the new content and let's the old one stay             | a<br/> b<br/> c<br/>                 | a<br/> X<br/> c<br/>             | a<br/> b<br/> c<br/>         |
| replace | Replaces the old content with the new one                      | a<br/> b<br/> c<br/>                 | a<br/> X<br/> c<br/>             | a<br/> X<br/> c<br/>         |
| merge   | Merges the two using a 3 way merge with the template as origin | a<br/> b<br/> c<br/>               | a<br/> X<br/> c<br/>              | a<br/> b<br/>X<br/> c<br/>       |

