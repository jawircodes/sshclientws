#!/bin/bash
#installer Tunnaapi 

cd

#Install Tunnapi
wget -O /usr/bin/sshclientws hhttps://raw.githubusercontent.com/jawircodes/sshclientws/main/sshclientws


#izin permision
chmod +x /usr/bin/sshclientws

#System tunnapi
wget -O /etc/systemd/system/sshclientws.service hhttps://raw.githubusercontent.com/jawircodes/sshclientws/main/sshclientws.service && chmod +x /etc/systemd/system/sshclientws.service

#restart service
systemctl daemon-reload

#Enable & Start & Restart 
systemctl enable sshclientws.service
systemctl start sshclientws.service
systemctl restart sshclientws.service

