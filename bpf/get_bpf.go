package bpf

var Get_src = `
#include <uapi/linux/ptrace.h>

#define LRU_BITS 24

typedef struct {
    u32 pid;
    u32 pad;
    u64 start_time_ns;
    u64 duration;
    int klen;
    char key[128];
} get_event_t;

typedef struct redisObject {
    unsigned type:4;
    unsigned encoding:4;
    unsigned lru:LRU_BITS; /* LRU time (relative to global lru_clock) or
                            * LFU data (least significant 8 bits frequency
                            * and most significant 16 bits access time). */
    int refcount;
    void *ptr;
} robj;

BPF_HASH(getcall,u64,get_event_t);
BPF_PERF_OUTPUT(duration_events);

int trace_start_time(struct pt_regs* ctx){
    u64 pid = bpf_get_current_pid_tgid();
    u64 start_time_ns = bpf_ktime_get_ns();
    
    robj* rObj = (robj*)PT_REGS_PARM2(ctx);

    get_event_t event = {
        .pid = pid>>32,
        .start_time_ns = start_time_ns,
    };

    event.klen = bpf_probe_read_user_str(&event.key,sizeof(event.key),(void*)rObj->ptr);
    getcall.update(&pid,&event);
    return 0;
}

int send_duration(struct pt_regs* ctx){
    u64 pid = bpf_get_current_pid_tgid();
    get_event_t* eventp = getcall.lookup(&pid);
    if (eventp==0){
        return 0;
    }
    get_event_t event = *eventp;
    event.duration = bpf_ktime_get_ns() - event.start_time_ns;
    duration_events.perf_submit(ctx,&event,sizeof(event));
    getcall.delete(&pid);
    return 0;
}`
