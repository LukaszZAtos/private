import os
import json

def read_json_file(file_path):
    with open(file_path, 'r') as file:
        return json.load(file)

def write_json_file(data, file_path):
    with open(file_path, 'w') as file:
        json.dump(data, file)

def process_file(file_path, unique_keys):
    data = read_json_file(file_path)
    unique_entries = {}
    
    for entry in data:
        key = tuple(entry[k] for k in unique_keys)
        if key not in unique_entries:
            unique_entries[key] = entry

    return list(unique_entries.values())

def main():
    work_dir = '/home/a692596/backup2perf'
    output_dir = '/home/nadabackup/in'
    unique_keys = ['primary_key1', 'primary_key2']  # Define your primary keys

    for filename in os.listdir(work_dir):
        if filename.endswith('.json'):
            file_path = os.path.join(work_dir, filename)
            processed_data = process_file(file_path, unique_keys)
            output_file_path = os.path.join(output_dir, filename + '_processed.json')
            write_json_file(processed_data, output_file_path)

if __name__ == "__main__":
    main()
