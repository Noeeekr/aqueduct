Traduzido em PT-BR:
[![pt-br](https://img.shields.io/badge/lang-pt--br-green.svg)](https://github.com/Noeeekr/aqueduct/blob/master/README.pt-br.md)

# Description
Aqueduct makes easy to transfer all types of files between desktops in the same network. It is currently being developed with GO programming language for linux operating systems and .

# Installation

> Download the latest version for your operating system [in this link](https://github.com/Noeeekr/aqueduct/releases).

After downloading the folder containing the files, extract the folder, find the file called config.env and locate the variable "path", you'll paste the path to the folder you want to share after it in the same line.
```
    path=/my/shared/folder/path
```
After that you're ready to start aqueduct with the following commands in your terminal:
 
```
    $ chmod 764 ./install.sh
    $ sudo ./install.sh
```
# Configuration
After installing aqueduct, you can do more configurations by starting the binary yourself with different flags. You can use the command below to see the available flags.
```
    $ ./aqueduct --help 
```

# Contribute
Read contributing.md for more information.
