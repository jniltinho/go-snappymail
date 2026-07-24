#!/usr/bin/env bash
# Unattended install of Zimbra FOSS 10.1.17.p2 (maldua build) on Ubuntu 24.04.
# Phase 1: install.sh -s (packages only, keystroke-fed)
# Phase 2: zmsetup.pl -c (config file)
set -euo pipefail

ZCS="zcs-10.1.17_GA_4200002.UBUNTU24_64.20260707094636"
URL="https://github.com/maldua/zimbra-foss/releases/download/zimbra-foss-build-ubuntu-24.04/10.1.17.p2/${ZCS}.tgz"
FQDN="mail.zimbra.test"
DOMAIN="zimbra.test"
ADMINPASS="Password1@"

if [ -d /opt/zimbra/bin ] && su - zimbra -c "zmcontrol status" 2>/dev/null | grep -q Running; then
  echo "Zimbra already installed and running — skipping."
  exit 0
fi

export DEBIAN_FRONTEND=noninteractive
apt-get update -qq
apt-get install -y -qq net-tools netcat-openbsd libidn12 libpcre3 libgmp10 libaio1t64 \
  libexpat1 libstdc++6 openssh-server rsyslog perl unzip 2>/dev/null || \
  apt-get install -y -qq net-tools netcat-openbsd libpcre3 libgmp10 libexpat1 libstdc++6 \
  openssh-server rsyslog perl unzip

# ── Fetch tarball (cache-aware) ────────────────────────────────
if [ ! -f "/zcs-cache/${ZCS}.tgz" ]; then
  echo "Downloading ${ZCS}.tgz (~236MB)..."
  curl -sL -o "/zcs-cache/${ZCS}.tgz" "$URL"
fi
mkdir -p /tmp/zcs
tar xzf "/zcs-cache/${ZCS}.tgz" -C /tmp/zcs

# ── Phase 1: packages (keystrokes: license y, repo y, pkgs default, dnscache n, modify y)
cd /tmp/zcs/${ZCS}
cat > /tmp/zcs-answers <<'EOF'
y
y
y
y
y
n
y
y
y
y
y
y
y
EOF
./install.sh -s --skip-activation-check < /tmp/zcs-answers

# ── Phase 2: zmsetup with config file ──────────────────────────
cat > /tmp/zmsetup.cfg <<EOF
AVDOMAIN="${DOMAIN}"
AVUSER="admin@${DOMAIN}"
CREATEADMIN="admin@${DOMAIN}"
CREATEADMINPASS="${ADMINPASS}"
CREATEDOMAIN="${DOMAIN}"
DOCREATEADMIN="yes"
DOCREATEDOMAIN="yes"
DOTRAINSA="yes"
EXPANDMENU="no"
HOSTNAME="${FQDN}"
HTTPPORT="8080"
HTTPPROXY="TRUE"
HTTPPROXYPORT="80"
HTTPSPORT="8443"
HTTPSPROXYPORT="443"
IMAPPORT="7143"
IMAPPROXYPORT="143"
IMAPSSLPORT="7993"
IMAPSSLPROXYPORT="993"
INSTALL_WEBAPPS="service zimlet zimbra zimbraAdmin"
JAVAHOME="/opt/zimbra/common/lib/jvm/java"
LDAPADMINPASS="${ADMINPASS}"
LDAPAMAVISPASS="${ADMINPASS}"
LDAPBESSEARCHSET="set"
LDAPHOST="${FQDN}"
LDAPPORT="389"
LDAPPOSTPASS="${ADMINPASS}"
LDAPREPPASS="${ADMINPASS}"
LDAPREPLICATIONTYPE="master"
LDAPROOTPASS="${ADMINPASS}"
LDAPSERVERID="2"
MAILBOXDMEMORY="3072"
MAILPROXY="TRUE"
MODE="https"
MYSQLMEMORYPERCENT="30"
POPPORT="7110"
POPPROXYPORT="110"
POPSSLPORT="7995"
POPSSLPROXYPORT="995"
PROXYMODE="https"
REMOVE="no"
RUNARCHIVING="no"
RUNAV="no"
RUNCBPOLICYD="no"
RUNDKIM="yes"
RUNSA="yes"
RUNVMHA="no"
SERVICEWEBAPP="yes"
SMTPDEST="admin@${DOMAIN}"
SMTPHOST="${FQDN}"
SMTPNOTIFY="yes"
SMTPSOURCE="admin@${DOMAIN}"
SNMPNOTIFY="no"
SNMPTRAPHOST="${FQDN}"
SPELLURL="http://${FQDN}:7780/aspell.php"
STARTSERVERS="yes"
STRICTSERVERNAME="yes"
SYSTEMMEMORY="10.0"
TRAINSAHAM="ham.acct@${DOMAIN}"
TRAINSASPAM="spam.acct@${DOMAIN}"
UIWEBAPPS="yes"
UPGRADE="yes"
VERSIONUPDATECHECKS="FALSE"
VIRUSQUARANTINE="virus-quarantine@${DOMAIN}"
ldap_dit_base_dn_config="cn=zimbra"
mailboxd_directory="/opt/zimbra/mailboxd"
mailboxd_keystore="/opt/zimbra/mailboxd/etc/keystore"
mailboxd_keystore_password="${ADMINPASS}"
mailboxd_truststore_password="changeit"
ssl_default_digest="sha256"
zimbraDNSMasterIP=""
zimbraDNSTCPUpstream="no"
zimbraDNSUseTCP="yes"
zimbraDNSUseUDP="yes"
zimbraDefaultDomainName="${DOMAIN}"
zimbraFeatureBriefcasesEnabled="Enabled"
zimbraFeatureTasksEnabled="Enabled"
zimbraIPMode="ipv4"
zimbraMailProxy="TRUE"
zimbraMtaMyNetworks="127.0.0.0/8 192.168.56.0/24 [::1]/128"
zimbraPrefTimeZoneId="America/Sao_Paulo"
zimbraReverseProxyLookupTarget="TRUE"
zimbraSSLCommonName="${FQDN}"
zimbraVersionCheckSendNotifications="FALSE"
EOF

/opt/zimbra/libexec/zmsetup.pl -c /tmp/zmsetup.cfg

# Classic UI as default for the web client (matches the go-snappymail zimbra skin)
su - zimbra -c "zmprov mc default zimbraPrefClientType advanced" || true

echo "02-zimbra OK — https://192.168.56.30/ (admin@${DOMAIN} / ${ADMINPASS})"
