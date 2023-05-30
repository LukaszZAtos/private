---
- name: Monitor Ansible host
  hosts: localhost
  gather_facts: true

  tasks:
    - name: Check if monitoring script is installed
      stat:
        path: /opt/monitoring_script.sh
      register: monitoring_script

    - name: Install monitoring script
      copy:
        src: monitoring_script.sh
        dest: /opt/monitoring_script.sh
        mode: 0755
      when: monitoring_script.stat.exists == False

    - name: Run monitoring script
      shell: /opt/monitoring_script.sh
      register: monitoring_result
      changed_when: false

    - name: Print monitoring result
      debug:
        var: monitoring_result.stdout_lines


#!/bin/bash

# Check Ansible controller connectivity
ansible_controller="ansible_controller"
ansible_controller_port=22

if nc -z -w 1 "$ansible_controller" "$ansible_controller_port"; then
  echo "Ansible controller ($ansible_controller) is reachable."
else
  echo "Ansible controller ($ansible_controller) is not reachable."
fi

# Check managed node connectivity
managed_nodes=("managed_node1" "managed_node2" "managed_node3")
managed_node_port=22

for node in "${managed_nodes[@]}"; do
  if nc -z -w 1 "$node" "$managed_node_port"; then
    echo "Managed node ($node) is reachable."
  else
    echo "Managed node ($node) is not reachable."
  fi
done

# Check Ansible controller process
ansible_controller_process="ansible-controller"
if pgrep -x "$ansible_controller_process" >/dev/null; then
  echo "Ansible controller process ($ansible_controller_process) is running."
else
  echo "Ansible controller process ($ansible_controller_process) is not running."
fi

# Check managed node processes
managed_node_processes=("nginx" "apache2" "postgres")
for process in "${managed_node_processes[@]}"; do
  if pgrep -x "$process" >/dev/null; then
    echo "Managed node process ($process) is running."
  else
    echo "Managed node process ($process) is not running."
  fi
done