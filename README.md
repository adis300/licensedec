# licensedec

This is a package to decorate all C-type source code with a Lisence templates and replace the existing template.

Written in golang


All arguments are optional

```
go get -u github.com/licensedec

licensedec -root=/YourAbsolutePath -license=MyLicense.txt -recursive=true -exts=m,h

```

## Meaning of all parameters
| Parameter |     Description     |  Default value |
|----------|:-------------:|:-----:|
| root |  Absolute path | pwd |
| license |    Path to your license template(relative or absolute)   |   LICENSE |
| recursive | Search sub folder? | false |
| clean | Clean temporary file? | true |
| exts | Valid file extensions to search and edit  | m,h,js,swift,go,cpp,mm,hh,hpp,mpp,java |

# Another example
```
cd yourSourceCodeDir

licensedec -recursive=true -license=MyLicense.txt
```
