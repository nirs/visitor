# SPDX-FileCopyrightText: The RamenDR authors
# SPDX-License-Identifier: Apache-2.0

name := visitor
version := 0.4.0
pkgname := $(name)-$(version)

all: linux windows

linux:
	CGO_ENABLED=0 GOOS=$@ go build -o dist/$(name)-$@-amd64

windows:
	GOOS=$@ go build -o dist/$(name)-$@-amd64.exe

rpm:
	git archive --prefix $(pkgname)/ HEAD > $(pkgname).tar.gz
	rpmbuild -ta $(pkgname).tar.gz
