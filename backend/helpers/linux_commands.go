package helpers

import (
	"fmt"
	"os/exec"
	"samey/config"
)

// CreateLinuxUserWithPassword creates a Linux user on the host system with sudo privileges.
func CreateLinuxUserWithPassword(username, publicKey string) error {
	sudoPassword := config.GetSudoPassword()

	// Command to create a Linux user with sudo access
	cmdCreateUser := fmt.Sprintf("sudo useradd -m -s /bin/bash -G sudo %s", username)
	if err := runHostCommand(cmdCreateUser, sudoPassword); err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	// Command to add public key to authorized_keys
	cmdAddSSHKey := fmt.Sprintf(`sudo mkdir -p /home/%s/.ssh && echo '%s' | sudo tee -a /home/%s/.ssh/authorized_keys && sudo chown -R %s:%s /home/%s/.ssh && sudo chmod 600 /home/%s/.ssh/authorized_keys`, username, publicKey, username, username, username, username, username)
	if err := runHostCommand(cmdAddSSHKey, sudoPassword); err != nil {
		return fmt.Errorf("failed to add SSH key: %v", err)
	}

	return nil
}

// runHostCommand runs a command on the host system outside the Docker container, passing the sudo password
func runHostCommand(command, sudoPassword string) error {
	// Prepend sudo password and run the command
	cmd := exec.Command("bash", "-c", fmt.Sprintf("echo '%s' | sudo -S %s", sudoPassword, command))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command error: %v, output: %s", err, string(output))
	}
	return nil
}
