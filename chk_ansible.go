package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type AnsibleJob struct {
	ID        int    `json:"id"`
	Status    string `json:"status"`
	Started   string `json:"started"`
	Finished  string `json:"finished"`
	Created   string `json:"created"`
	Modified  string `json:"modified"`
	Failed    bool   `json:"failed"`
	FailedMsg string `json:"failed_msg"`
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run monitor_ansible.go <AnsibleTowerURL> <AnsibleTokenFile> <AnsibleJobTemplate>")
		os.Exit(1)
	}

	ansibleTowerURL := os.Args[1]
	ansibleTokenFile := os.Args[2]
	ansibleJobTemplate := os.Args[3]

	ansibleToken, err := readTokenFromFile(ansibleTokenFile)
	if err != nil {
		fmt.Printf("Failed to read Ansible token: %v\n", err)
		os.Exit(1)
	}

	if err := checkAnsibleHealth(ansibleTowerURL, ansibleToken); err != nil {
		fmt.Printf("Failed to check Ansible Tower health: %v\n", err)
		os.Exit(1)
	}

	job, err := getAnsibleJob(ansibleTowerURL, ansibleToken, ansibleJobTemplate)
	if err != nil {
		fmt.Printf("Failed to get job details: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Job ID: %d\nStatus: %s\nStarted: %s\nFinished: %s\nFailed: %t\n", job.ID, job.Status, job.Started, job.Finished, job.Failed)
	if job.Failed {
		fmt.Printf("Failure Message: %s\n", job.FailedMsg)
		os.Exit(1)
	}

	os.Exit(0)
}

func readTokenFromFile(filename string) (string, error) {
	tokenBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(tokenBytes), nil
}

func checkAnsibleHealth(ansibleTowerURL, ansibleToken string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/ping/", ansibleTowerURL), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ansibleToken))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ansible Tower is not healthy (status code: %d)", resp.StatusCode)
	}

	return nil
}

func getAnsibleJob(ansibleTowerURL, ansibleToken, ansibleJobTemplate string) (*AnsibleJob, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/job_templates/%s/last_job/", ansibleTowerURL, ansibleJobTemplate), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ansibleToken))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var job AnsibleJob
	err = json.Unmarshal(body, &job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}
