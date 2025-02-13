#!/usr/bin/env bash

project_name="aqueduct"

echo "Iniciando instalação"
echo ""

Root=$(realpath $0);
Root=$(dirname "$Root")

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
mv $Root/$project_name.service /etc/systemd/system

# Public files
mv -r $Root/public /etc/$project_name/                          

# Configuration files
mv $Root/config.env /etc/$project_name/

# Binaries
mv $Root/$project_name /usr/local/bin/

echo "--- Instalação concluída"
echo "--- A configuração padrão para Aqueduct foi definida, para uma configuração personalizada edite o arquivo de configuração em /etc/$project_name/config.env"
echo
echo "--- Para executar o programa digite: "
echo 
echo "  $ sudo ./aqueduct"
echo
echo "--- Para executar o programa sempre que o sistema for iniciado: "
echo 
echo "  [ Reiniciar o daemon gerenciador de serviços ]"
echo "  $ sudo systemctl daemon-reload "
echo
echo "  [ Habilitar programa sempre que o sistema iniciar ]"
echo "  $ sudo systemctl enable aqueduct.service"
echo 
echo "  [ Habilitar o programa agora ]" 
echo "  $ sudo systemctl start aqueduct.service"
