# check if vmlinux.h is already present and is a non-empty file
if [ ! -f vmlinux.h ] || [ ! -s vmlinux.h ]; then
    echo "vmlinux.h not found. generating..."
    bpftool btf dump file /sys/kernel/btf/vmlinux format c >vmlinux.h
fi
