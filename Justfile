bin := 'show_slow_command_notif'

build:
  go build -o {{bin}} .

install: build
  cp {{bin}} ~/.bin/show_slow_command_notif

fmt:
  go fmt *.go
