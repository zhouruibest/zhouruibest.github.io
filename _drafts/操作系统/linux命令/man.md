# DESCRIPTION
This page describes the conventions that should be employed when writing man pages for the Linux man-pages project, which documents the user-space API provided by the Linux kernel and the GNU C library. The project thus provides most of the pages in Section 2, many of the pages that appear in Sections 3, 4, and 7, and a few of the pages that appear in Sections 1, 5, and 8 of the man pages on a Linux system. The conventions described on this page may also be useful for authors writing man pages for other projects.

# Sections of the manual pages
The manual Sections are traditionally defined as follows:

## 1 User commands (Programs)
Commands that can be executed by the user from within a shell.
## 2 System calls
Functions which wrap operations performed by the kernel.
## 3 Library calls
All library functions excluding the system call wrappers (Most of the libc functions).
## 4 Special files (devices)
Files found in /dev which allow to access to devices through the kernel.
## 5 File formats and configuration files
Describes various human-readable file formats and configuration files.
## 6 Games
Games and funny little programs available on the system.
## 7 Overview, conventions, and miscellaneous
Overviews or descriptions of various topics, conventions, and protocols, character set standards, the standard filesystem layout, and miscellaneous other things.
## 8 System management commands
Commands like mount(8), many of which only root can execute.