#!/usr/bin/env bash

project_name="aqueduct"

echo "Iniciando instalação"
echo ""

Root=$(realpath $0);
Root=$(dirname "$Root")

if [ ! -d /etc/fileserver ]; then 
    mkdir /etc/$project_name
fi

if [ ! -d /var/log/fileserver ]; then
    mkdir /var/log/$project_name
fi

if [ ! -f /var/log/fileserver/output.log ]; then
    touch /var/log/$project_name/output.log
fi

cp -r $Root/files/* /etc/$project_name
cp -r $Root/$project_name /usr/local/bin/$project_name

echo "--- Instalação concluída"
echo "--- Agora defina os valores de configuração em /etc/$project_name/config.conf" 