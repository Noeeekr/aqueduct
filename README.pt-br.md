Translated to EN:
[![en](https://img.shields.io/badge/lang-en-red.svg)](https://github.com/Noeeekr/aqueduct/blob/master/readme.md)

# Descrição
Aqueduct facilita a transferência de todos os tipos de arquivos entre computadores na mesma rede. Atualmente, está sendo desenvolvido na linguagem de programação GO para sistemas operacionais Linux.

# Instalação

> Baixe a versão mais recente para o seu sistema operacional [neste link](https://github.com/Noeeekr/aqueduct/releases).

Após baixar a pasta contendo os arquivos, extraia-a, encontre o arquivo chamado config.env e localize a variável "path", você colará o caminho para a pasta que você deseja compartilhar logo após ela na mesma linha.
```
    path=/pasta/para/compartilhar/
```
Depois disso, você está pronto para iniciar o Aqueduct com os seguintes comandos no seu terminal:
 
```
    $ chmod 764 ./install.sh
    $ sudo ./install.sh
```
# Configuração
Após instalar o Aqueduct, você pode realizar mais configurações iniciando o binário diretamente com diferentes flags. Você pode usar o comando abaixo para ver as flags disponíveis.
```
    $ ./aqueduct --help
```

# Contribua
Leia contributing.md para mais informações.