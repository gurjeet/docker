# ZFS Driver for Docker

This directory contains the storage driver for Docker that uses ZFS as the
backend.

# Usage

Create a zfs mount point using the following command (your ZFS pool name may be different; mine is zstore)

	sudo zfs create -o mountpoint=/var/lib/docker/zfs zstore/docker

Then launch Docker and force it to use the ZFS storage driver
(you may want to change your init/Upstart scripts, instead)

	docker -d -s zfs -p /var/run/docker.pid

# Hacking

In development mode, I use the following commands to build my version of Docker.
(I usually chain them together in a single command, but have split them here for clarity)

	sudo service docker start
	sleep 1
	sudo docker run -privileged -v `pwd`:/go/src/github.com/dotcloud/docker docker hack/make.sh binary
	RC=$?; echo build exit code: $RC
	sudo service docker stop

And then execute it like so

	TS=$(date '+%Y%m%d_%H%M%S')
	sudo ${PWD}/bundles/*/binary/docker* -d -D -s zfs -p /var/run/docker_g.pid > /tmp/docker.dev.$TS.log 2>&1 &
	less +F /tmp/docker.dev.$TS.log

The chained versions of the above commands are

	sudo service docker start; sleep 1; sudo docker run -privileged -v `pwd`:/go/src/github.com/dotcloud/docker docker hack/make.sh binary; RC=$?; echo build exit code: $RC; sudo service docker stop
	TS=$(date '+%Y%m%d_%H%M%S'); sudo ${PWD}/bundles/*/binary/docker* -d -D -s zfs -p /var/run/docker_g.pid > /tmp/docker.dev.$TS.log 2>&1 & :; less +F /tmp/docker.dev.$TS.log

And to terminate the docker binary launched above, use `fg` command to bring the
Docker process for foreground, and then hit `Ctrl-C` to terminate it.
