package bpf

var Set_src = `
#include <uapi/linux/ptrace.h>

typedef struct {
    u32 pid;
    u32 pad;
    u64 start_time_ns;
    u64 duration;
    int klen;
    char key[128];
} set_event_t;

BPF_HASH(setcall,u64,set_event_t);
BPF_PERF_OUTPUT(set_events);

int trace_start(struct pt_regs* ctx){
    u64 pid = bpf_get_current_pid_tgid();
    u64 start_time_ns = bpf_ktime_get_ns();
    
    robj* rObj = (robj*)PT_REGS_PARM2(ctx);

    set_event_t event = {
        .pid = pid>>32,
        .start_time_ns = start_time_ns,
    };

    event.klen = bpf_probe_read_user_str(&event.key,sizeof(event.key),(void*)rObj->ptr);
    setcall.update(&pid,&event);
    return 0;
}

int trace_end(struct pt_regs* ctx){
    u64 pid = bpf_get_current_pid_tgid();
    set_event_t* eventp = setcall.lookup(&pid);
    if (eventp==0){
        return 0;
    }
    set_event_t event = *eventp;
    event.duration = bpf_ktime_get_ns() - event.start_time_ns;
    set_events.perf_submit(ctx,&event,sizeof(event));
    setcall.delete(&pid);
    return 0;
}
`
