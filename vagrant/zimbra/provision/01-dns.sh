#!/usr/bin/env bash
# Local DNS for Zimbra: A + MX for zimbra.test via dnsmasq on 127.0.0.1.
set -euo pipefail

IP="192.168.56.30"
FQDN="mail.zimbra.test"
DOMAIN="zimbra.test"

# /etc/hosts: remove 127.0.x.x hostname lines, pin the FQDN
sed -i "/${FQDN}/d; /127.0.1.1/d" /etc/hosts
echo "${IP} ${FQDN} mail" >> /etc/hosts

export DEBIAN_FRONTEND=noninteractive
apt-get update -qq
apt-get install -y -qq dnsmasq >/dev/null

# Free port 53 from systemd-resolved
mkdir -p /etc/systemd/resolved.conf.d
cat > /etc/systemd/resolved.conf.d/no-stub.conf <<EOF
[Resolve]
DNSStubListener=no
EOF
systemctl restart systemd-resolved

cat > /etc/dnsmasq.d/zimbra.conf <<EOF
domain=${DOMAIN}
local=/${DOMAIN}/
address=/${FQDN}/${IP}
mx-host=${DOMAIN},${FQDN},10
server=1.1.1.1
server=8.8.8.8
listen-address=127.0.0.1
bind-interfaces
EOF

systemctl enable --now dnsmasq
systemctl restart dnsmasq

rm -f /etc/resolv.conf
echo "nameserver 127.0.0.1" > /etc/resolv.conf

# sanity
host -t MX ${DOMAIN} 127.0.0.1
host ${FQDN} 127.0.0.1
echo "01-dns OK"
