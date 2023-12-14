from scapy.all import sniff, DNS, DNSQR
import subprocess
import os

# Configuration
VPN_INTERFACE = 'tun0'  # Replace with your VPN interface
DOMAINS_TO_BLOCK = ['googlevideo.com', 'anotherdomain.com']  # List of domains to block
BLOCKLIST_PATH = '/etc/dnsmasq.blocklist'  # Path to dnsmasq blocklist

def is_subdomain(domains, subdomain):
    return any(subdomain.endswith('.' + domain) or subdomain == domain for domain in domains)

def update_dnsmasq_blocklist(subdomain):
    with open(BLOCKLIST_PATH, 'a+') as file:
        file.seek(0)
        if f"address=/{subdomain}/" not in file.read():
            file.write(f"address=/{subdomain}/\n")
            subprocess.run(["sudo", "systemctl", "restart", "dnsmasq"], check=True)
            print(f"Blocked {subdomain}")

def dns_query(packet):
    if packet.haslayer(DNS) and packet.getlayer(DNS).qr == 0:
        domain_name = str(packet[DNSQR].qname.decode('utf-8')).rstrip('.')
        if is_subdomain(DOMAINS_TO_BLOCK, domain_name):
            print(f"Detected DNS query for {domain_name}")
            update_dnsmasq_blocklist(domain_name)

# Make sure the dnsmasq blocklist file exists
os.makedirs(os.path.dirname(BLOCKLIST_PATH), exist_ok=True)
open(BLOCKLIST_PATH, 'a').close()

# Start packet sniffing
sniff(iface=VPN_INTERFACE, filter="port 53", prn=dns_query)
