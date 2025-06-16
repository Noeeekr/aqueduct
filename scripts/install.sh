#!/usr/bin/env bash

project_name="aqueduct"

echo "Iniciando instalação"
echo ""

Root=$(realpath $0);
Root=$(dirname "$Root")

Path=$(grep path $Root/config.env | cut -d "=" -f2)
if [ -z "$Path" ]; then
    echo "Para instalar o programa, defina o caminho da pasta de compartilhamento no arquivo config.env"
    exit 0
fi

if [ ! -d /etc/$project_name ]; then 
    mkdir /etc/$project_name
fi

if [ ! -d /var/log/$project_name ]; then
    mkdir /var/log/$project_name
fi

if [ ! -f /var/log/$project_name/output.log ]; then
    touch /var/log/$project_name/output.log
fi

# Systemd
cp -r $Root/$project_name.service /etc/systemd/system

# Public files
cp -r $Root/public /etc/$project_name/                          

# Configuration files
cp -r $Root/config.env /etc/$project_name/

# Binaries
cp -r $Root/$project_name /usr/local/bin/

echo "Instalação concluída"
echo "Para informações sobre como executar o programa digite: aqueduct"