package ansb

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

const template string = `
[instances:vars]
ansible_python_interpreter=/usr/bin/python3
ansible_user=root
ansible_ssh_extra_args='-o StrictHostKeyChecking=no'`

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func CreateAnsibleHostTemplate() {
	if FileExists("hosts") {
		os.Remove("hosts")
	}
	//_ = os.WriteFile("hosts", template, 0644)
	f, _ := os.Create("hosts")
	_, _ = f.WriteString(template)
	f.Sync()
	defer f.Close()
}

func AppendAnsible(val string) {
	//fmt.Printf("Called from output: %s \n", val)

	file, err := os.Open("hosts")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Read the entire file into memory
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Prepend new data to the in-memory buffer
	newData := []byte(val + "\n")
	data = append(newData, data...)

	// Open the file in write mode
	file, err = os.OpenFile("hosts", os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Write the entire buffer back to the file
	_, err = file.Write(data)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func RunPlayBook() {
	//Get Directroy-Path
	fmt.Println("Running Ansible Playbook / waiting for Resources to be ready...")
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	//Run Playbook in Bash
	cmd := exec.Command("bash", "-c", "sleep 60;ANSIBLE_PYTHON_INTERPRETER=auto_silent ANSIBLE_HOST_KEY_CHECKING=False ansible-playbook -i hosts ansb/instances.yaml")
	cmd.Dir = path

	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		fmt.Println("could not run command: ", err)
	}
}
