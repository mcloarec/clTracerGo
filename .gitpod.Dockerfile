FROM gitpod/workspace-full-vnc

RUN sudo mkdir /docker-content
RUN sudo apt-get update
RUN sudo apt-get install nasm
RUN sudo apt-get install qemu-system-gui