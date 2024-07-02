<br/>
<p align="center">
<img src="assets/go2tv-logo-color.svg" width="225" alt="Go2TV logo">
</a>
</p>
<br/>
<div align="center">
<p>

[![Go Report Card](https://goreportcard.com/badge/github.com/alexballas/Go2TV)](https://goreportcard.com/report/github.com/alexballas/Go2TV)
[![Release Version](https://img.shields.io/github/v/release/alexballas/Go2TV?label=Release)](https://github.com/alexballas/Go2TV/releases/latest)
[![Tests](https://github.com/alexballas/go2tv/actions/workflows/go.yml/badge.svg?branch=devel)](https://github.com/alexballas/go2tv/actions/workflows/go.yml)

[![Build for ARMv6 (32-bit)](https://github.com/alexballas/go2tv/actions/workflows/build-arm.yml/badge.svg?branch=devel)](https://github.com/alexballas/go2tv/actions/workflows/build-arm.yml)
[![Build for ARMv8 (64-bit)](https://github.com/alexballas/go2tv/actions/workflows/build-arm64.yml/badge.svg?branch=devel)](https://github.com/alexballas/go2tv/actions/workflows/build-arm64.yml)
[![Build for Android](https://github.com/alexballas/go2tv/actions/workflows/build-android.yml/badge.svg?branch=devel)](https://github.com/alexballas/go2tv/actions/workflows/build-android.yml)
[![Build for Linux](https://github.com/alexballas/go2tv/actions/workflows/build-linux.yml/badge.svg?branch=devel)](https://github.com/alexballas/go2tv/actions/workflows/build-linux.yml)
[![Build for MacOS Intel](https://github.com/alexballas/go2tv/actions/workflows/build-mac-intel.yml/badge.svg?branch=devel)](https://github.com/alexballas/go2tv/actions/workflows/build-mac-intel.yml)
[![Build for MacOS Apple Silicon](https://github.com/alexballas/go2tv/actions/workflows/build-mac.yml/badge.svg?branch=devel)](https://github.com/alexballas/go2tv/actions/workflows/build-mac.yml)
[![Build for Windows](https://github.com/alexballas/go2tv/actions/workflows/build-windows.yml/badge.svg?branch=devel)](https://github.com/alexballas/go2tv/actions/workflows/build-windows.yml)
</p>
Cast your media files to UPnP/DLNA Media Renderers and Smart TVs.
</div>

---
GUI mode
-----
![](https://i.imgur.com/Ga3hLJM.gif)

![](https://i.imgur.com/Pw44BYD.png)
![](https://i.imgur.com/JeUxGGd.png)

CLI mode
-----
![](https://i.imgur.com/BsMevHi.gif)

Parameters
-----
``` console
$ go2tv -h
Usage of go2tv:
  -l    List all available UPnP/DLNA Media Renderer models and URLs.
  -s string
        Local path to the subtitles file.
  -t string
        Cast to a specific UPnP/DLNA Media Renderer URL.
  -tc
        Use ffmpeg to transcode input video file.
  -u string
        HTTP URL to the media file. URL streaming does not support seek operations. (Triggers the CLI mode)
  -v string
        Local path to the video/audio file. (Triggers the CLI mode)
  -version
        Print version.
```

Allowed media files in the GUI
-----
- mp4, avi, mkv, mpeg, mov, webm, m4v, mpv, mp3, flac, wav, jpg, jpeg, png

This is a GUI only limitation.

Build requirements and dependencies
-----
- Go v1.19+
- ffmpeg (optional)

**Build using Docker**

Since the repo provides a [Dockerfile](./Dockerfile), you can build a Go2TV Docker image and run it with just Docker installed (no build requirements and deps above needed). Also, no Git repo cloning is needed (Docker will do it behind the scenes). Just issue:
``` console
$ docker build --force-rm [--pull] -t go2tv github.com/alexballas/go2tv#main
```
Notice the branch name after the `#`, as the above will build `main`. You can also build `devel` if you want to build the latest code. Usage under Docker is outside this document's scope, check Docker docs for more information, specially volume mounts and networking. [x11docker](https://github.com/mviereck/x11docker) might come handy to run GUI mode, although it's not tested, since main Docker usage is CLI.

Quick Start
-----
Download the app here https://github.com/alexballas/Go2TV/releases/latest. A single executable. No installation or external dependencies.

**Transcoding**

Go2TV supports live video transcoding, if ffmpeg is installed. When transcoding, SEEK operations are not available. Transcoding offers the maximum compatibility with the various file formats and devices. Only works with video files.

**MacOS potential issues**

If you get the "cannot be opened because the developer cannot be verified" error, you can apply the following workaround.
- Control-click the app icon, then choose Open from the shortcut menu.
- Click Open.

If you get the "go2tv is damaged and can't be opened. You should move it to the Bin." error you can apply the following workaround.
- Launch Terminal and then issue the following command: `xattr -cr /path/to/go2tv.app`.

Tested on
-----
- Samsung UE50JU6400
- Samsung UE65KS7000
- Android - BubbleUPnP app

Author
------

Alexandros Ballas <alex@ballas.org>

## Q & A
报错
```text
# github.com/go-gl/gl/v2.1/gl
# [pkg-config --cflags  -- gl gl]
Package gl was not found in the pkg-config search path.
Perhaps you should add the directory containing `gl.pc'
to the PKG_CONFIG_PATH environment variable
No package 'gl' found
Package gl was not found in the pkg-config search path.
Perhaps you should add the directory containing `gl.pc'
to the PKG_CONFIG_PATH environment variable
No package 'gl' found
# github.com/go-gl/glfw/v3.3/glfw
In file included from ./glfw/src/internal.h:188,
                 from ./glfw/src/context.c:30,
                 from /home/ap/data/software/go/pkg/mod/github.com/go-gl/glfw/v3.3/glfw@v0.0.0-20240307211618-a69d953ea142/c_glfw.go:4:
./glfw/src/x11_platform.h:33:10: fatal error: X11/Xlib.h: No such file or directory
   33 | #include <X11/Xlib.h>
      |          ^~~~~~~~~~~~
compilation terminated.
make: *** [Makefile:5：build] 错误 1
```
解决方法：`sudo apt-get install golang gcc libgl1-mesa-dev xorg-dev`