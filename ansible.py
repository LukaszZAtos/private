import argparse
import requests
import urllib3
import sys

def get_last_job_id(api_url, api_token, template_id):
    headers = {
        "Authorization": f"Bearer {api_token}"
    }

    api_endpoint = f"{api_url}/api/v2/job_templates/{template_id}/jobs/"
    urllib3.disable_warnings()
    response = requests.get(api_endpoint, headers=headers, verify=False)

    if response.status_code == 200:
        launches = response.json().get("results", [])
        if launches:
            last_job_id = launches[-1]["id"]
            return last_job_id
        else:
            print("No jobs found for the specified template.")
            return None
    else:
        print(f"Failed to retrieve job launches. Status code: {response.status_code}")
        sys.exit(1)  # Exit with code 1 to indicate an error

def check_ansible_job_status(api_url, api_token, job_id):
    headers = {
        "Authorization": f"Bearer {api_token}"
    }

    api_endpoint = f"{api_url}/api/v2/jobs/{job_id}/"
    urllib3.disable_warnings()
    response = requests.get(api_endpoint, headers=headers, verify=False)

    if response.status_code == 200:
        job_details = response.json()
        job_status = job_details.get("status", "Unknown")
        print(f"Job ID: {job_id}, Status: {job_status}")
    else:
        print(f"Failed to retrieve job status. Status code: {response.status_code}")
        sys.exit(1)  # Exit with code 1 to indicate an error

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Check Ansible job status through API.")
    parser.add_argument("-u", "--url", type=str, help="API URL")
    parser.add_argument("-t", "--token", type=str, help="API Token")
    parser.add_argument("-i", "--jobid", type=int, help="Ansible Job ID")
    parser.add_argument("-tid", "--templateid", type=int, help="Ansible Job Template ID")
    args = parser.parse_args()

    if args.url and args.token:
        if args.jobid:
            check_ansible_job_status(args.url, args.token, args.jobid)
        elif args.templateid:
            last_job_id = get_last_job_id(args.url, args.token, args.templateid)
            if last_job_id:
                check_ansible_job_status(args.url, args.token, last_job_id)
            else:
                sys.exit(1)  # Exit with code 1 to indicate an error
    else:
        print("Please provide API URL, API Token, and either Ansible Job ID or Template ID using -u, -t, and -i or -tid parameters.")
        sys.exit(1)  # Exit with code 1 to indicate an error
