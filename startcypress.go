package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func startXvnc() *exec.Cmd {
	xvnc := exec.Command(
		"Xvnc",
		os.Getenv("DISPLAY"),
		"-alwaysshared",
		"-depth",
		"16",
		"-geometry",
		os.Getenv("VNC_GEOMETRY"),
		"-securitytypes",
		"none",
		"-auth",
		fmt.Sprintf("%s/.Xauthority", os.Getenv("HOME")),
		"-fp",
		"catalogue:/etc/X11/fontpath.d",
		"-pn",
		"-rfbport",
		os.Getenv("VNC_PORT"),
		"-rfbwait",
		"30000")
	fmt.Println("Starting Xvnc")
	xvnc.Start()
	return xvnc
}

func waitForPort() {
	n := 1
	for n < 3 {
		conn, _ := net.DialTimeout("tcp", net.JoinHostPort("localhost", os.Getenv("VNC_PORT")), time.Second)
		if conn != nil {
			conn.Close()
			break
		}
		n++
		time.Sleep(time.Second)
	}
}

func startFluxbox() *exec.Cmd {
	fluxbox := exec.Command("fluxbox")
	fmt.Println("Starting fluxbox")
	fluxbox.Start()
	return fluxbox
}

func printCypressOutput(cypressStdout io.ReadCloser) {
	scanner := bufio.NewScanner(cypressStdout)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}
}

func startCypress() *exec.Cmd {
	args := os.Args[1:]
	args = append([]string{"cypress"}, args...)
	cypress := exec.Command("npx", args...)
	cypressStdout, _ := cypress.StdoutPipe()
	go printCypressOutput(cypressStdout)
	cypress.Start()
	return cypress
}

func startProcesses() (*exec.Cmd, *exec.Cmd, *exec.Cmd) {
	xvnc := startXvnc()
	waitForPort()
	fluxbox := startFluxbox()
	cypress := startCypress()
	return xvnc, fluxbox, cypress
}

func stopProcesses(commands ...*exec.Cmd) {
	for _, command := range commands {
		fmt.Println("Stopping", command.Args[0])
		command.Process.Kill()
		command.Wait()
	}
}

func waitForSignals(done chan bool) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	done <- true
}

func waitForCypress(cypress *exec.Cmd, done chan bool) {
	cypress.Wait()
	done <- true
}

func waitForSignalsOrCypress(xvnc *exec.Cmd, fluxbox *exec.Cmd, cypress *exec.Cmd) {
	done := make(chan bool, 1)
	go waitForCypress(cypress, done)
	go waitForSignals(done)
	<-done
	stopProcesses(xvnc, fluxbox, cypress)
}

func main() {
	xvnc, fluxbox, cypress := startProcesses()
	waitForSignalsOrCypress(xvnc, fluxbox, cypress)
	os.Exit(cypress.ProcessState.ExitCode())
}
