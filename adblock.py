from scapy.all import sniff, IP, UDP, DNS, DNSQR
import re
import subprocess

def dns_query(packet):
    if packet.haslayer(DNS) and packet.getlayer(DNS).qr == 0:  # DNS query
        domain = str(packet[DNSQR].qname)
        if 'pattern_of_interest' in domain:  # Check for your domain pattern
            update_dnsmasq_blocklist(domain)

def update_dnsmasq_blocklist(domain):
    blocklist_path = "/etc/dnsmasq.blocklist"  # Path to your dnsmasq blocklist file
    block_entry = f"address=/{domain}/#\n"     # Format for dnsmasq blocklist entry

    try:
        # Check if the domain is already in the blocklist
        with open(blocklist_path, 'r+') as file:
            if block_entry not in file.read():
                file.write(block_entry)
                print(f"Added {domain} to blocklist.")

        # Restart dnsmasq to apply the changes
        subprocess.run(["sudo", "systemctl", "restart", "dnsmasq"], check=True)
        print("dnsmasq restarted successfully.")
    except Exception as e:
        print(f"Error updating dnsmasq blocklist: {e}")
    # Function to update dnsmasq blocklist
    # This involves file operations and possibly executing system commands

# Start sniffing on VPN interface (replace 'vpn_interface_name' with your actual VPN interface)
sniff(iface='tun0', filter='tcp port 443', prn=dns_query)
