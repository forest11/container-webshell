package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/forest11/container-webshell/models"
)

// 获取exec id
func GetDockerExecId(host string, port string, containerid string) (string, error) {
	log.SetFlags(log.Llongfile)
	client := &http.Client{}
	baseURL := fmt.Sprintf("http://%s:%s/containers/%s/exec", host, port, containerid)
	request, err := http.NewRequest("POST", baseURL,
		strings.NewReader("{\"Tty\": true, \"Cmd\": [\"/bin/sh\"], \"AttachStdin\": true, \"AttachStderr\": true, \"Privileged\": true, \"AttachStdout\": true}"),
	)
	if err != nil {
		log.Println(err)
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return "", err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
	}
	v := &models.DockerContect{}
	err = json.Unmarshal(content, v)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return v.Id, nil
}

// 启动exec
func StartDockerExec(conn net.Conn, host string, port string, execId string) {
	data := "{\"Tty\":true}"
	_, err := conn.Write([]byte(
		fmt.Sprintf("POST /exec/%s/start HTTP/1.1\r\nHost: %s\r\nContent-Type: application/json\r\nContent-Length: %s\r\n\r\n%s",
			execId, fmt.Sprintf("%s:%s", host, port), fmt.Sprint(len([]byte(data))), data)))

	if err != nil {
		log.Println(err)
	}

}

// 调整tty 窗口大小
func ResizeContainer(host string, port string, execId string, width string, height string) {
	baseURL := fmt.Sprintf("http://%s:%s/exec/%s/resize?h=%s&w=%s", host, port, execId, width, height)
	request, err := http.NewRequest("POST", baseURL, nil)
	if err != nil {
		log.Println(err)
	}
	request.Header.Set("Content-Type", "text/plain")
	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		log.Println(err)
	}
}
