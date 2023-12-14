from scapy.all import sniff, DNS, DNSQR
import subprocess
import os
import time

# Configuration
VPN_INTERFACE = 'tun0'
DOMAINS_TO_BLOCK = ['googlevideo.com', 'anotherdomain.com']
BLOCKLIST_PATH = '/etc/dnsmasq.blocklist'
RESTART_INTERVAL = 60  # Time in seconds between dnsmasq restarts

last_restart_time = time.time()

def is_subdomain(domains, subdomain):
    return any(subdomain.endswith('.' + domain) or subdomain == domain for domain in domains)

def update_dnsmasq_blocklist(subdomain):
    global last_restart_time
    updated = False
    with open(BLOCKLIST_PATH, 'a+') as file:
        file.seek(0)
        if f"address=/{subdomain}/" not in file.read():
            file.write(f"address=/{subdomain}/\n")
            updated = True
            print(f"Blocked {subdomain}")

    current_time = time.time()
    if updated and (current_time - last_restart_time > RESTART_INTERVAL):
        subprocess.run(["sudo", "systemctl", "restart", "dnsmasq"], check=True)
        last_restart_time = current_time
        print("dnsmasq restarted successfully.")

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
